#!/bin/bash

GO=''

while getopts d:b:g: flag
do
    case "${flag}" in
        d) distributions=${OPTARG};;
        b) binaries=${OPTARG};;
        g) GO=${OPTARG};;
        *) exit 1;;
    esac
done

[[ -n "$GO" ]] || GO='go'

if [[ -z $distributions && -z $binaries ]]; then
    echo "List of distributions and binaries to generate the goreleaser not provided. Use '-d' to specify the names of the distributions and '-b' to specify the names of binaries. Ex.:"
    echo "$0 -d otelcol"
    exit 1
fi

echo "Artifacts to generate: $distributions,$binaries";

for distribution in $(echo "$distributions,$binaries" | tr "," "\n")
do
    if [[ "$distribution" == "builder" || "$distribution" == "opampsupervisor" ]]; then
      target_path="./cmd"
    else
      target_path="./distributions"
    fi

    if [[ "$distribution" == "otelcol-contrib" ]]; then
        ${GO} run cmd/goreleaser/main.go -d "${distribution}" --generate-build-step > "${target_path}/${distribution}/.goreleaser-build.yaml"
    fi

    ${GO} run cmd/goreleaser/main.go -d "${distribution}" > "${target_path}/${distribution}/.goreleaser.yaml"
done
