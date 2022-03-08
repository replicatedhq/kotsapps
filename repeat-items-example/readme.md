<h1>Repeatable Items Example App</h1>

This is a sample application that shows an example of using [repeatable items](https://docs.replicated.com/reference/custom-resource-config#repeatable-items). 

In this example, we are using repeatable items to provide one or more ports to use for our `nginx` deployment as well as the corresponding `service`.

First, we define our config fields per the documentation. Here is a snippet of the [config.yaml](manifests/config.yaml) file:

```yaml

    - name: nginx-deployments
      title: Example variadic deployment of nginx
      description: Each provided port spins up another nginx with a NodePort service attached
      items:
        - name: nginx
          title: NginX port
          help_text: Select a port for the nginx deployment(s)
          type: text
          repeatable: true
          templates:
              - apiVersion: v1
                kind: Service
                name: my-nodeport
                yamlPath: spec.ports[0]
              - apiVersion: apps/v1
                kind: Deployment
                name: my-nginx
                yamlPath: spec.template.spec.containers[0].ports[0]
          valuesByGroup:
            nginx-deployments:
              nginx-port: "80"

```
In the yaml above, we have defined two templates. Each one of these has a corresponding manifest that matches the `apiVersion`, `kind`, and `name`. The `yamlPath` is the yaml path to where the value shuld go. In our case, we are targetting an array field called `ports`.

So to summarize, the definition in the `config.yaml` file is for one field that can be repeatable. The value(s) from this field will be populated in both a `deployment` and `service` manifests.

Below is a snippet of the [nginx deployment](manifests/deployment.yaml) that has the templating needed to use the value from the repeatable items:

``` yaml

spec:
  selector:
    matchLabels:
      app: variadic_example
  template:
    metadata:
        labels:
          app: variadic_example
    spec:
      containers:
        - name: nginx
          image: nginx
          ports:
            - containerPort: repl{{ ConfigOption "[[repl .nginx ]]" | ParseInt }}
         
 ```
 
 Here is the snippet for the [nodeport service](manifests/service.yaml). Note that in this one we use the random field name that is assigned to the instance of the field as the value for the `name` field.
 
 ```yaml
 spec:
  type: NodePort
  ports:
  - port: repl{{ ConfigOption "[[repl .nginx ]]" | ParseInt }}
    name: '[[repl .nginx ]]'
 
 ```
 
Once you deploy the application, you will be able to define multiple instances of the NginX Port field. Below is a screenshot of a deployed app with several instances of the field with values populated.
![img/fields-rendered-with-values.png]

When we examine the files in the Admin Console we can see how it has assigned a name to each instance of the field.
!(img/config-w-values.png)

Below is the `service.yaml` file templated with the ports listed.
!(img/service-tempated.png)

When we describe the `pod` and the `service`, we can also see the ports:
!(img/describe-service.png)
!(img/describe-deployment.png)

