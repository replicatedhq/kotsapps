#!/bin/bash
set -e


# By tagging the release on the stable brange, it triggers CI to automatically make a stable release on all apps. 
# Note: Tag must start with a "v" for this to work. 
git checkout stable
git tag $1
git push origin stable
