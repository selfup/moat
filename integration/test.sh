#!/usr/bin/env bash

set -e

go install

if [[ -d archive ]]
then
    rm -rf archive
fi

if [[ -d fixtures ]]
then
    rm -rf fixtures
fi

WOW="wow this is going to be encrypted and saved to a cloud service directory"

moat -home="archive" -service="fixtures"

echo $WOW >> archive/Moat/wow.txt

moat -home="archive" -service="fixtures" -cmd=push

rm archive/Moat/wow.txt

moat -home="archive" -service="fixtures" -cmd=pull

cat archive/Moat/wow.txt | grep -q "$WOW"

if [[ $? == '0' ]]
then
    echo "moat passed"
else
    echo "moat failed"
    exit 1
fi

rm -rf archive fixtures
