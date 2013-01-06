#!/usr/bin/env python
# -*- coding: utf-8 -*-
import os
from bottle import route, request, static_file

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


class WordSegmenter(object):
    def __init__(self, vocabulary):
        self.vocabulary = vocabulary
    
    def segment(self, sentence):
        """Break a sentence into a list of words."""
        if not sentence:
            return []
        if sentence in self.vocabulary:
            return [sentence]
        else:
            longest_word = sentence[:1]
            for i in xrange(2, len(sentence)):
                maybe_word = sentence[:i]
                if maybe_word in self.vocabulary:
                    longest_word = maybe_word
                    continue
                else:
                    break
            return [longest_word] + self.segment(sentence[len(longest_word):])


class ChineseWordSegmenter(WordSegmenter):
    """ChineseWordSegmenter"""
    def __init__(self):
        super(ChineseWordSegmenter, self).__init__(self.load_chinese_vocabulary())
    
    def load_chinese_vocabulary(self):
        try:
            return ChineseWordSegmenter.__vocabulary
        except AttributeError:
            import gzip
            vocabulary = dict()
            root_dir = os.path.dirname(__file__)
            cedict_path = os.path.join(root_dir, 'dict', 'cedict_1_0_ts_utf-8_mdbg.txt.gz')
            with gzip.open(cedict_path) as cedict:
                for line in cedict:
                    if line.startswith('#'):
                        continue
                    simplified_hanzi = line.split()[1].decode('utf-8')
                    meaning = line[line.find('/'):].strip()
                    vocabulary[simplified_hanzi] = meaning
            ChineseWordSegmenter.__vocabulary = vocabulary
            return vocabulary
            
    def lookup_meaning(self, words):
        return [(word, self.vocabulary.get(word, "?")) for word in words]


if __name__ == '__main__':
    import argparse
    parser = argparse.ArgumentParser(description='Process some integers.')
    parser.add_argument('--port', type=int, default=8080, help='port to bind the HTTP server')
    args = parser.parse_args()

    import bottle
    bottle.debug(True)
    bottle.run(host='localhost', port=args.port, reloader=True)
