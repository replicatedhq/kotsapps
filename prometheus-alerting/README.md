End to End Prometheus Alerting Example
===========

This example demonstrates how to use several KOTS features together with the bundled [kube-prometheus](https://github.com/prometheus-operator/kube-prometheus) stack that comes with [kurl.sh](https://kurl.sh) to configure application-level alerting via webhooks or SMTP emails.


### The Flaky App

In this example, we'll monitor a flaky app called `flaky-app`. Most notably, this app has a indicator `temperature_celsius`. When this value is above `85`, a warning should be triggered. Above a value of `90`, a critical alert should be triggered.

The deployment and service can be found in [./manifests/flaky-app.yaml](./manifests/flaky-app.yaml). The golang source code can be found in [./cmd/flaky-app](./cmd/flaky-app). The frontend is static HTML in [bad_javascript.go](./cmd/flaky-app/bad_javascript.go) (it certainly cannot be described as good javascript by any measure). 

The applicaiton stores a single temperature value in memory and has controls and API endpoints to modify the temperature up or down.

![app dashboard](./doc/healthy-app.png)

Changing the temperature will show a visual state change in the application, and we'll explore how this affects the monitoring systems on the backend.

![app dashboard](./doc/warning-app.png)

Most notably, the applicaiton will expose the value of this temperature at `/metrics` for prometheus to pick up

![app dashboard](./doc/exposed-temp.png)

### Monitoring Metrics

Service monitor

Changing temp

Graphs in Prom

Graphs in KOTS Dashboard

### Temperature Alerting

Alert Rules CR


Alerts in Prom

Alerts in 



### Configuring Alert Sinks

Webhooks and SMTP

Using RequestBin and Mailcatcher.me to preview

Config Screen

Alert Manager Secret (patch config, may not work with future versions of prometheus) -- alternative CRD method might work, but not tested


Alertmanager rendered config

Previewing alerts


![app dashboard](./doc/app-dashboard.png)

