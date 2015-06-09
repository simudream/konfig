// Package engine provides engine struct.
package engine

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"regexp"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/resourced/configurator/role"
	"github.com/resourced/configurator/stack"
)

// New is the constructor for a new engine.
func New(root string) (*Engine, error) {
	engine := &Engine{Root: root}

	logicDirs, err := engine.readDir("logic")
	if err != nil {
		return nil, err
	}

	stackFiles, err := engine.readDir("stacks")
	if err != nil {
		return nil, err
	}

	roleFiles, err := engine.readDir("roles")
	if err != nil {
		return nil, err
	}

	hostname, err := os.Hostname()
	if err != nil {
		return nil, err
	}

	engine.Logic = logicDirs
	engine.Stacks = stackFiles
	engine.Roles = roleFiles

	engine.DryRun = true
	engine.PythonPath = "/usr/bin/python"
	engine.PipPath = "/usr/local/bin/pip"
	engine.RubyPath = "/usr/bin/ruby"
	engine.BundlePath = "bundle"
	engine.Hostname = hostname

	return engine, nil
}

type Engine struct {
	// Root is the root of project directory.
	Root string

	// PythonPath is the path to python executable.
	PythonPath string

	// PipPath is the path to pip executable.
	PipPath string

	// RubyPath is the path to ruby executable.
	RubyPath string

	// BundlePath is the path to bundle executable.
	BundlePath string

	// DryRun is the dry run flag, default is true.
	DryRun bool

	// Hostname is the host's name.
	Hostname string

	Logic  []os.FileInfo
	Stacks []os.FileInfo
	Roles  []os.FileInfo
}

func (e *Engine) readDir(dirname string) ([]os.FileInfo, error) {
	return ioutil.ReadDir(path.Join(e.Root, dirname))
}

// RunLogic allows engine to execute one logic layer.
func (e *Engine) RunLogic(name string) ([]byte, error) {
	pythonExecPath := path.Join(e.Root, "logic", name, "__init__.py")
	_, pyErr := os.Stat(pythonExecPath)
	if pyErr == nil {
		return e.RunPythonLogic(name)
	}

	rubyExecPath := path.Join(e.Root, "logic", name, name+".rb")
	_, rbErr := os.Stat(rubyExecPath)
	if rbErr == nil {
		return e.RunRubyLogic(name)
	}

	if os.IsNotExist(pyErr) || os.IsNotExist(rbErr) {
		return nil, errors.New(fmt.Sprintf("Logic must be implemented in Python(%v/__init__.py) or Ruby(%v/%v.rb)", name, name, name))
	}

	return nil, nil
}

// InstallPythonLogicDependencies allows engine to installs dependencies for a logic written in python.
func (e *Engine) InstallPythonLogicDependencies(name string) ([]byte, error) {
	reqPath := path.Join(e.Root, "logic", name, "requirements.txt")
	if e.DryRun {
		return []byte(e.PipPath + " install -r " + reqPath), nil
	}

	return exec.Command(e.PipPath, "install", "-r", reqPath).CombinedOutput()
}

// RunPythonLogic allows engine to run a logic written in python.
func (e *Engine) RunPythonLogic(name string) ([]byte, error) {
	execPath := path.Join(e.Root, "logic", name, "__init__.py")
	if e.DryRun {
		return []byte(e.PythonPath + " " + execPath), nil
	}

	return exec.Command(e.PythonPath, execPath).CombinedOutput()
}

// InstallRubyLogicDependencies allows engine to installs dependencies for a logic written in ruby.
func (e *Engine) InstallRubyLogicDependencies(name string) ([]byte, error) {
	logicPath := path.Join(e.Root, "logic", name)
	if e.DryRun {
		return []byte("cd " + logicPath + "; " + e.BundlePath), nil
	}

	cmd := exec.Command(e.BundlePath)
	cmd.Path = logicPath

	return cmd.CombinedOutput()
}

// RunRubyLogic allows engine to run a logic written in ruby.
func (e *Engine) RunRubyLogic(name string) ([]byte, error) {
	execPath := path.Join(e.Root, "logic", name, name+".rb")
	if e.DryRun {
		return []byte(e.RubyPath + " " + execPath), nil
	}

	return exec.Command(e.RubyPath, execPath).CombinedOutput()
}

// ReadStack allows engine to read a particular stack defined in TOML file.
func (e *Engine) ReadStack(name string) (stack.Stack, error) {
	var stk stack.Stack

	if !strings.HasSuffix(name, ".toml") {
		name = name + ".toml"
	}

	stackPath := path.Join(e.Root, "stacks", name)
	if _, err := toml.DecodeFile(stackPath, &stk); err != nil {
		return stk, err
	}

	return stk, nil
}

// RunStack allows engine to run a particular stack.
func (e *Engine) RunStack(name string) ([]byte, error) {
	stk, err := e.ReadStack(name)
	if err != nil {
		return nil, err
	}

	allOutput := make([]byte, 0)

	for _, step := range stk.Steps {
		if strings.HasPrefix(step, "logic/") {
			logicName := strings.Replace(step, "logic/", "", -1)

			output, err := e.RunLogic(logicName)
			if err != nil {
				return output, err
			}

			allOutput = append(allOutput, output...)
			allOutput = append(allOutput, []byte("\n")...)
		}
	}

	return allOutput, nil
}

// ReadRole allows engine to read a particular role defined in TOML file.
func (e *Engine) ReadRole(name string) (role.Role, error) {
	var rl role.Role

	if !strings.HasSuffix(name, ".toml") {
		name = name + ".toml"
	}

	rolePath := path.Join(e.Root, "roles", name)
	if _, err := toml.DecodeFile(rolePath, &rl); err != nil {
		return rl, err
	}

	return rl, nil
}

func (e *Engine) checkIfRoleMatched(rl role.Role) bool {
	hostnameMatcherIsProvided := true
	if rl.Matchers.Hostname == nil || len(rl.Matchers.Hostname) < 2 {
		hostnameMatcherIsProvided = false
	}

	tagsMatcherIsProvided := true
	if rl.Matchers.Tags == nil || len(rl.Matchers.Tags) < 1 {
		tagsMatcherIsProvided = false
	}

	// There are no matchers provided.
	if !hostnameMatcherIsProvided && !tagsMatcherIsProvided {
		return false
	}

	if hostnameMatcherIsProvided {
		return e.checkIfRoleMatchedByHostname(rl)
	} else if tagsMatcherIsProvided {
		return e.checkIfRoleMatchedByTags(rl)
	}

	return false
}

func (e *Engine) checkIfRoleMatchedByHostname(rl role.Role) bool {
	hostnameMatcherOperator := rl.Matchers.Hostname[0]
	hostnameMatcherValue := rl.Matchers.Hostname[1]

	// If value == "$HOSTNAME", substitute it to e.Hostname
	if hostnameMatcherValue == "$HOSTNAME" {
		hostnameMatcherValue = e.Hostname
	}

	if hostnameMatcherOperator == "=" {
		return e.Hostname == os.ExpandEnv(hostnameMatcherValue)

	} else if hostnameMatcherOperator == "~" {
		reg, err := regexp.Compile(hostnameMatcherValue)
		if err != nil {
			return false
		}

		return reg.MatchString(e.Hostname)
	}

	return false
}

func (e *Engine) checkIfRoleMatchedByTags(rl role.Role) bool {
	// TODO(didip): Fetch tags data first.
	return false
}

// RunRole allows engine to run a particular role.
func (e *Engine) RunRole(name string) ([]byte, error) {
	rl, err := e.ReadRole(name)
	if err != nil {
		return nil, err
	}

	allOutput := make([]byte, 0)

	if e.checkIfRoleMatched(rl) {
		for _, step := range rl.Steps {
			if strings.HasPrefix(step, "stacks/") {
				stackName := strings.Replace(step, "stacks/", "", -1)

				output, err := e.RunStack(stackName)
				if err != nil {
					return output, err
				}

				allOutput = append(allOutput, output...)
				allOutput = append(allOutput, []byte("\n")...)
			}

			if strings.HasPrefix(step, "logic/") {
				logicName := strings.Replace(step, "logic/", "", -1)

				output, err := e.RunLogic(logicName)
				if err != nil {
					return output, err
				}

				allOutput = append(allOutput, output...)
				allOutput = append(allOutput, []byte("\n")...)
			}
		}
	}

	return allOutput, nil
}

func (e *Engine) RunRoles() ([]byte, error) {
	allOutput := make([]byte, 0)

	for _, roleFile := range e.Roles {
		// roleFile must always be a file.
		if roleFile.IsDir() {
			continue
		}

		output, err := e.RunRole(roleFile.Name())
		if err != nil {
			return output, err
		}

		allOutput = append(allOutput, output...)
		allOutput = append(allOutput, []byte("\n")...)
	}

	return allOutput, nil
}
