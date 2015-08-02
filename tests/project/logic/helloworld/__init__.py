from __future__ import absolute_import
from __future__ import division
from __future__ import print_function
from __future__ import unicode_literals
from __future__ import with_statement

import os.path
import sys
import json
sys.path.append(os.path.join(os.path.dirname(os.path.realpath(__file__)), '..'))

import base


class HelloWorld(base.Base):
    def dryrun(self):
        output = {
            'users': self.data.get('users', []),
            'template_name': 'helloworld.txt.tmpl',
            'target_path': '/tmp/helloworld.txt',
            'message': 'Hello world'
        }
        output_json = json.dumps(output)
        print(output_json)
        return output_json

    def run(self):
        self.write_file('helloworld.txt.tmpl', '/tmp/helloworld.txt', users=self.data.get('users', []))
        return self.dryrun()


if __name__ == "__main__":
    HelloWorld().run()
