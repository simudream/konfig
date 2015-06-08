from __future__ import absolute_import
from __future__ import division
from __future__ import print_function
from __future__ import unicode_literals
from __future__ import with_statement

import os
import os.path
import sys
import json


class Base(object):
    def __init__(self):
        self.dry_run = True
        self.data = {}
        self._read_data()

    def _read_data(self):
        data_dir = os.path.join(self.current_dir(), 'data')
        if not os.path.exists(data_dir):
            return

        for filename in os.listdir(data_dir):
            if filename.lower().startswith('readme'):
                continue

            full_filename = os.path.join(data_dir, filename)
            with open(full_filename) as f:
                self.data[full_filename] = f.read()

                if full_filename.endswith('.json'):
                    self.data[full_filename] = json.loads(self.data[full_filename])

    def current_dir(self):
        return os.path.dirname(os.path.realpath(sys.modules[self.__module__].__file__))

    def exec_or_print(self, command):
        if self.dry_run:
            print(command)
            return 0
        else:
            return os.system(command)

    def init(self):
        pass

    def run(self):
        pass
