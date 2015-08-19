[![GoDoc](https://godoc.org/github.com/didip/konfig?status.svg)](http://godoc.org/github.com/didip/konfig)
[![license](http://img.shields.io/badge/license-MIT-red.svg?style=flat)](https://raw.githubusercontent.com/didip/konfig/master/LICENSE)

**Konfig:** Task runner, file stamper, and more...

It is quite possible the simplest configuration management you've ever seen.


## Comparison to existing solutions:

* We are very serious about simple workflow. Logic->Stack, that's all.

    * **Logic:** Write a python script to accomplish a particular task. Use JSON for stdin & stdout.

    * **Stack:** This is how you can create a mixins of logic and other stacks in sequential order.

* **IT DOES NOT USE DSL**

* Dry run switch is `true` by default.

* Built-in testing story. Every logic **MUST** contain test files.

* It uses git for versioning and deploying on a remote host.


## Example when running locally
```
konfig -root=./tests/project -cmd=run -stack=helloworld -dryrun=false
INFO[0000] Starting stack: helloworld                                                dryrun=false
INFO[0000] Reading data for stack: helloworld                                        dryrun=false
INFO[0000] Running tests for logic: helloworld                                       dryrun=false
INFO[0000] Installing dependencies for logic: helloworld                             dryrun=false
INFO[0002] Executed: pip install -r tests/project/logic/helloworld/requirements.txt  dryrun=false
INFO[0002] Checking if test files exist for logic: helloworld                        dryrun=false
INFO[0002] Executed: python tests/project/logic/helloworld/__init__test.py           dryrun=false
INFO[0002] Starting logic: helloworld                                                dryrun=false
INFO[0002] Installing dependencies for logic: helloworld                             dryrun=false
INFO[0002] Executed: pip install -r tests/project/logic/helloworld/requirements.txt  dryrun=false
INFO[0002] Executed: python tests/project/logic/helloworld/__init__.py --no-dryrun   dryrun=false
```


## Prerequisites

* Git for versioning and deploying on a remote host.

* Python and Pip for implementing logic.


## Installation

Download the binary release [here](https://github.com/didip/konfig/releases) and starts using it.


## Usage
```
konfig -h

# Creating a new project
konfig -root=/path/to/project -cmd=new

# Running on a host
konfig -root=/path/to/project -cmd=run -stack=stack-name -dryrun=false

# Pulling down remote git project onto -root
konfig -root=/path/to/project -cmd=pull -git=https://github.com/path/to/project/repo.git

# Pulling down remote git project onto -root & then running it
# You can run this command under minutely cron to replicate `puppet agent -t` behavior
konfig -root=/path/to/project -cmd=pull-run -git=https://github.com/path/to/project/repo.git -stack=stack-name -dryrun=false
```
