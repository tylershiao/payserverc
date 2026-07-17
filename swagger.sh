#!/bin/bash
if [[ $# -gt 0 ]]
then
        cd docs/$1
        swagger generate spec -o swagger.json
        exit 0
fi

cd docs/web
swagger generate spec -o swagger.json
cd ../backend 
swagger generate spec -o swagger.json