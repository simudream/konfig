ResourceD new title: Toolkit of a Happy DevOps

Config management.

/roles
    prod-docker.toml
        steps = [
            "stacks/users.toml",
            "stacks/docker.toml"
        ]

        [matchers]
        hostname = ["=", "$HOSTNAME"]
        tags = ["aaa", "bbb"]


    staging-docker.toml

/stacks
    users.toml
        steps = ["/logic/users"]
    docker.toml
        steps = [
            "/logic/docker"
        ]

/logic
    base.py
    base.rb
    /users
        /data
        /templates
        /users.py > base.py

    /docker
        /data
        /templates
        /docker.py > base.py
