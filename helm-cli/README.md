KOTS HELM CLI
==================

Example project that shows how to package the HELM CLI as a KOTS application. In this example, the HELM CLI is used to install and upgrade a Grafana Helm Chart. 

**NOTE - This is an experimanetal project and NOT the way KOTS is currently designed to deploy Helm Charts. To learn about how KOTS supports Helm Charts, please review https://kots.io/vendor/guides/helm-chart/ and https://kots.io/vendor/helm/using-helm-charts/.**

### Desired Outcome

The desired outcome of this is to learn if it is feasable to have KOTS deploy an application that only consists of the HELM CLI and the charts it will need to install and or upgrade. In this scenario, the deployment of the actual application (Grafana) is managed by the Helm CLI and not KOTS as KOTS will simply manage the deployment of the Helm CLI.

### Helm CLI Container

The container is built using a Dockerfile based on the one found in this [GitHub Repository](https://github.com/alpine-docker/helm/blob/master/README.md). The only modifications are that the `ENTRYPOINT` and `CMD` lines are removed (since these will be passed in the Pod Definition file) and added a `COPY` command to copy the Grafana chart at build time. The Grafana chart is included in the [app](https://github.com/cremerfc/helm-cli-kots/tree/main/app) directory of this repo, which also contains the [Dockerfile](https://github.com/cremerfc/helm-cli-kots/blob/main/app/Dockerfile).

The reason for including the chart at build time is for airgap installations. While the chart could be provided to the end user by other means and then mount it later, this could add complexity and possible points of failure. 

#### Defining the Application in KOTS

An application in KOTS is basically a set of YAML files. In this case, this is a very small and simple application that only consists of the container above. This container will then be tasked with deploying the actual Chart.

Because of the nature of this Application, we are deploying this container as a [Kubernets Job](https://kubernetes.io/docs/concepts/workloads/controllers/job/) and is defined [here](https://github.com/cremerfc/helm-cli-kots/blob/main/manifests/helm-cli-job.yaml). In order to have this Job properly recreated each time an upgrade is performed by KOTS, or the application is "Redeployed", we added the KOTS [label annotation](https://kots.io/vendor/packaging/cleaning-up-jobs/) that tells KOTS to delete the job once it finishes. If this is not included, the Job will remain and any upgrades of the Job will fail.

#### Giving the Job the Proper Permissions

To ensure that the Pod executing the commands has the proper permissions to install/upgrade the Helm chart in the cluster we add the following to the Job defintion file:
```yaml
      serviceAccountName: kotsadm
      automountServiceAccountToken: true

```

#### Executing the Helm Commands

The following lines are added to the Job defintion file:

```yaml

       command: ["/bin/sh"]
       args: ["-c", "helm upgrade --install grafana /apps/grafana-6.1.16.tgz -f kots-values.yaml"]

```

We are using the `Upgrade` command with the `--install` option which allows us to handle both a new install and an upgrade of the chart. The `kots-values.yaml` file is how we pass any and all values we want to override when deploying in this manner.

The `kots-values.yaml` file is mounted as a Kubernetes `ConfigMap` which is defined as part of the application.

Since this is more of an experiment and things are bound to go wrong, we need as much information when troubleshooting. To help us, we added the following to the job definition file:

```yaml
           terminationMessagePolicy: FallbackToLogsOnError
```

This setting will write the last chunck of the pod's log (2048 bytes or 80 lines, whichever is smaller) to the [termination-log](https://kubernetes.io/docs/tasks/debug-application-cluster/determine-reason-pod-failure/). This came in really handy when the container would fail due to a Helm error.


### Managing Images

When it comes to images and KOTS, there are two main use cases that you need to consider:

* How the images will be available to the application at install/deploy time (private images vs. public images)
* How to use the right tag for those images depending on the install/deploy method (online vs. airgap)


#### Online Installs

Since this KOTS application is comprised of only the Helm CLI container, KOTS does not know about the images that will be pulled when the Helm CLI container runs and installs/upgrades the chart. Since this example applicaiton installs the Grafana chart, which pulls only public images no further changes would be needed.  This is assuming, of course, there aren't any other network policies or anything else in place in the environment that would limit the access to these public images.

However, more than likely you will have private images. While you could add image tags to pull directly from your registry and include your secrets, a much better way is to let Replicated manage your images. By letting Replicated manage your images, the pull secrets will become invalid when the license expires and you'll be able to create airgap bundles.

The first step is to set up your private image repo as described [here](https://kots.io/vendor/packaging/private-images/). You have the option of either using Replicated's registry to store your images, or as a proxy to your private images. This will allow the application to leverage KOTS secrets to pull the images at deploy time, instead of yours.

Regardless, we'll need to ensure that the correct image tags are being used by the chart. In the Grafana Chart, the images are managed in the Values file which we can take advantage of by using a [ConfigMap](https://github.com/cremerfc/helm-cli-kots/blob/main/manifests/helm-values-config-map.yaml) to override these values. This is the same ConfigMap that is mounted as the `kots-values.yaml` file that is being passed with the `helm` [command](https://github.com/cremerfc/helm-cli-kots#executing-the-helm-commands).

If we were to pull all of the images referenced in Helm Chart and then push them into an ECR repo '(<aws-account-id>.dkr.ecr.<zone>us-east-2.amazonaws.com/demo-apps/)', we would then handle the `grafana` image in the configMap as shown here:

```yaml
           image:
             repository: proxy.replicated.com/proxy/helm-cli-kots/429114214526.dkr.ecr.us-east-2.amazonaws.com/demo-apps/grafana
      pullSecrets: 
        - kotsadm-replicated-registry
```

#### Airgap Installs

KOTS manages airgap installations by including all of the images (including those for the Application) in an airgap bundle. In a traditional app where KOTS manages the deployment of the app, KOTS will automatically pull those containers when an airgap bundle is built. However, as mentioned above, since our app only consists of the Helm CLI container, that is the only Application container that would be included in the bundle, leaving out all of the containers that the Helm Chart would require.

In order to have KOTS include these containers in the airgap bundle we need to add them to the `AdditionalImages` [section](https://kots.io/reference/v1beta1/application/#additionalimages) of the KOTS [Application](https://kots.io/reference/v1beta1/application/) Defintion file. 

Below are all of the images referenced in the Chart's `Values.yaml` file added to the [replicated-app.yaml](https://github.com/cremerfc/helm-cli-kots/blob/main/manifests/replicated-app.yaml) file:

```yaml
  additionalImages: 
    - grafana/grafana:7.3.5
    - bats/bats:v1.1.0
    - curlimages/curl:7.73.0
    - busybox:1.31.1
    - kiwigrid/k8s-sidecar:1.1.0
    - grafana/grafana-image-renderer:latest
```

Note that all of the image tags above are for public images. If the containers are in the private registry we mentioned above, then the image tags would be as follows:

```yaml
  additionalImages: 
    - <aws-account-id>.dkr.ecr.<zone>us-east-2.amazonaws.com/demo-apps/grafana:7.3.5
    - <aws-account-id>.dkr.ecr.<zone>us-east-2.amazonaws.com/demo-apps/bats:v1.1.0
    - <aws-account-id>.dkr.ecr.<zone>us-east-2.amazonaws.com/demo-apps/curl:7.73.0
    - <aws-account-id>.dkr.ecr.<zone>us-east-2.amazonaws.com/demo-apps/busybox:1.31.1
    - <aws-account-id>.dkr.ecr.<zone>us-east-2.amazonaws.com/demo-apps/k8s-sidecar:1.1.0
    - <aws-account-id>.dkr.ecr.<zone>us-east-2.amazonaws.com/demo-apps/grafana-image-renderer:latest
```

And a reminder that for the above to work, it is assumed that you have set up a connection to this ECR repository as described [here](https://kots.io/vendor/packaging/private-images/).

When KOTS deploys the application to an airgap environment, it will first push all of the images to a local registry. KOTS provides [template function contexts](https://kots.io/reference/template-functions/contexts/) you can use.

In most common scenarios, the application may or may not be installed in an airgap environment. In order to account for this, we use the `{{repl if` statement to determine if there is a [local registry](https://kots.io/reference/template-functions/config-context/#haslocalregistry) defined (this is always true for airgap installs). If true, we use the [local registry address](https://kots.io/reference/template-functions/config-context/#localregistryaddress) to replace the image tag with the local registry tag.

Below is a snippet of the ConfigMap image section:


```yaml
           image:
             repository: {{repl if HasLocalRegistry }}{{repl LocalRegistryAddress}}{{repl else}}grafana{{repl end}}/grafana
             pullSecrets:
             - your-registry-secret-if-needed
```

In the case of a private repository, instead of `grafana` we would replace this with something like `registry.replicated.com/your-app-slug/` if you use the Replicated registry to store the images, or with your own private registry address. Note that you may also need to include a `pullSecret` as well.

```yaml
           image:
             repository: {{repl if HasLocalRegistry }}{{repl LocalRegistryAddress}}{{repl else}}proxy.replicated.com/proxy/helm-cli-kots/<aws-account-id>.dkr.ecr.<zone>us-east-2.amazonaws.com/demo-apps{{repl end}}/grafana
             pullSecrets:
             - your-registry-secret-if-needed
```


#### Overriding Values at Install/Upgrade Time

As mentioned above, values that we want to override at install/upgrade time will be passed in the `kots-values.yaml` file, and its contents are defined in the [helm-values-config-map](https://github.com/cremerfc/helm-cli-kots/blob/main/manifests/helm-values-config-map.yaml) file.

One of the advantages of using KOTS, is that it can provide the end user with a web UI to enter the values to use at runtime. To do this we use the KOTS [Config](https://kots.io/reference/v1beta1/config/) custom resource as defined in the [Config.yaml](https://github.com/cremerfc/helm-cli-kots/blob/main/manifests/config.yaml) file in this repository. Those values are then mapped to the [helm-values-config-map](https://github.com/cremerfc/helm-cli-kots/blob/main/manifests/helm-values-config-map.yaml) file.

### Conclusion 

No real impediments were discovered, other than having to manually manage the images (Public v. Private imgaes, airgap v. online). 

However, all this testing was done with a fairly simple chart. To further test this out, a more complex chart (with lots of hooks!) is needed.

