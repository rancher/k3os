#!/usr/bin/env bash

for key in $( set | awk '{FS="="}  /^OS_/ {print $1}' ); do echo $key ; done

