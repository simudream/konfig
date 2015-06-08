from __future__ import absolute_import
from __future__ import division
from __future__ import print_function
from __future__ import unicode_literals
from __future__ import with_statement

import os.path
import sys
import pwd
sys.path.append(os.path.join(os.path.dirname(os.path.realpath(__file__)), '..'))

import base


class User(base.Base):
    def _new_user_cmd(self, name, primary_group=None):
        if primary_group:
            return "useradd -g {0} {1}".format(primary_group, name)
        else:
            return "useradd {0}".format(name)

    def _update_user_primary_group_cmd(self, name, primary_group):
        return "usermod -g {0} {1}".format(primary_group, name)

    def _add_user_to_groups_cmd(self, name, groups):
        return "usermod -a -G {0} {1}".format(','.join(groups), name)

    def run(self):
        for _, data in self.data.items():
            try:
                pwd.getpwnam(data['name'])

                if 'groups' in data and len(data['groups']) > 0:
                    # Update user primary group
                    self.exec_or_print(self._update_user_primary_group_cmd(data['name'], data['groups'][0]))

            except KeyError:
                # User does not exist.
                if 'groups' in data and len(data['groups']) > 0:
                    # Create user and assign primary group
                    self.exec_or_print(self._new_user_cmd(data['name'], data['groups'][0]))

                    # Add user to secondary groups
                    if len(data['groups']) > 1:
                        self.exec_or_print(self._add_user_to_groups_cmd(data['name'], data['groups'][1:]))

                else:
                    # Create user only
                    self.exec_or_print(self._new_user_cmd(data['name']))


if __name__ == "__main__":
    User().run()
