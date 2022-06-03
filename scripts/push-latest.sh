#!/bin/bash

image=$1
candidate=$2

function usage {
    echo "usage: $0 <docker image> <candidate version>"
    exit 1
}

# ensure parameters for image and candidate are set
if [ -z $image ] || [ -z $candidate ]; then
    usage
fi

# tags usually start with "v"
if [[ ${candidate} == v* ]]; then
    candidate="${candidate:1}"
fi

# pull the latest for the given image
docker pull ${image}:latest

# retrieve the version information from labels on the image
current=`docker inspect -f '{{ index .Config.Labels "org.opencontainers.image.version" }}' ${image}:latest`

# use sort's version sort to compare the candidate and the current version numbers
if [ "$(printf '%s\n' "$candidate" "$current" | sort -V | head -n1)" = "$candidate" ]; then 
    echo "Candidate [${candidate}] is older or less than [${current}], latest will *not* be updated"
else
    echo "Candidate [${candidate}] is newer than [${current}], latest will be updated"
    docker tag ${image}:${candidate} ${image}:latest
    docker push ${image}:latest
fi
