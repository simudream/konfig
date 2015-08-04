Stacks allow you to mix and match multiple stacks or logic.

### To create a new stack

1. Create a subdirectory inside `stacks` directory and give it the name of your stack.

2. In there, create a `.toml` file with the same name.

For example, given that the stack name is helloworld,
```
# stacks/helloworld/helloworld.toml

# steps are series of stacks or logic.
# When there's an error, the entire flow is stopped.
steps = [
    "stacks/otherstack",
    "logic/sysctl",
    "logic/iptables"
]
```

### Passing data to your stack

Each step may require metadata to complete (e.g. metadata required for stamping files).

To pass data to your stack:

1. Create `data` subdirectory under `stacks` directory.

2. In there, create `.json` file containing your data. We recommend creating & naming the JSON file per category that makes sense to you.

For a complete example, take a look at `tests/project` directory.
