# https://kudo.dev/docs/developing-operators/getting-started.html#package-structure
# Kotsadm: http://34.83.176.50:8800 / yqKI5ywYT

# Prerequisites

1. Install the KUDO kubectl plugin. 
```s
brew tap kudobuilder/tap
brew install kudo-cli
``` 
2. Initialize KUDO and create Kudo manifest (including CRDs): 
```s
kubectl kudo init --dry-run --output yaml > ./manifests/kudo.yaml
```
3. Create a VPC in GCE that is airgapped (or close enough, only allowing 22 and 8800 ingress)
4. Create an instance in GCE that uses the VPC. 
5. Install to the airgapped VM
```s
curl -sSL https://k8s.kurl.sh/kudo-nginx-operator-austin | sudo bash
```
6. Create manifests for the operator and deployments
```s
kubectl kudo install ./manifests/kudo-operator --dry-run --output yaml > ./manifests/kudo-sdf.yaml
```
