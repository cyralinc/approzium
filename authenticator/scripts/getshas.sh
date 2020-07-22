#!/usr/bin/env bash

cd pkg
rm *.zip
for D in `find . -type d`
do
  shasum -a 256 $D/authenticator*
done
cd -
