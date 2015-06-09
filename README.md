[![GoDoc](https://godoc.org/github.com/resourced/configurator?status.svg)](http://godoc.org/github.com/resourced/configurator)
[![license](http://img.shields.io/badge/license-MIT-red.svg?style=flat)](https://raw.githubusercontent.com/resourced/configurator/master/LICENSE)

**ResourceD Configurator:** The simplest configuration management in the world.


## Comparison to existing solutions:

* We are very serious about simple workflow. Logic->Stack->Role, that's all.

    * **Logic:** This is where you implement the management process. You can use Ruby or Python.

    * **Stack:** This is how you can create a mixins of logic or other stacks.

    * **Role:** This is how a particular host is matched to stacks. Think of `site.pp`

* **IT DOES NOT USE DSL.**

* Dry run switch is `true` by default.

* It uses git to store and version metadata.

* It has two mode: Agent and SSH, just like others out there, but here's the difference:

    * When in Agent mode, there's only one binary to download and install.

    * When in SSH mode, there's only one binary to download and install.

* It understands EC2 tags and use them for role matching.


## Example when running locally
```
configurator -root=./tests/project -cmd=run
INFO[0000] Role helloworld-staging.toml matched.                              dryrun=true hostname=didip-mac-mini.local matcher:=[= $HOSTNAME]
INFO[0000] Running role: helloworld-staging.toml                              dryrun=true
INFO[0000] Running stack: helloworld.toml                                     dryrun=true
INFO[0000] Running logic: helloworld-py                                       dryrun=true
INFO[0000] pip install -r tests/project/logic/helloworld-py/requirements.txt  dryrun=true
INFO[0000] python tests/project/logic/helloworld-py/__init__.py               dryrun=true
INFO[0000] Running logic: helloworld-rb                                       dryrun=true
INFO[0000] cd tests/project/logic/helloworld-rb && bundle                     dryrun=true
INFO[0000] ruby tests/project/logic/helloworld-rb/helloworld-rb.rb            dryrun=true
```


## Prerequisites

* Git

* Python and Pip, or

* Ruby and Bundler.


## SSH mode installation

Download the binary release [here](https://github.com/resourced/configurator/releases) and starts using it.
```
configurator -h

# Creating a new project
configurator -root=/path/to/project -cmd=new

# Cleaning dirty local changes on remote host
configurator -root=/path/to/project -cmd=clean

# Running on remote host
configurator -root=/path/to/project
```


## Agent mode installation

1. Download ResourceD binary release [here](https://github.com/resourced/resourced/releases).

2. Configure ResourceD to run the Configurator.
