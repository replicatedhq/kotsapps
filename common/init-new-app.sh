 #!/bin/bash
set -e

# Usage info and args
usage()
{
cat << EOF
usage: $0 <app_slug>

Initialize a new kotsapp directory matching the provided app_slug value. 
EOF
}
if [[ -z $1 ]] 
then
     usage
     exit 1
else 
    app_slug=$1
fi

# Make the directory based on the appname, starting from git root
cd $(git rev-parse --show-toplevel)
mkdir $app_slug
cd $app_slug

# Copy Templates
cp ../common/Makefile.default ./Makefile 
cp ../common/workflow.default.yml ../.github/workflows/${app_slug}.yml

# Set workflow based on current directory
sed -i'' -e "s/{{app_slug}}/${app_slug}/g" ../.github/workflows/${app_slug}.yml
rm ../.github/workflows/${app_slug}.yml-e

# Get manifests from default yml location
git clone git@github.com:replicatedhq/kots-default-yaml.git 
mkdir manifests
cp ./kots-default-yaml/*.yaml ./manifests
rm -rf kots-default-yaml

