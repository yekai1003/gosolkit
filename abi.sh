#!/bin/bash
orgfile=$1
targetfile=$2
awk '/const .+ABI = .+/{print substr($4,2,length($4)-2)}' $orgfile > $targetfile
