# https://kudo.dev/docs/developing-operators/getting-started.html#package-structure

# Prerequisites

The KUDO kubectl plugin. 
```s
brew tap kudobuilder/tap
brew install kudo-cli
```
or 
```s
kubectl krew install kudo
```


# Setup
1. Create a VPC in GCE that is airgapped (or close enough, only allowing 22 and 8800 ingress/outbound. All other outbound denied)
2. Create an instance in GCE that uses the VPC. 
3. Download airgap installer from vendor onto your desktop, and transfer the installer to the machine:
```s
INSTANCE=austins-kotsapps-kudo
ZONE=us-west1-b
PROJECT=smart-proxy-839

curl -o replicated-installer.tar.gz https://kurl.sh/bundle/kudo-nginx-operator-austin
gcloud compute scp replicated-installer.tar.gz ${INSTANCE}:replicated-installer.tar.gz --zone $ZONE
```
4. Logon to the instance, extract the installer and begin installation: 
```s
gcloud beta compute ssh --zone $ZONE $INSTANCE --project $PROJECT
mkdir replicated && mv replicated-installer.tar.gz replicated && cd replicated
tar xvf replicated-installer.tar.gz
cat ./install.sh | sudo bash -s airgap
```
5. Upload license and airgap bundle. 

5. Verify you have kudo instances generated: 
```s
kubectl kudo get instances 
```
5. Verfy you have two kudo pods: 
```s
kubectl get pods 
```


# How I generated the Operator YAML using KUDO
1. Initialize KUDO and create Kudo manifest (including CRDs). Or use what is in repo 
```s
kubectl kudo init --dry-run --output yaml > ./manifests/kudo.yaml
```
2. Extract the YAML from those manifests (for some reason, kudo doesn't have dry-run install). Replace first three vars with appropriate values
```s
OPERATOR_FILENAME=kudo-operator.yaml
NAMESPACE=kudo
OPERATOR_NAME=first-operator

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
