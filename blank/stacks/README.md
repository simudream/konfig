Stacks allow you to mix and match multiple stacks or logic.

### To create a new stack

* Create a `.toml` file under `/stacks`. Here's an example:
    ```
    # helloworld.toml

    # steps are series of stacks or logic.
    # When there's an error, Cofigurator stops running the entire steps.
    steps = [
        "stacks/default.toml",
        "logic/sysctl",
    ]
    ```
