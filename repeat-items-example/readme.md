<h1>Repeatable Items Example App</h1>

This is a sample application that shows an example of using [repeatable items](https://docs.replicated.com/reference/custom-resource-config#repeatable-items). 

In this example, we are using repeatable items to provide one or more ports to use for our `nginx` deployment as well as the corresponding `service`.

First, we define our config fields per the documentation:

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
