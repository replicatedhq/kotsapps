Troubleshooting example: Multi-container Deployments
========================


This repo contains an example deployment and [troubleshoot.sh](https://troubleshoot.sh) Collector. It demonstrates collecting logs from a multi-container deployment with two `initContainer` entries and two `container` entries.


### Testing

Apply the deployment:

```
kubectl apply -f deployment.yaml
```

after a few seconds, run the support bundle

```
kubectl support-bundle support-bundle.yaml
```

Untar the bundle

```
tar xzvf support-bundle.tar.gz app-logs failing-app-logs
```



Review the healthy deployment logs

```
tail app-logs/**/*
```

You should see something like

```
==> app-logs/many-init-containers-d8b7b5b5b-gwp5d/first.log <==
Sat Apr 11 18:08:30 UTC 2020 I'm First...
first done

==> app-logs/many-init-containers-d8b7b5b5b-gwp5d/fourth.log <==
Sat Apr 11 18:17:13 UTC 2020 fourth container waiting
Sat Apr 11 18:17:23 UTC 2020 fourth container waiting
Sat Apr 11 18:17:33 UTC 2020 fourth container waiting

==> app-logs/many-init-containers-d8b7b5b5b-gwp5d/second.log <==
Sat Apr 11 18:08:41 UTC 2020 I'm Second
second done

==> app-logs/many-init-containers-d8b7b5b5b-gwp5d/third.log <==
Sat Apr 11 18:17:12 UTC 2020 third container waiting
Sat Apr 11 18:17:22 UTC 2020 third container waiting
Sat Apr 11 18:17:32 UTC 2020 third container waiting
Sat Apr 11 18:17:42 UTC 2020 third container waiting
```

Review the unhealthy deployment logs


```
tail failing-app-logs/**/*
```

You should see something like

```
==> failing-app-logs/failing-init-container-567c7d6db8-9wftq/first.log <==
Sat Apr 11 18:27:34 UTC 2020 I'm First...
first done

==> failing-app-logs/failing-init-container-567c7d6db8-9wftq/second-previous.log <==
Sat Apr 11 18:27:56 UTC 2020 I'm Second
whoops, something broke!

==> failing-app-logs/failing-init-container-567c7d6db8-9wftq/second.log <==
Sat Apr 11 18:28:22 UTC 2020 I'm Second


==> failing-app-logs/failing-init-container-567c7d6db8-9wftq/third-errors.json <==
[
  "failed to get log stream: container \"third\" in pod \"failing-init-container-567c7d6db8-9wftq\" is waiting to start: PodInitializing"
]
```
