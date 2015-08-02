// Package engine provides engine struct.
package engine

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/Sirupsen/logrus"
	"github.com/resourced/resourced-stacks/stack"
	"github.com/robertkrimen/otto"
)

// New is the constructor for a new engine.
func New(root, conditions string) (*Engine, error) {
	hostname, err := os.Hostname()
	if err != nil {
		return nil, err
	}

	engine := &Engine{Root: root, Hostname: hostname}

	err = os.MkdirAll(root, 0755)
	if err != nil {
		return nil, err
	}

	logicDirs, err := engine.readDir("logic")
	if err != nil {
		return nil, err
	}

	stackFiles, err := engine.readDir("stacks")
	if err != nil {
		return nil, err
	}

	engine.Logic = logicDirs
	engine.Stacks = stackFiles

	engine.DryRun = true
	engine.PythonPath = "python"
	engine.PipPath = "pip"
	engine.RubyPath = "ruby"
	engine.BundlePath = "bundle"
	engine.jsVM = otto.New()

	engine.SetConditions(conditions)

	return engine, nil
}

type Engine struct {
	Root string

	// PythonPath is the path to python executable.
	PythonPath string

	// PipPath is the path to pip executable.
	PipPath string

	// RubyPath is the path to ruby executable.
	RubyPath string

	// BundlePath is the path to bundle executable.
	BundlePath string

	// Conditions to match before running stacks/logic.
	Conditions string

	DryRun bool

	Hostname string

	EC2Tags []map[string]string

	GitBranch string

	Logic  []os.FileInfo
	Stacks []os.FileInfo

	jsVM *otto.Otto
}

// SetConditions format and assigns JS conditions.
func (e *Engine) SetConditions(conditions string) {
	if conditions == "" {
		conditions = "true"
	}

	e.Conditions = conditions
}

func (e *Engine) EvalConditions() (bool, error) {
	e.jsVM.Set("name", e.Hostname)
	e.jsVM.Set("tags", make(map[string]string))

	value, err := e.jsVM.Run(e.Conditions)
	if err != nil {
		return false, err
	}
	return value.ToBoolean()
}

// IsGitRepo checks if Root is a git repo.
func (e *Engine) IsGitRepo() bool {
	_, err := os.Stat(path.Join(e.Root, ".git"))
	if err != nil {
		return false
	}
	return true
}

func (e *Engine) CleanProject() error {
	if !e.IsGitRepo() {
		return nil
	}

	cmd := exec.Command("git", "reset", "--hard")
	cmd.Path = e.Root

	output, err := cmd.CombinedOutput()
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Fatal(string(output))

		return err
	}

	logrus.Info(string(output))
	return nil
}

func (e *Engine) NewProject() error {
	// 1. Create tmp directory.
	dir, err := ioutil.TempDir(os.TempDir(), "resourced-stacks")
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Fatal(err)
	}
	defer os.RemoveAll(dir)

	// 2. git clone to /tmp directory.
	output, err := exec.Command("git", "clone", "https://github.com/resourced/resourced-stacks.git", dir).CombinedOutput()
	if err != nil {
		os.RemoveAll(dir)

		logrus.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Fatal(string(output))
	}

	logrus.Info(string(output))

	// 3. mv blank template folder to Root
	logrus.Infof("Moving %v to %v...", path.Join(dir, "blank"), e.Root)
	err = os.Rename(path.Join(dir, "blank"), e.Root)
	if err != nil {
		os.RemoveAll(dir)

		if !strings.Contains(err.Error(), "directory not empty") {
			logrus.WithFields(logrus.Fields{
				"error": err.Error(),
			}).Fatal(string(output))
		}
	}
	return nil
}

func (e *Engine) readDir(dirname string) ([]os.FileInfo, error) {
	return ioutil.ReadDir(path.Join(e.Root, dirname))
}

// RunLogic allows engine to execute one logic layer.
func (e *Engine) RunLogic(name string) ([]byte, error) {
	logrus.WithFields(logrus.Fields{
		"dryrun": e.DryRun,
	}).Infof("Running logic: %v", name)

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
		err := errors.New(fmt.Sprintf("Logic must be implemented in Python(%v/__init__.py)", name))

		logrus.WithFields(logrus.Fields{
			"dryrun": e.DryRun,
			"error":  err.Error(),
		}).Errorf("Unable to run logic: %v", name)

		return nil, err
	}

	return nil, nil
}

// InstallPythonLogicDependencies allows engine to installs dependencies for a logic written in python.
func (e *Engine) InstallPythonLogicDependencies(name string) ([]byte, error) {
	reqPath := path.Join(e.Root, "logic", name, "requirements.txt")
	if e.DryRun {
		logrus.WithFields(logrus.Fields{
			"dryrun": e.DryRun,
		}).Infof("%v install -r %v", e.PipPath, reqPath)

		return nil, nil
	}

	_, err := os.Stat(reqPath)
	if err != nil {
		return make([]byte, 0), nil
	}

	return exec.Command(e.PipPath, "install", "-r", reqPath).CombinedOutput()
}

// RunPythonLogic allows engine to run a logic written in python.
func (e *Engine) RunPythonLogic(name string) ([]byte, error) {
	_, err := e.InstallPythonLogicDependencies(name)
	if err != nil {
		return nil, err
	}

	execPath := path.Join(e.Root, "logic", name, "__init__.py")
	if e.DryRun {
		logrus.WithFields(logrus.Fields{
			"dryrun": e.DryRun,
		}).Infof("%v %v", e.PythonPath, execPath)

		return nil, nil
	}

	return exec.Command(e.PythonPath, execPath).CombinedOutput()
}

// InstallRubyLogicDependencies allows engine to installs dependencies for a logic written in ruby.
func (e *Engine) InstallRubyLogicDependencies(name string) ([]byte, error) {
	installCommand := fmt.Sprintf("%v install --path vendor && %v package", e.BundlePath, e.BundlePath)

	logicPath := path.Join(e.Root, "logic", name)
	if e.DryRun {
		logrus.WithFields(logrus.Fields{
			"dryrun": e.DryRun,
		}).Infof("cd %v && %v", logicPath, installCommand)

		return nil, nil
	}

	_, err := os.Stat(logicPath)
	if err != nil {
		return make([]byte, 0), nil
	}

	cmd := exec.Command(installCommand)
	cmd.Path = logicPath

	return cmd.CombinedOutput()
}

// RunRubyLogic allows engine to run a logic written in ruby.
func (e *Engine) RunRubyLogic(name string) ([]byte, error) {
	_, err := e.InstallRubyLogicDependencies(name)
	if err != nil {
		return nil, err
	}

	execPath := path.Join(e.Root, "logic", name, name+".rb")
	if e.DryRun {
		logrus.WithFields(logrus.Fields{
			"dryrun": e.DryRun,
		}).Infof("%v %v", e.RubyPath, execPath)

		return nil, nil
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
		logrus.WithFields(logrus.Fields{
			"dryrun": e.DryRun,
			"error":  err.Error(),
		}).Errorf("Unable to decode %v", stackPath)

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

	logrus.WithFields(logrus.Fields{
		"dryrun": e.DryRun,
	}).Infof("Running stack: %v", name)

	for _, step := range stk.Steps {
		if strings.HasPrefix(step, "stacks/") {
			stackName := strings.Replace(step, "stacks/", "", -1)

			output, err := e.RunStack(stackName)
			if err != nil {
				logrus.WithFields(logrus.Fields{
					"dryrun": e.DryRun,
					"error":  err.Error(),
				}).Errorf("Unable to run stack: %v", stackName)

				return output, err
			}
		}

		if strings.HasPrefix(step, "logic/") {
			logicName := strings.Replace(step, "logic/", "", -1)

			output, err := e.RunLogic(logicName)
			if err != nil {
				logrus.WithFields(logrus.Fields{
					"dryrun": e.DryRun,
					"error":  err.Error(),
				}).Errorf("Unable to run logic: %v", logicName)

				return output, err
			}
		}
	}

	return make([]byte, 0), nil
}
