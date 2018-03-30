#!/usr/bin/env bash

set -euo pipefail

coverageFile="c.out"
coveragePartialFile="c.partial"

if [[ -f ${coverageFile} ]]; then
	rm ${coverageFile}
fi
touch ${coverageFile}

for d in $(go list ./... | grep -v vendor); do
    go test -coverprofile=${coveragePartialFile} -covermode=atomic ${d}

    if [[ -f ${coveragePartialFile} ]]; then
        cat ${coveragePartialFile} >> ${coverageFile}
        rm ${coveragePartialFile}
    fi
done
