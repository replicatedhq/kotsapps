# https://kudo.dev/docs/developing-operators/getting-started.html#package-structure
# Kotsadm: http://34.83.176.50:8800 / yqKI5ywYT

# Prerequisites

The KUDO kubectl plugin. 
```s
brew tap kudobuilder/tap
brew install kudo-cli
```


# Setup
1. Create a VPC in GCE that is airgapped (or close enough, only allowing 22 and 8800 ingress)
2. Create an instance in GCE that uses the VPC. 
3. Install to the airgapped VM
```s
curl -sSL https://k8s.kurl.sh/kudo-nginx-operator-austin | sudo bash
```
4. Verify you have kudo instances generated: 
```s
kubectl kudo get instances -n kudo
```
5. Verfy you have two kudo pods: 
```s
kubectl get pods -n kudo
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