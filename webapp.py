#!/usr/bin/env python
# -*- coding: utf-8 -*-
import os
from bottle import route, request, static_file

from shenmeci import ChineseWordSegmenter

STATIC_PATH = os.path.join(os.path.dirname(__file__), 'static')


@route('/static/<filepath:path>')
def server_static(filepath):
    return static_file(filepath, root=STATIC_PATH)

@route('/')
def index():
    return static_file('index.html', root=STATIC_PATH)

@route('/segment')
def segment():
    query = request.query.q
    segmenter = ChineseWordSegmenter()
    words = segmenter.segment(query)
    result = [dict(z=z, m=m) for z, m in segmenter.lookup_meaning(words)]
    return {u"r": result}


if __name__ == '__main__':
    import argparse
    parser = argparse.ArgumentParser(description='Run web server.')
    parser.add_argument('--port', type=int, default=8080, help='port to bind the HTTP server')
    args = parser.parse_args()

    import bottle
    bottle.debug(True)
    bottle.run(host='localhost', port=args.port, reloader=True)
