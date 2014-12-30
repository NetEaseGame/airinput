# coding: utf8

import flask
from flask import request
# from flask import request, flash, redirect, url_for, render_template
import proto

bp = flask.Blueprint('home', __name__)

# @bp.route('/tmp/<path:path>')
# def static_proxy(path):
#     return flask.send_from_directory(utils.TMPDIR, path)

@bp.route('/')
def home():
    return flask.render_template('index.html', phones=proto.PHONES)

@bp.route('/connect')
def connect():
    if request.headers.getlist("X-Forwarded-For"):
        ip = request.headers.getlist("X-Forwarded-For")[0]
    else:
        ip = request.remote_addr
    serialno = request.args['serialno']
    print ip, request.args['serialno']
    proto.PHONES[serialno] = ip
    proto.save()
    return 'good'

@bp.route('/phone/<serialno>')
def phone(serialno):
    return flask.render_template('phone.html', 
        ip=proto.PHONES[serialno], serialno=serialno)
    
@bp.route('/list')
def list_phones():
    return flask.render_template('list.html', phones=proto.PHONES)

@bp.route('/all')
def all_phones():
    return 'good'
