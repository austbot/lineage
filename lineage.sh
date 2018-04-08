#!/usr/bin/env bash
line=$(head -n 1 $1)
whtfile="whitelist.txt"
function li() {
    echo "info: ${1}"
}

function le() {
    echo "error: ${1}"
    if [[ $2 ]]; then
        exit $2
    fi
}

function chkw() {
    if [ -f "${whtfile}" ]; then
      li "Found ${whtfile}"
    else
      le "${whtfile} missing" 1
    fi
}

li "Found base image ${line}"
if [[ $line =~ ^FROM\ (.*)$ ]]; then
    image="${BASH_REMATCH[1]}"
    li "checking ${image} for approval"
    if [[ $image =~ ^[^\/]*\/.*$ ]]; then
        li "Image is not from std library, checking manual whitelist"
        match=`grep $image "${whtfile}"`
        if [ $match ]; then
            li "Image approved match found"
            exit 0
        else
            le "Image not approved" 2
        fi
    else
        li "Image approved its from the std library on docker hub"
    fi
    exit 0
fi
