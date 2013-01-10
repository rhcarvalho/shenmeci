# -*- coding: utf-8 -*-

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
