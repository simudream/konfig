Role is where the rubber meets the road so to speak. Configurator loops through list of all roles and apply the first matching role to the target host.

### To create a new role

* Create a `.toml` file under `/roles`. Here's an example:
    ```
    # helloworld-staging.toml

    # steps are series of stacks or logic this role will apply to target host.
    # You can have multiple stacks or logic.
    steps = [
        "stacks/helloworld.toml",
        "logic/base",
    ]

    # matchers is how you tie-in role to a specific host.
    # There are two types of matchers: hostname or tags.
    [matchers]
    # There are two kind of operator you can use for hostname: = or ~
    # hostname = ["=", "$HOSTNAME"] basically means apply this to every target host.
    hostname = ["=", "$HOSTNAME"]

    # Currently matches only EC2 tags
    tags = ["key:value", "key2:value2"]
    ```
