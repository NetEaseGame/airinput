#!/usr/bin/env python
# -*- coding: utf-8 -*-
#
# webair
#
import flask

app = flask.Flask(__name__)

from routers import home, api
app.register_blueprint(home.bp, url_prefix='')
app.register_blueprint(api.bp, url_prefix='/api')


if __name__ == '__main__':
    app.run(port=5000, host='0.0.0.0', debug=True)
