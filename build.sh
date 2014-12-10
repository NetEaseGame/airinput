#! /bin/sh
#
# build.sh
# Copyright (C) 2014 hzsunshx <hzsunshx@onlinegame-13-180>
#
# Distributed under terms of the MIT license.
#

export GOOS=android
export GOARCH=arm

CGO_ENABLED=1 go build
