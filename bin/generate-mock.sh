#!/usr/bin/env bash

set -euo pipefail

for file in `find . -name '*.go' | grep service`; do
    if `grep -q 'interface {' ${file}`; then
        dest=${file//service\//}
        mockgen -source=${file} -destination=test/mock/${dest}
    fi
done
