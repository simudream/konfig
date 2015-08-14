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


class User(base.Base):
    def _new_user_cmd(self, name, primary_group=None):
        if primary_group:
            return "useradd -g {0} {1}".format(primary_group, name)
        else:
            return "useradd {0}".format(name)

    def _add_user_to_groups_cmd(self, name, groups):
        return "usermod -a -G {0} {1}".format(','.join(groups), name)

    def run(self):
        report = {'success': [], 'failure': []}

        for user in self.data["users"]:
            if 'shell' not in user:
                user['shell'] = '/bin/bash'

            if 'groups' in user and len(user['groups']) > 0:
                # Create user and assign primary group
                output, exit_code = self.exec_with_dryrun(self._new_user_cmd(user['name'], user['groups'][0]))

                if exit_code == 0:
                    report['success'].append(user)
                else:
                    report['failure'].append(user)

                # Add user to secondary groups
                if len(user['groups']) > 1:
                    output, exit_code = self.exec_with_dryrun(self._add_user_to_groups_cmd(user['name'], user['groups'][1:]))

                if exit_code == 0:
                    report['success'].append(user)
                else:
                    report['failure'].append(user)

            else:
                # Create user only
                output, exit_code = self.exec_with_dryrun(self._new_user_cmd(user['name']))
                if exit_code == 0:
                    report['success'].append(user)
                else:
                    report['failure'].append(user)

        return report


if __name__ == "__main__":
    print(json.dumps(User().run()))
