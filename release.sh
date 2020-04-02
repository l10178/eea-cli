#!/bin/bash

# local release

current_dir=$(cd $(dirname $0); pwd)

cd ${current_dir}

goreleaser=/usr/bin/goreleaser

${goreleaser} --snapshot --skip-publish --rm-dist

cp -f ${current_dir}/dist/eea_linux_amd64/eea /usr/local/bin/eea && chmod +x /usr/local/bin/eea
cp -f ${current_dir}/dist/eear_linux_amd64/eear /usr/local/bin/eear && chmod +x /usr/local/bin/eear

