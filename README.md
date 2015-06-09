[![GoDoc](https://godoc.org/github.com/resourced/configurator?status.svg)](http://godoc.org/github.com/resourced/configurator)
[![license](http://img.shields.io/badge/license-MIT-red.svg?style=flat)](https://raw.githubusercontent.com/resourced/configurator/master/LICENSE)

**ResourceD Configurator:** The simplest configuration management in the world.

## Comparison to existing solutions:

* We are very serious about simple workflow. Logic->Stack->Role, that's all.

    * **Logic:** This is where you implement the management process. You can use Ruby or Python.

    * **Stack:** This is how you can create a mixins of logic or other stacks.

    * **Role:** This is how a particular host is matched to stacks. Think of `site.pp`

* IT DOES NOT USE DSL.

* It uses git to store and version metadata.

* It has two mode: Agent and SSH, just like others out there, but here's the difference:

    * When in Agent mode, there's only one binary to download and install.

    * When in SSH mode, there's only one binary to download and install.

* It understands EC2 tags and use them for role matching.


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
