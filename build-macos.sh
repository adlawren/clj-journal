#!/usr/env/bin bash

GRAAL_VER="20.2.0"
GRAAL_DIR="graalvm-ce-java11-${GRAAL_VER}"
NAME="app-0.1.0-SNAPSHOT-standalone"
PROJECT_NAME="$(basename $(pwd))"

export PATH="/Library/Java/JavaVirtualMachines/graalvm-ce-java11-${GRAAL_VER}/Contents/Home/bin:${PATH}"

# Install GraalVM if it isn't installed
if [ ! $(which gu) ]
then
    # See https://www.graalvm.org/docs/getting-started-with-graalvm/macos/

    GRAAL_TARBALL="graalvm-ce-java11-darwin-amd64-${GRAAL_VER}.tar.gz"

    wget "https://github.com/graalvm/graalvm-ce-builds/releases/download/vm-${GRAAL_VER}/${GRAAL_TARBALL}"

    tar -xvzf "${GRAAL_TARBALL}"

    echo "Enter password to move GraalVM to virtual machines directory:"
    sudo mv "${GRAAL_DIR}" /Library/Java/JavaVirtualMachines/

    rm -f "${GRAAL_TARBALL}"
fi

gu install native-image

lein uberjar

native-image --initialize-at-build-time -H:+ReportUnsupportedElementsAtRuntime -jar "./target/${NAME}.jar"

mv "${NAME}" "${PROJECT_NAME}-macos"
