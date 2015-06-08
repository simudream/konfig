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

func (e *Engine) InstallPythonLogicDependencies(name string) ([]byte, error) {
	reqPath := path.Join(e.Root, "logic", name, "requirements.txt")
	if e.DryRun {
		return []byte(e.PipPath + " install -r " + reqPath), nil
	}

	return exec.Command(e.PipPath, "install", "-r", reqPath).CombinedOutput()
}

func (e *Engine) RunLogic(name string) ([]byte, error) {
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
