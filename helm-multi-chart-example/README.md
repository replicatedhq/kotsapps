Replicated KOTS Multi Chart Example
==================

Example project to show how to manage an application that is comprised of multiple charts. Keep in mind that my multiple charts, we are not referring to [subcharts](https://helm.sh/docs/chart_template_guide/subcharts_and_globals/#helm), rather multiple independent charts. Think of three microservices, each packaged as its own chart and may or may not contain a subchart.

#### Things to note

- For each chart, there is a KOTS [HelmChart](https://kots.io/reference/v1beta1/helmchart/) definition file in the `kots/manifests/` directory. 
- As of this initial commit, no other KOTS features have been added or configured.


#### Sample Charts

This example uses three charts available from the [prometheus community repo](https://github.com/prometheus-community/helm-charts).

The Charts are:

- [alertmanager](https://github.com/prometheus-community/helm-charts/tree/main/charts/alertmanager)
- [prometheus](https://github.com/prometheus-community/helm-charts/tree/main/charts/prometheus)
- [prometheus-postgres-exporter](https://github.com/prometheus-community/helm-charts/tree/main/charts/prometheus-postgres-exporter)


#### Folder Structure

In this KOTS example, the directory for each of the Helm Charts is at the root of this repostiroy. There is also a `kots` directory which contains the replicated YAML files. The Replicated YAML files include a [KURL](https://kurl.sh/) installer for embedded installs, and in the `kots/manifests/` directory are the YAMLS for the KOTS application. This is also the directory where at build time, Helm should place the packaged Helm Charts.

#### CI/CD Process

The CI/CD process would consist of two steps:

1. For each Chart, package it and put the `tar.gz` output in the `kots/manifests` directory.

2. Run the Replicated CLI to create a release that includes the `tar.gz` for each chart and promote to channel.

Assuming that that the current directory is the root directory, a CI/CD process would look similar to this:

``` shell
  helm dependencies update alertmanager
  helm package alertmanager -d kots/manifests/

  helm dependencies update prometheus
  helm package prometheus -d kots/manifests/

  helm dependencies update prometheus-postgres-exporter
  helm package prometheus-postgres-exporter -d kots/manifests/

  replicated release create --yaml-dir=kots/manifests/ --auto -y
```

### Integrating with CI 

This project contains a basic shell script based on the commands described above. 

### Verifying the Release in Replicated Vendor Portal

Once a release is created in Replicated, log in to the [Replicated Vendor Portal](https://vendor.replicated.com/) to check that for each chart, the proper KOTS `HelmChart` definition file is associated.

![](/docs/multi-chart-release.png)

### Deploying the application

After deplying the application, a quick check of `kubectl get pods` shows that for the most part, everything has been deployed:

![](/docs/multi-chart-pods.png)