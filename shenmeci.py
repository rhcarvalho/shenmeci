# -*- coding: utf-8 -*-

# Shenmeci: Chinese word segmentation and Chinese-English online dictionary.
# Copyright (C) 2013  Rodolfo Henrique Carvalho
# https://github.com/rhcarvalho/shenmeci
#
# This file is part of Shenmeci.
#
# Shenmeci is free software: you can redistribute it and/or modify
# it under the terms of the GNU General Public License as published by
# the Free Software Foundation, either version 3 of the License, or
# (at your option) any later version.
#
# This program is distributed in the hope that it will be useful,
# but WITHOUT ANY WARRANTY; without even the implied warranty of
# MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
# GNU General Public License for more details.
#
# You should have received a copy of the GNU General Public License
# along with this program.  If not, see <http://www.gnu.org/licenses/>.

import os
import gzip
from collections import defaultdict


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


def dict_path(file):
    parent_directory = os.path.dirname(__file__)
    path = os.path.join(parent_directory, "dict", file)
    return path

def load_cedict():
    global load_cedict
    # TODO: this could be in some sort of configuration file
    cedict_file = "cedict_1_0_ts_utf-8_mdbg.txt.gz"
    cedict_path = dict_path(cedict_file)
    vocabulary = defaultdict(list)
    with gzip.open(cedict_path) as cedict:
        for line in cedict:
            # skip comments
            if line.startswith('#'):
                continue
            simplified_hanzi = line.split()[1].decode('utf-8')
            meaning = line[line.find('/'):].strip()
            vocabulary[simplified_hanzi].append(meaning)
    # transform vocabulary into a dict of (unicode-string, byte-string) pairs
    vocabulary = dict((k, '/'.join(v)) for k, v in vocabulary.iteritems())
    # change global binding to always return the vocabulary without recomputing
    # (this may be not a very good idea)
    load_cedict = lambda: vocabulary
    return vocabulary


class ChineseWordSegmenter(WordSegmenter):
    """ChineseWordSegmenter"""
    def __init__(self):
        super(ChineseWordSegmenter, self).__init__(load_cedict())
    
    def lookup_meaning(self, words):
        return [(word, self.vocabulary.get(word, "?")) for word in words]
