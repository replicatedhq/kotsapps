#!/bin/sh

#Very basic script that shows

helm dependencies update alertmanager
helm package alertmanager -d kots/manifests/

helm dependencies update prometheus
helm package prometheus -d kots/manifests/

helm dependencies update prometheus-postgres-exporter
helm package prometheus-postgres-exporter -d kots/manifests/

replicated release create --yaml-dir=kots/manifests/ --auto -y
