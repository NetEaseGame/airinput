#!/usr/bin/env python
# -*- coding: utf-8 -*-

import json
import flask

from flask import request
import requests

import proto

app = None #airtest.connect(monitor=False)

bp = flask.Blueprint('api', __name__)

@bp.route('/')
def home():
    return 'API documentation'

def runjs(ip, fn, *args):
    argstr = ', '.join(map(repr, args))
    requests.post('http://{ip}:21000/runjs'.format(ip=ip), data='%s(%s)' %(fn, argstr))

@bp.route('/<serialno>/touch', methods=['POST'])
def touch(serialno):
    ip = proto.PHONES[serialno]
    pt = json.loads(request.form['data'])
    print pt
    runjs(ip, 'tap', pt['x'], pt['y'], 100)
    return 'ok'

@bp.route('/<serialno>/drag', methods=['POST'])
def drag(serialno):
    ip = proto.PHONES[serialno]
    twopt = json.loads(request.form['data'])
    start, end = twopt['start'], twopt['end']
    runjs(ip, 'drag', start['x'], start['y'], end['x'], end['y'], 10, 500)
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
