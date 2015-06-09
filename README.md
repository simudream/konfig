[![GoDoc](https://godoc.org/github.com/resourced/configurator?status.svg)](http://godoc.org/github.com/resourced/configurator)
[![license](http://img.shields.io/badge/license-MIT-red.svg?style=flat)](https://raw.githubusercontent.com/resourced/configurator/master/LICENSE)

**ResourceD Configurator:** The simplest configuration management in the world.

## Strengths compared to existing solutions:

* It simplify your workflow by a lot. All the scripting logic is done in Python or Ruby.

* We are very serious about simple workflow. Logic -> Stack -> Role, that's all you have to understand.

    * **Logic:** This is where you implement the installation process. You can use Ruby or Python.

    * **Stacks:** This is where you create a mixins of logic or other stacks.

    * **Role:** This is how configurator understand how to apply stacks to a particular host.

* It uses git to store and version metadata.

* Similar to existing solutions out there, it has two mode: Agent and SSH. But here's the difference:

    * When in Agent mode, there's only one binary to download and install.

    * When in SSH mode, there's only one binary to download and install as well.

* It understands EC2 tags and use them for role matching.


## SSH mode Installation

Download the binary release [here](https://github.com/resourced/configurator/releases) and starts using it.
```
configurator -h
```

## Agent mode Installation

1. Download ResourceD binary release [here](https://github.com/resourced/resourced/releases).

2. Configure ResourceD to run the Configurator.
