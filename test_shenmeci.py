# -*- coding: utf-8 -*-
import unittest
from shenmeci import WordSegmenter


class BaseWordSegmenterTestCase(unittest.TestCase):
    def setUp(self):
        vocabulary = set("""
        A
        AB
        ABCD
        C
        """.split())
        self.segmenter = WordSegmenter(vocabulary)


class SegmentSingleWordTests(BaseWordSegmenterTestCase):
    def test_empty_input(self):
        sentence = ""
        self.assertEqual([], self.segmenter.segment(sentence))
    
    def test_single_character(self):
        sentence = "A"
        self.assertEqual(["A"], self.segmenter.segment(sentence))
    
    def test_two_character_word(self):
        sentence = "AB"
        self.assertEqual(["AB"], self.segmenter.segment(sentence))


class SegmentTwoWordsTests(BaseWordSegmenterTestCase):
    def test_two_single_character_words(self):
        sentence = "AC"
        self.assertEqual(["A", "C"], self.segmenter.segment(sentence))
    
    def test_two_character_word_plus_single_character_word(self):
        sentence = "ABC"
        self.assertEqual(["AB", "C"], self.segmenter.segment(sentence))


class SegmentThreeWordsTests(BaseWordSegmenterTestCase):
    def test_three_single_character_words(self):
        sentence = "ACD"
        self.assertEqual(["A", "C", "D"], self.segmenter.segment(sentence))


class SegmentLongWordsTests(BaseWordSegmenterTestCase):
    def test_long_word_prefix_wont_match(self):
        sentence = "ABCDZ"
        self.assertEqual(["ABCD", "Z"], self.segmenter.segment(sentence))


if __name__ == '__main__':
    unittest.main()
