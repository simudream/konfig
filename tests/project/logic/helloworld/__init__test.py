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

import helloworld

class TestHelloWorld(unittest.TestCase):
    def setUp(self):
        s = StringIO.StringIO('{"users": [{"name": "didip"}, {"name": "brotato"}], "motd": {"message": "this will override data defined in logic"}}')
        sys.stdin = s
        self.logic = helloworld.HelloWorld()
        sys.stdin = sys.__stdin__

    def test_read_list_data_from_json_file_and_stdin(self):
        self.assertEqual(len(self.logic.data["users"]), 3)

        found_brad = False
        found_didip = False
        found_brotato = False

        for user in self.logic.data["users"]:
            if user['name'] == 'brad':
                found_brad = True
            if user['name'] == 'didip':
                found_didip = True
            if user['name'] == 'brotato':
                found_brotato = True

        self.assertTrue(found_brad, 'brad is expected to be found')
        self.assertTrue(found_didip, 'didip is expected to be found')
        self.assertTrue(found_brotato, 'brotato is expected to be found')

    def test_read_dict_data_from_json_file_and_stdin(self):
        self.assertTrue('motd' in self.logic.data)
        self.assertEqual(self.logic.data['motd']['message'], 'this will override data defined in logic')

        without_stdin_logic = helloworld.HelloWorld()
        self.assertTrue('motd' in without_stdin_logic.data)
        self.assertEqual(without_stdin_logic.data['motd']['message'], 'Awesome Sauce')


    def test_write_file_from_template(self):
        self.logic.write_file('helloworld.txt.tmpl', '/tmp/helloworld.txt', users=self.logic.data['users'])

        self.assertTrue(os.path.exists('/tmp/helloworld.txt'))
        os.remove('/tmp/helloworld.txt')


if __name__ == '__main__':
    unittest.main()
