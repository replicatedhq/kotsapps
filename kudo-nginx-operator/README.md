# https://kudo.dev/docs/developing-operators/getting-started.html#package-structure

# Existing Cluster Install Steps
1. Run: 
```s
curl https://kots.io/install | bash kubectl kots install kudo-nginx-operator
```
2. Follow steps on command line, and upload license file (License-Stable.yaml) to admin console

# New VM Install Steps (Embed the k8s cluster along with the application):
1. Run: 
```s
curl -sSL https://k8s.kurl.sh/kudo-nginx-operator | sudo bash
```
2. Follow steps on command line, and upload license file (License-Stable.yaml) to admin console

# Airgapped Install Steps onto a new machine
1. Create an airgapped network, or one that is close enough to verify all outbound is denied with minimal ingress (e.g., only allow 8800/22 ingress). 
2. Create an instance within the network. 
3. Download airgap installer from vendor onto your desktop, and transfer the installer to the machine:
```s
curl -o replicated-installer.tar.gz https://kurl.sh/bundle/kudo-nginx-operator-austin
gcloud compute scp replicated-installer.tar.gz ${INSTANCE}:replicated-installer.tar.gz --zone $ZONE
```
4. Logon to the instance, extract the installer and begin installation: 
```s
mkdir replicated && mv replicated-installer.tar.gz replicated && cd replicated
tar xvf replicated-installer.tar.gz
cat ./install.sh | sudo bash -s airgap
```
5. Follow steps on command line, and upload license file (License-Stable.yaml) and airgap bundle to admin console
6. Verify you have kudo instances generated: 
```s
kubectl kudo get instances 
```
7. Verfy you have two kudo pods: 
```s
kubectl get pods 
```

# How I generated the Operator YAML using KUDO

## Prerequisites

The KUDO kubectl plugin. 
```s
brew tap kudobuilder/tap
brew install kudo-cli
```
or 
```s
kubectl krew install kudo
```

## Steps: 
1. Initialize KUDO and create Kudo manifest (including CRDs). Or use what is in repo 
```s
kubectl kudo init --dry-run --output yaml > ./manifests/kudo.yaml
```
2. Extract the YAML from those manifests (for some reason, kudo doesn't have dry-run install). Replace first three vars with appropriate values
```s
OPERATOR_FILENAME=kudo-operator.yaml
NAMESPACE=kudo
OPERATOR_NAME=first-operator

kudo install ./kudo-operator
cd manifests
kubectl get operators.kudo.dev $OPERATOR_NAME -n $NAMESPACE --export -n $NAMESPACE -o yaml > $OPERATOR_FILENAME && echo --- >> $OPERATOR_FILENAME
TMP_OPVERNAME=$(kubectl get operatorversions.kudo.dev -n kudo -o jsonpath="{.items[?(@.spec.operator.name==\"${OPERATOR_NAME}\")].metadata.name}")
kubectl get operatorversions.kudo.dev ${TMP_OPVERNAME} --export -n $NAMESPACE -o yaml >> $OPERATOR_FILENAME && echo --- >> $OPERATOR_FILENAME
kubectl get instances.kudo.dev ${OPERATOR_NAME}-instance --export -n $NAMESPACE -o yaml >> $OPERATOR_FILENAME && echo --- >> $OPERATOR_FILENAME
cd ..
```
3. Remove kudo and all yaml from the cluster
```s
kubectl kudo init --dry-run --output yaml | kubectl delete -f -
```
