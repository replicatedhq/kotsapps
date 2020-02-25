 #!/bin/bash
set -e

git clone git@github.com:replicatedhq/kots-default-yaml.git 
cp ./kots-default-yaml/*.yaml .
rm -rf kots-default-yaml