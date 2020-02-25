# Example mongoDB application running custom support analyzers using mtools

cd manifests
../common/fetch-default-yaml.sh
helm fetch stable/mongodb
cp ../common/Makefile.default Makefile

