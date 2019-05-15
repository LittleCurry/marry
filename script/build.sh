#!/bin/bash

git pull origin develop

function build_docker {
    TAG=$1

    IMAGE_NAME="101.37.151.52:5000/cloud_server/me_photo"

    IMAGE_FULL_NAME="$IMAGE_NAME:$TAG"
    HAS_OLD_IMAGES=$(docker images|grep $IMAGE_NAME|grep $TAG|wc -l)
    echo $HAS_OLD_IMAGES

    if [ $HAS_OLD_IMAGES -ne "0" ]; then
        echo "Remove docker image..."
        docker rmi -f $IMAGE_FULL_NAME
    fi

    echo "Building docker image..."
    docker build -t $IMAGE_FULL_NAME .
    echo "Push image to reigstry"
    docker push $IMAGE_FULL_NAME

    if [ "$TAG" != "latest" ]; then
    echo "Remove $IMAGE_FULL_NAME"
    docker rmi -f $IMAGE_FULL_NAME
    fi


}

set -e

apidoc -i $GOPATH/src/gitlab.com/SiivaVideoStudio/cloud_server/me_photo/ -o $GOPATH/src/gitlab.com/SiivaVideoStudio/cloud_server/me_photo/apidoc

DATETAG=$(date +"%y%m%d%H%M%S")

cd $GOPATH/src/gitlab.com/SiivaVideoStudio/cloud_server/me_photo

echo "Building application..."
CGO_ENABLED=0 GOOS=linux go build -o resource/main .

build_docker $DATETAG

build_docker "latest"

docker rmi $(docker images -f "dangling=true" -q)

rm -r ./apidoc

echo "Cleanup resources..."
rm resource/main

echo "Done"
