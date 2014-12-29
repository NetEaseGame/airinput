#!/usr/bin/env python
# -*- coding: utf-8 -*-

import json
import flask
import airtest
# import airtest.image as aim
# import cv2
import time
from flask import request

import requests

app = None #airtest.connect(monitor=False)

bp = flask.Blueprint('api', __name__)

@bp.route('/')
def home():
    return 'API documentation'

pressed = False

lasttime = time.time()
eventfd = open('event.txt', 'w')

def runjs(fn, *args):
    global lasttime, eventfd

    argstr = ', '.join(map(repr, args))
    print fn, args
    ctime = time.time()
    print >>eventfd, '%d' %((ctime - lasttime)*1000)
    print >>eventfd, '%s(%s)' %(fn, argstr)
    eventfd.flush()
    lasttime = ctime
    requests.post('http://10.242.134.91:21000/runjs', data='%s(%s)' %(fn, argstr))
    # requests.post('http://10.242.119.210:21000/runjs', data='%s(%s)' %(fn, argstr))

@bp.route('/touch', methods=['POST'])
def touch():
    pt = json.loads(request.form['data'])
    print pt
    runjs('tap', pt['x'], pt['y'], 100)
    return 'ok'

@bp.route('/drag', methods=['POST'])
def drag():
    twopt = json.loads(request.form['data'])
    start, end = twopt['start'], twopt['end']
    runjs('drag', start['x'], start['y'], end['x'], end['y'], 10, 500)
    return 'ok'


@bp.route('/run')
def run_code():
    global app
    code = flask.request.args.get('code')
    try:
        exec code        
    except Exception, e:
        return flask.jsonify(dict(success=False, message=str(e)))
    return flask.jsonify(dict(success=True, message=""))

@bp.route('/connect')
def connect():
    global app
    device = flask.request.args.get('device')
    devno = flask.request.args.get('devno')
    try:
        app = airtest.connect(devno, device=device, monitor=False)
    except Exception, e:
        return flask.jsonify(dict(success=False, message=str(e)))
        
    return flask.jsonify(dict(success=True, message="连接成功"))