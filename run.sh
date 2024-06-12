#!/usr/bin/env zsh

os=$(uname -s | tr '[:upper:]'  '[:lower:]')
arch=$(arch)
app="${0:a:h}/build/couture_${os}_${arch}"
[[ -f "${app}" && -x "${app}" ]] || {
    echo "Error: ${app} not found or not executable"
    exit 1
}

"${app}" "$@"
