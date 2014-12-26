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
    requests.post('http://10.242.116.68:21000/runjs', data='%s(%s)' %(fn, argstr))
    # requests.post('http://10.242.119.210:21000/runjs', data='%s(%s)' %(fn, argstr))

@bp.route('/touch', methods=['POST'])
def touch():
    global pressed
    for pt in json.loads(request.form['event']):
        if pt['pressed']:
            if pressed:
                runjs('move', pt['x'], pt['y'])
                # r = requests.post('http://10.242.116.68:21000/runjs', data='move(%d, %d)' %(pt['x'], pt['y']))
            else:
                runjs('press', pt['x'], pt['y'])
            pressed = True
        else:
            runjs('release')
            pressed = False
    return 'ok'

# @bp.route('/crop')
# def crop():
#     rget = flask.request.args.get

#     filename = rget('filename')
#     screen = rget('screen')
#     x, y = int(rget('x')), int(rget('y'))
#     width, height = int(rget('width')), int(rget('height'))

#     screen_file = screen.lstrip('/').replace('/', os.sep)
#     # screen_path = os.path.join(utils.selfdir(), screen_file)
#     # output_path = os.path.join(utils.workdir(), filename)
#     assert os.path.exists(screen_path)

#     im = cv2.imread(screen_path)
#     cv2.imwrite(output_path, im[y:y+height, x:x+width])
#     return flask.jsonify(dict(success=True, 
#         message="文件已保存: "+output_path.encode('utf-8')))

# @bp.route('/cropcheck')
# def crop_check():
#     rget = flask.request.args.get
    
#     screen = rget('screen')
#     x, y = int(rget('x')), int(rget('y'))
#     width, height = int(rget('width')), int(rget('height'))

#     screen_file = screen.lstrip('/').replace('/', os.sep)
#     screen_path = os.path.join(utils.selfdir(), screen_file)

#     im = cv2.imread(screen_path)
#     im = im[y:y+height, x:x+width]  # crop image
#     siftcnt = aim.sift_point_count(im)
#     return flask.jsonify(dict(siftcnt=siftcnt))

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