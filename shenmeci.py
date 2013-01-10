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
        words = []
        while sentence:
            longest_word = sentence[:1]
            for i in xrange(2, len(sentence) + 1):
                maybe_word = sentence[:i]
                if maybe_word in self.vocabulary:
                    longest_word = maybe_word
                else:
                    break
            words.append(longest_word)
            sentence = sentence[len(longest_word):]
        return words


class WordTranslator(object):
    def __init__(self, vocabulary):
        self.vocabulary = vocabulary

    def lookup_meanings(self, words):
        return [(word, self.vocabulary.get(word, "?")) for word in words]


class CEDICT(object):
    def __init__(self):
        # TODO: this could be in some sort of configuration file
        parent_directory = os.path.dirname(__file__)
        dictionaries_directory = os.path.join(parent_directory, "dict")
        cedict_file = "cedict_1_0_ts_utf-8_mdbg.txt.gz"
        self._cedict_path = os.path.join(dictionaries_directory, cedict_file)
        self.vocabulary = []
        self.load_cedict()

    def load_cedict(self):
        """Load a CEDICT Chinese-English dictionary file."""
        cedict_path = self._cedict_path
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
        self.vocabulary = vocabulary


cedict = CEDICT()
ChineseWordSegmenter = WordSegmenter(cedict.vocabulary)
ChineseEnglishWordTranslator = WordTranslator(cedict.vocabulary)
