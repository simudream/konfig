from __future__ import absolute_import
from __future__ import division
from __future__ import print_function
from __future__ import unicode_literals
from __future__ import with_statement

import os
import os.path
import sys
import argparse
import json
import socket
import jinja2
import subprocess


class Base(object):
    def __init__(self):
        self._setup_argsparser()
        self.hostname = socket.gethostname()
        self._read_data()
        self._setup_template()

    def _setup_argsparser(self):
        self.argsparser = argparse.ArgumentParser(description='Logic runner for {0}'.format(self.__class__.__name__))
        self.argsparser.add_argument('--dryrun', dest='dryrun', action='store_true')
        self.argsparser.add_argument('--no-dryrun', dest='dryrun', action='store_false')
        self.argsparser.set_defaults(dryrun=True)
        self.args = self.argsparser.parse_args()

    def _setup_template(self):
        templates_dir = os.path.join(self.current_dir(), 'templates')
        if os.path.isdir(templates_dir):
            self.template = jinja2.Environment(loader=jinja2.FileSystemLoader(templates_dir), trim_blocks=True)

    def _read_data(self):
        self.data = {}
        self._read_stdin()
        self._read_data_dir()

    def _read_stdin(self):
        if not sys.stdin.isatty(): # Avoid blocking on empty stdin
            try:
                input_data = sys.stdin.read()
                if input_data:
                    self.data = json.loads(input_data) or {}
            except ValueError:
                print('{"error": "Input from stdin is not in JSON format"}')
                sys.exit(1)

    def _read_data_dir(self):
        data_dir = os.path.join(self.current_dir(), 'data')
        if not os.path.exists(data_dir):
            return

        for filename in os.listdir(data_dir):
            if filename.endswith('.json'):
                key = filename.replace('.json', '')

                full_filename = os.path.join(data_dir, filename)
                with open(full_filename) as f:
                    new_data = json.loads(f.read())

                    if key in self.data:
                        # Merge or append new items to existing list
                        if isinstance(self.data[key], list):
                            if isinstance(new_data, list):
                                self.data[key].extend(new_data)
                            else:
                                self.data[key].append(new_data)

                    else:
                        self.data[key] = new_data

    def current_dir(self):
        return os.path.dirname(os.path.realpath(sys.modules[self.__module__].__file__))

    def write_file(self, template_filename, target_path, **kwargs):
        if self.template is None:
            print('{"error": "jinja2 Environment object is missing. Templates subdirectory must exist inside your logic directory."}')
            sys.exit(1)

        file_content = self.template.get_template(template_filename).render(**kwargs)
        with open(target_path, "wb") as fh:
            fh.write(file_content)

    def exec_with_dryrun(self, command):
        if self.args.dryrun:
            return command, 0
        else:
            p = subprocess.Popen(command, stdout=subprocess.PIPE, shell=True)
            (output, err) = p.communicate()
            exit_code = p.wait()
            return output, exit_code

    def dryrun(self):
        output = '{"message": "Success"}'
        print(output)
        return output

    def run(self):
        return self.dryrun()


if __name__ == "__main__":
    Base().run()
