@echo off

echo Android AirInput -- Initializing...

echo Waiting for device to be connected...
adb.exe wait-for-device

echo - Installing native service...
adb.exe push air-native "/data/local/tmp/air-native"
adb.exe shell chmod 755 "/data/local/tmp/air-native"

echo Starting...
rem adb.exe shell kill -9 "/data/local/tmp/air-native"
rem adb.exe shell "/data/local/tmp/air-native -daemon"
adb.exe shell "/data/local/tmp/air-native -remote=mt.nie.netease.com:5000"

echo Service started successfully.
