KUDO Nginx Operator
==================

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
