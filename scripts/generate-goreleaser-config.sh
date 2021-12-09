#!/bin/bash
REPO_DIR="$( cd "$(dirname $( dirname "${BASH_SOURCE[0]}" ))" &> /dev/null && pwd )"
GEN_YAML_DIR="${REPO_DIR}/.generated-yaml"
TEMPLATES_DIR="${REPO_DIR}/scripts/goreleaser-templates"
MAIN_TEMPLATE="${TEMPLATES_DIR}/goreleaser.yaml"
YQ=$(which yq)

while getopts d:y: flag
do
    case "${flag}" in
        d) distributions=${OPTARG};;
        y) YQ=${OPTARG};;
    esac
done

if [[ -z $distributions ]]; then
    echo "List of distributions to use with goreleaser not provided. Use '-d' to specify the names of the distributions. Ex.:"
    echo "$0 -d opentelemetry-collector"
    exit 1
fi

mkdir -p "${GEN_YAML_DIR}"
touch "${GEN_YAML_DIR}/last-generation"

templates=$(ls ${TEMPLATES_DIR}/*.template.yaml | xargs -n 1 basename | sed 's/.template.yaml//gi')
for template in $templates
do
    for distribution in $(echo $distributions | tr "," "\n")
    do
        export CONTAINER_BASE_NAME="otel/{distribution}"
        DIST_CONF="${REPO_DIR}/distributions/${distribution}/distribution.conf"

        if [[ -f "${DIST_CONF}" ]]; then
            set -o allexport
            source "${DIST_CONF}"
            set +o allexport
        fi

        sed "s/{distribution}/${distribution}/gi" "${TEMPLATES_DIR}/${template}.template.yaml" | envsubst > "${GEN_YAML_DIR}/${distribution}-${template}.yaml"
        if [[ $? -ne 0 ]]; then
            echo "❌ ERROR: failed to generate '${template}' YAML snippets for '${distribution}'."
            exit 1
        fi
    done
done

set -e -o pipefail

${YQ} eval-all '. as $item ireduce ({}; . *+ $item)' "${MAIN_TEMPLATE}" "${GEN_YAML_DIR}"/*.yaml > .goreleaser.yaml
echo "✅ SUCCESS: goreleaser YAML generated"
