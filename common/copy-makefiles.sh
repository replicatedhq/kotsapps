 #!/bin/bash
set -e

cd $(git rev-parse --show-toplevel)
for d in */; do cp ./common/Makefile.default "$d/Makefile"; done