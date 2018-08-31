#!/bin/sh

#if [ -z "$GOPATH" ]; then
  #printf "No GOPATH detected! "
  #printf "Setup GOPATH before build: export GOPATH=/your/go/path\\n"
  #exit 1
#fi

#if [ ! -d "$GOPATH/src/golang.org/x/oauth2" ];then
  #git clone https://github.com/golang/oauth2.git $GOPATH/src/golang.org/x/oauth2
  #[ $? -ne 0 ] && exit 1
#fi

#if [ ! -d "$GOPATH/src/leb.io/aeshash" ]
  #git clone https://github.com/tildeleb/aeshash.git $GOPATH/src/leb.io/aeshash
  #[ $? -ne 0 ] && exit 1
#fi

#if [ ! -d "$GOPATH/src/leb.io/cuckoo" ]
  #git clone https://github.com/tildeleb/cuckoo.git $GOPATH/src/leb.io/cuckoo
  #[ $? -ne 0 ] && exit 1
#fi

#if [ ! -d "$GOPATH/src/leb.io/hashland" ]
  #git clone https://github.com/tildeleb/hashland.git $GOPATH/src/leb.io/hashland
  #[ $? -ne 0 ] && exit 1
#fi

#if [ ! -d "$GOPATH/src/leb.io/hrff" ]
  #git clone https://github.com/tildeleb/hrff.git $GOPATH/src/leb.io/hrff
  #[ $? -ne 0 ] && exit 1
#fi

#if [ ! -d "$GOPATH/src/leb.io/hashes" ]
  #mkdir $GOPATH/src/leb.io/hashes
  #cp -r $GOPATH/src/hashland/siphash $GOPATH/src/leb.io/hashes/
  #[ $? -ne 0 ] && exit 1
#fi

#if [ ! -d "$GOPATH/src/cloud.google.com/go" ]
  #git clone https://github.com/GoogleCloudPlatform/google-cloud-go.git $GOPATH/src/cloud.google.com/go
  #[ $? -ne 0 ] && exit 1
#fi

#if [ ! -d "$GOPATH/src/google.golang.org/api" ]
  #git clone https://github.com/google/google-api-go-client.git $GOPATH/src/google.golang.org/api
  #[ $? -ne 0 ] && exit 1
#fi

#if [ ! -d "$GOPATH/src/golang.org/x/build" ]
  #git clone https://github.com/golang/build.git $GOPATH/src/golang.org/x/build
  #[ $? -ne 0 ] && exit 1
#fi

glide install
