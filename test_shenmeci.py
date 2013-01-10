# -*- coding: utf-8 -*-
import unittest
from shenmeci import DAWGWordSegmenter as WordSegmenter


class BaseWordSegmenterTestCase(unittest.TestCase):
    def setUp(self):
        vocabulary = set(u"""
        A
        AB
        ABCD
        C
        """.split())
        self.segmenter = WordSegmenter(vocabulary)


class SegmentSingleWordTests(BaseWordSegmenterTestCase):
    def test_empty_input(self):
        sentence = u""
        self.assertEqual([], self.segmenter.segment(sentence))
    
    def test_single_character(self):
        sentence = u"A"
        self.assertEqual([u"A"], self.segmenter.segment(sentence))
    
    def test_two_character_word(self):
        sentence = u"AB"
        self.assertEqual([u"AB"], self.segmenter.segment(sentence))


class SegmentTwoWordsTests(BaseWordSegmenterTestCase):
    def test_two_single_character_words(self):
        sentence = u"AC"
        self.assertEqual([u"A", u"C"], self.segmenter.segment(sentence))
    
    def test_two_character_word_plus_single_character_word(self):
        sentence = u"ABC"
        self.assertEqual([u"AB", u"C"], self.segmenter.segment(sentence))


class SegmentThreeWordsTests(BaseWordSegmenterTestCase):
    def test_three_single_character_words(self):
        sentence = u"ACD"
        self.assertEqual([u"A", u"C", u"D"], self.segmenter.segment(sentence))


class SegmentLongWordsTests(BaseWordSegmenterTestCase):
    def test_long_word_prefix_wont_match(self):
        sentence = u"ABCDZ"
        self.assertEqual([u"ABCD", u"Z"], self.segmenter.segment(sentence))


if __name__ == '__main__':
    unittest.main()
