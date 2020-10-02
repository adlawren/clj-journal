#!/bin/bash

echo
echo "cli-lein-template :: Creating uberjar..."
echo

lein uberjar


echo
echo "cli-lein-template :: Create native executable..."
echo

mkdir -p build
cd build

# @todo: use the correct version
NAME=app-0.1.0-SNAPSHOT-standalone

echo "Password needed to run docker:"
sudo docker run -it \
        --rm \
        -v $HOME:$HOME \
        -w `pwd` \
        -v `pwd`:`pwd` \
        -u $UID:$GID \
        quay.io/quarkus/centos-quarkus-native-s2i:graalvm-1.0.0-rc16 \
        sh -c "native-image -H:+ReportUnsupportedElementsAtRuntime -jar ../target/$NAME.jar"

mv ./$NAME ../clj-journal

cd ../
rm -r build
