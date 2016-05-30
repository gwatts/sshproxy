#!/bin/bash

# Simple script to build releases for various platforms

set -e

BASEDIR=/tmp/sshproxy-release

mkdir -p $BASEDIR/src $BASEDIR/pkg $BASEDIR/bin
mkdir -p $BASEDIR/src/github.com/gwatts/sshproxy

export GOPATH=$BASEDIR
cd $BASEDIR/src/github.com/gwatts
git clone https://github.com/gwatts/sshproxy.git
cd $BASEDIR/src/github.com/gwatts/sshproxy
glide install


VERSION=$(perl -ne '/version\s+=\s+"([^"]+)/ && print "$1"' $BASEDIR/src/github.com/gwatts/sshproxy/version.go)

if [ -z $VERSION ]; then
    echo "Failed to detect version"
    exit 1
fi

echo "Processing $VERSION"

cd $BASEDIR
rm -f sshproxy sshproxy.exe
mkdir $VERSION

echo "Linux build"
GOOS=linux GOARCH=amd64 go build -v github.com/gwatts/sshproxy
tar -cvzf $VERSION/sshproxy-linux-${VERSION}.tar.gz ./sshproxy
rm sshproxy

echo "Mac"
GOOS=darwin GOARCH=amd64 go build -v github.com/gwatts/sshproxy
tar -cvzf $VERSION/sshproxy-mac-${VERSION}.tar.gz ./sshproxy
rm sshproxy

echo "Windows 64bit"
GOOS=windows GOARCH=amd64 go build -v github.com/gwatts/sshproxy
zip $VERSION/sshproxy-windows-64bit-${VERSION}.zip ./sshproxy.exe
rm sshproxy.exe

echo "Windows 32bit"
GOOS=windows GOARCH=386 go build -v github.com/gwatts/sshproxy
zip $VERSION/sshproxy-windows-32bit-${VERSION}.zip ./sshproxy.exe
rm sshproxy.exe
