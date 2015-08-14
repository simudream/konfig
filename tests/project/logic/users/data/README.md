Put all users metadata here.

### JSON format
{
    "name": "",
    "uid": 1,
    "groups": [200, "developers", "designers"],
    "pub_key": "",
    "shell": "/bin/bash"
}

* name: The name of user.

* uid: The user id.

* groups: The groups user belong to. The first one is primary, everything else is secondary.

* pub_key: The public key of that user.

* shell: The shell of that user. Default is `/bin/bash`.