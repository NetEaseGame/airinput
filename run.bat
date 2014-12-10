@echo off
adb push airtoolbox /data/local/tmp/
adb shell chmod 755 /data/local/tmp/airtoolbox
adb shell /data/local/tmp/airtoolbox
