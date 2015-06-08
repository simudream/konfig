package role

import (
	"github.com/resourced/configurator/stack"
)

type Role struct {
	stack.Stack

	Matchers struct {
		Hostname []string `toml:"hostname"`
		Tags     []string `toml:"tags"`
	}
}
