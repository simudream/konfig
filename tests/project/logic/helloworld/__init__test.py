#!/usr/bin/env python

import sys
import StringIO
import unittest
import os.path
import json
sys.path.append(os.path.join(os.path.dirname(os.path.realpath(__file__)), '..'))

import helloworld

class TestHelloWorld(unittest.TestCase):
    def setUp(self):
        s = StringIO.StringIO('{"users": [{"name": "didip"}, {"name": "brotato"}]}')
        sys.stdin = s
        self.logic = helloworld.HelloWorld()
        sys.stdin = sys.__stdin__

    def test_write_file_from_template(self):
        self.logic.write_file('helloworld.txt.tmpl', '/tmp/helloworld.txt', users=self.logic.data['users'])

        self.assertTrue(os.path.exists('/tmp/helloworld.txt'))
        os.remove('/tmp/helloworld.txt')


if __name__ == '__main__':
    unittest.main()
