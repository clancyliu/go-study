#!/bin/bash

set -e

file=$(defaults read /Applications/Navicat\ for\ MySQL.app/Contents/Info.plist)

regex="CFBundleShortVersionString = \"([^\.]+)"
[[ $file =~ $regex ]]

version=${BASH_REMATCH[1]}

echo "Detected Navicat Premium version $version"

case $version in
    "16")
        file=~/Library/Preferences/com.navicat.NavicatForMySQL.plist
        ;;
    "15")
        file=~/Library/Preferences/com.navicat.NavicatForMySQL.plist
        ;;
    *)
        echo "Version '$version' not handled"
        exit 1
       ;;
esac

echo -n "Reseting trial time..."

regex="([0-9A-Z]{32}) = "
[[ $(defaults read $file) =~ $regex ]]

hash=${BASH_REMATCH[1]}

if [ ! -z $hash ]; then
    defaults delete $file $hash
fi

regex="\.([0-9A-Z]{32})"
[[ $(ls -a ~/Library/Application\ Support/PremiumSoft\ CyberTech/Navicat\ CC/Navicat\ for\ MySQL/ | grep '^\.') =~ $regex ]]

hash2=${BASH_REMATCH[1]}

if [ ! -z $hash2 ]; then
    rm ~/Library/Application\ Support/PremiumSoft\ CyberTech/Navicat\ CC/Navicat\ for\ MySQL/.$hash2
fi

echo " Done"