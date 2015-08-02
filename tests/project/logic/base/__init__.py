from __future__ import absolute_import
from __future__ import division
from __future__ import print_function
from __future__ import unicode_literals
from __future__ import with_statement

import os
import os.path
import sys
import json
import socket
import jinja2


class Base(object):
    def __init__(self):
        self.dryrun_flag = True
        self.hostname = socket.gethostname()
        self.data = {}
        self._read_data()
        self._setup_template()

    def _setup_template(self):
        templates_dir = os.path.join(self.current_dir(), 'templates')
        if os.path.isdir(templates_dir):
            self.template = jinja2.Environment(loader=jinja2.FileSystemLoader(templates_dir), trim_blocks=True)

    def _read_data(self):
        self._read_stdin()
        self._read_data_dir()

    def _read_stdin(self):
        try:
            self.data = json.loads(sys.stdin.read())
        except ValueError:
            print('{"error": "Input from stdin is not in JSON format"}')
            sys.exit(1)

    def _read_data_dir(self):
        data_dir = os.path.join(self.current_dir(), 'data')
        if not os.path.exists(data_dir):
            return

        for filename in os.listdir(data_dir):
            if filename.lower().startswith('readme'):
                continue

            full_filename = os.path.join(data_dir, filename)
            with open(full_filename) as f:
                if full_filename.endswith('.json'):
                    for key, value in json.loads(f.read()).items():
                        self.data[key] = value

    def current_dir(self):
        return os.path.dirname(os.path.realpath(sys.modules[self.__module__].__file__))

    def write_file(self, template_filename, target_path, **kwargs):
        if self.template is None:
            print('{"error": "jinja2 Environment object is missing. Templates subdirectory must exist inside your logic directory."}')
            sys.exit(1)

        file_content = self.template.get_template(template_filename).render(**kwargs)
        with open(target_path, "wb") as fh:
            fh.write(file_content)

    def exec_or_print(self, command):
        if self.dryrun_flag:
            print(command)
            return 0
        else:
            return os.system(command)

    def dryrun(self):
        output = '{"message": "Success"}'
        print(output)
        return output

    def run(self):
        return self.dryrun()
