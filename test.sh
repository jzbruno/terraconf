#!/usr/bin/env bash

set -euo pipefail

coverageFile="coverage.txt"
profileFile="profile.out"

echo "" > ${coverageFile}

for d in $(go list ./... | grep -v vendor); do
    go test -race -coverprofile=${profileFile} -covermode=atomic ${d}

    if [[ -f ${profileFile} ]]; then
        cat ${profileFile} >> ${coverageFile}
        rm ${profileFile}
    fi
done
