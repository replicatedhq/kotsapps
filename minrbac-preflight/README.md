# minrbac-preflight


This app showcases a minimal example of using

- `requireMinimalRBACPrviliges` flag
- `exclude` filters on Preflight checks with a hidden config option to control whether cluster-admin checks are run


The product is an app that has the following properties

- Rich preflight checks of node resources can be run in an embedded kURL scenario
- for customer-facing UI installs, preflight checks will be be prompted during the UI flow
- End-to-end automatable installs in the minimal RBAC scenario where preflight checks would otherwise need to be run manually


### Automated existing-cluster install

```
export APP_SLUG=myapp
export APP_CHANNEL=Unstable
export SHARED_PASSWORD=password # or something
export LICENSE_FILE=$HOME/Desktop/license.yaml

cat <<EOF > /tmp/config-values.yaml
apiVersion: kots.io/v1beta1
kind: ConfigValues
spec:
  values:
    skip_cluster_preflights:
      value: "1"
EOF

kubectl kots install $APP_SLUG/$APP_CHANNEL \
    --namespace $APP_SLUG \
    --shared-password=$SHARED_PASSWORD \
    --license-file $LICENSE_FILE \
    --config-values /tmp/config-values.yaml
 ```

 This should run the app with no manual intervention
