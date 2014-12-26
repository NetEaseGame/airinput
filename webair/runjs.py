#! /usr/bin/env python
# -*- coding: utf-8 -*-
#
# Author:  hzsunshx
# Created: 2014-12-26 16:56

"""
runjs
"""

import requests
import time

for code in open('event.txt'):
    code = code.strip()
    if code.isdigit():
        time.sleep(int(code)/1000.0)
    else:
        print code
        requests.post('http://10.242.116.68:21000/runjs', data='%s' %(code))
