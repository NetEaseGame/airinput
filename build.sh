#! /bin/sh
#
# build.sh
# Copyright (C) 2014 hzsunshx <hzsunshx@onlinegame-13-180>
#
# Distributed under terms of the MIT license.
#

GOOS=android GOARCH=arm CGO_ENABLED=1 go build -o air-native
cp air-native dist/
