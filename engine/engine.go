package engine

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/resourced/configurator/stack"
)

func New(root string) (*Engine, error) {
	engine := &Engine{Root: root}

	logicDirs, err := engine.ReadDir("logic")
	if err != nil {
		return nil, err
	}

	stackFiles, err := engine.ReadDir("stacks")
	if err != nil {
		return nil, err
	}

	roleFiles, err := engine.ReadDir("roles")
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

	return engine, nil
}

type Engine struct {
	Root       string
	PythonPath string
	PipPath    string
	RubyPath   string
	BundlePath string
	DryRun     bool
	Logic      []os.FileInfo
	Stacks     []os.FileInfo
	Roles      []os.FileInfo
}

func (e *Engine) ReadDir(dirname string) ([]os.FileInfo, error) {
	return ioutil.ReadDir(path.Join(e.Root, dirname))
}

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

func (e *Engine) InstallPythonLogicDependencies(name string) ([]byte, error) {
	reqPath := path.Join(e.Root, "logic", name, "requirements.txt")
	if e.DryRun {
		return []byte(e.PipPath + " install -r " + reqPath), nil
	}

	return exec.Command(e.PipPath, "install", "-r", reqPath).CombinedOutput()
}

func (e *Engine) RunPythonLogic(name string) ([]byte, error) {
	execPath := path.Join(e.Root, "logic", name, "__init__.py")
	if e.DryRun {
		return []byte(e.PythonPath + " " + execPath), nil
	}

	return exec.Command(e.PythonPath, execPath).CombinedOutput()
}

func (e *Engine) InstallRubyLogicDependencies(name string) ([]byte, error) {
	logicPath := path.Join(e.Root, "logic", name)
	if e.DryRun {
		return []byte("cd " + logicPath + "; " + e.BundlePath), nil
	}

	cmd := exec.Command(e.BundlePath)
	cmd.Path = logicPath

	return cmd.CombinedOutput()
}

func (e *Engine) RunRubyLogic(name string) ([]byte, error) {
	execPath := path.Join(e.Root, "logic", name, name+".rb")
	if e.DryRun {
		return []byte(e.RubyPath + " " + execPath), nil
	}

	return exec.Command(e.RubyPath, execPath).CombinedOutput()
}

func (e *Engine) ReadStack(name string) (stack.Stack, error) {
	var stk stack.Stack

	stackPath := path.Join(e.Root, "stacks", name+".toml")
	if _, err := toml.DecodeFile(stackPath, &stk); err != nil {
		return stk, err
	}

	return stk, nil
}

func (e *Engine) RunStack(name string) ([]byte, error) {
	stk, err := e.ReadStack(name)
	if err != nil {
		return nil, err
	}

	for _, step := range stk.Steps {
		if strings.HasPrefix(step, "logic/") {
			logicName := strings.Replace(step, "logic/", "", -1)

			output, err := e.RunLogic(logicName)
			log.Printf(string(output))

			if err != nil {
				return output, err
			}
		}
	}

	return nil, nil
}
