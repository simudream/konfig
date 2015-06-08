from __future__ import absolute_import
from __future__ import division
from __future__ import print_function
from __future__ import unicode_literals
from __future__ import with_statement

import os.path
import sys
sys.path.append(os.path.join(os.path.dirname(os.path.realpath(__file__)), '..'))

import base


class HelloWorld(base.Base):
    def run(self):
        print("Hello World")


if __name__ == "__main__":
    HelloWorld().run()
