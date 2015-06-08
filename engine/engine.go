package engine

import (
	"io/ioutil"
	"os"
	"os/exec"
	"path"
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
	engine.RubyPath = "/usr/bin/ruby"

	return engine, nil
}

type Engine struct {
	Root       string
	PythonPath string
	RubyPath   string
	DryRun     bool
	Logic      []os.FileInfo
	Stacks     []os.FileInfo
	Roles      []os.FileInfo
}

func (e *Engine) ReadDir(dirname string) ([]os.FileInfo, error) {
	return ioutil.ReadDir(path.Join(e.Root, dirname))
}

func (e *Engine) RunLogic(name string) ([]byte, error) {
	logicPath := path.Join(e.Root, "logic", name, "__init__.py")
	if e.DryRun {
		return []byte(e.PythonPath + " " + logicPath), nil
	}

	return exec.Command(e.PythonPath, logicPath).CombinedOutput()
}
