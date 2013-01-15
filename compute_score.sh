#!/usr/bin/env bash

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

DATA=icwb2-data
GOLD=$DATA/gold
TESTING=$DATA/testing

SEGMENTED=$DATA/segmented.utf8
RESULT=$DATA/result.utf8

score=$DATA/scripts/score
shenmeci=./shenmeci.py

$shenmeci $TESTING/cityu_test.utf8 $SEGMENTED &&

perl $score $GOLD/{cityu_training_words.utf8,cityu_test_gold.utf8} $SEGMENTED | tee $RESULT | grep ===
