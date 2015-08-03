[![GoDoc](https://godoc.org/github.com/resourced/resourced-stacks?status.svg)](http://godoc.org/github.com/resourced/resourced-stacks)
[![license](http://img.shields.io/badge/license-MIT-red.svg?style=flat)](https://raw.githubusercontent.com/resourced/resourced-stacks/master/LICENSE)

**ResourceD Stacks:**

Task runner, file stamper, quite possible the simplest configuration management you've ever seen.


## Comparison to existing solutions:

* We are very serious about simple workflow. Logic->Stack, that's all.

    * **Logic:** Write a python script to accomplish a particular task. Use JSON format for both stdin and stdout.

    * **Stack:** This is how you can create a mixins of logic or other stacks.

* **IT DOES NOT USE DSL.**

* Dry run switch is `true` by default.

* It uses git to store and version metadata.

* It understands EC2 tags and use them for role matching.


## Example when running locally
```
resourced-stacks -root=./tests/project -cmd=run -stack=helloworld -dryrun=false
INFO[0000] Starting stack: helloworld                                               dryrun=false
INFO[0000] Starting logic: helloworld                                               dryrun=false
INFO[0001] Executed: pip install -r tests/project/logic/helloworld/requirements.txt dryrun=false
INFO[0001] Executed: python tests/project/logic/helloworld/__init__.py --no-dryrun  dryrun=false
```


## Prerequisites

* Git

* Python and Pip, or

* Ruby and Bundler.


## CLI mode

Download the binary release [here](https://github.com/resourced/resourced-stacks/releases) and starts using it.
```
resourced-stacks -h

# Creating a new project
resourced-stacks -root=/path/to/project -cmd=new

# Cleaning dirty local changes on remote host
resourced-stacks -root=/path/to/project -cmd=clean

# Running on a host, by default -cmd=run
resourced-stacks -root=/path/to/project -stack=stack-name
```


## Agent mode

1. Download ResourceD binary release [here](https://github.com/resourced/resourced/releases).

2. Configure ResourceD Executor to run your stacks.

3. Run ResourceD.
