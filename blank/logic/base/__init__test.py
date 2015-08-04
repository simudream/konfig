#!/usr/bin/env python
from __future__ import absolute_import
from __future__ import division
from __future__ import print_function
from __future__ import unicode_literals
from __future__ import with_statement

import sys
import StringIO
import unittest
import os.path
import json
sys.path.append(os.path.join(os.path.dirname(os.path.realpath(__file__)), '..'))

import base

class TestBase(unittest.TestCase):
    def setUp(self):
        s = StringIO.StringIO('{"users": [{"name": "didip"}, {"name": "brotato"}]}')
        sys.stdin = s
        self.logic = base.Base()
        sys.stdin = sys.__stdin__

    def test_set_data_from_stdin(self):
        self.assertEqual(len(self.logic.data["users"]), 2)
        self.assertEqual(self.logic.data["users"][0]["name"], "didip")

    def test_set_non_json_data_from_stdin(self):
        try:
            s = StringIO.StringIO('awesome')
            sys.stdin = s
            b = base.Base()
            sys.stdin = sys.__stdin__
        except SystemExit:
            self.assertTrue("system exited because input from stdin is not in JSON format")

    def test_run_and_check_output(self):
        output = json.loads(self.logic.run())
        self.assertEqual(output["message"], "Success")


if __name__ == '__main__':
    unittest.main()
