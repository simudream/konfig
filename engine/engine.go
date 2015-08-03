// Package engine provides engine struct.
package engine

import (
	"bytes"
	"encoding/json"
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

	logicDirs, err := ioutil.ReadDir(path.Join(engine.Root, "logic"))
	if err != nil {
		return nil, err
	}

	stackFiles, err := ioutil.ReadDir(path.Join(engine.Root, "stacks"))
	if err != nil {
		return nil, err
	}

	engine.Logic = logicDirs
	engine.Stacks = stackFiles

	engine.DryRun = true
	engine.PythonPath = "python"
	engine.PipPath = "pip"
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

// RunLogic allows engine to execute one logic layer.
func (e *Engine) RunLogic(name string, data map[string]interface{}) ([]byte, error) {
	logrus.WithFields(logrus.Fields{
		"dryrun": e.DryRun,
	}).Infof("Starting logic: %v", name)

	pythonExecPath := path.Join(e.Root, "logic", name, "__init__.py")
	_, pyErr := os.Stat(pythonExecPath)
	if pyErr == nil {
		return e.RunPythonLogic(name, data)
	}

	if os.IsNotExist(pyErr) {
		err := errors.New(fmt.Sprintf("Logic must be implemented in Python(%v/__init__.py)", name))

		logrus.WithFields(logrus.Fields{
			"dryrun": e.DryRun,
			"error":  err.Error(),
			"path":   pythonExecPath,
		}).Errorf("Unable to run logic: %v", name)

		return nil, err
	}

	return nil, nil
}

// InstallPythonLogicDependencies allows engine to installs dependencies for a logic written in python.
func (e *Engine) InstallPythonLogicDependencies(name string) ([]byte, error) {
	reqPath := path.Join(e.Root, "logic", name, "requirements.txt")

	commandChunks := []string{e.PipPath, "install", "-r", reqPath}

	_, err := os.Stat(reqPath)
	if err != nil {
		return nil, err
	}

	if e.DryRun {
		logrus.WithFields(logrus.Fields{
			"dryrun": e.DryRun,
		}).Infof("Executing: " + strings.Join(commandChunks, " "))

		return nil, nil
	}

	output, err := exec.Command(commandChunks[0], commandChunks[1:]...).CombinedOutput()
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"dryrun": e.DryRun,
			"error":  err.Error(),
		}).Infof("Failed executing: " + strings.Join(commandChunks, " "))
	} else {
		logrus.WithFields(logrus.Fields{
			"dryrun": e.DryRun,
		}).Infof("Executed: " + strings.Join(commandChunks, " "))
	}

	return output, err
}

// RunPythonLogic allows engine to run a logic written in python.
func (e *Engine) RunPythonLogic(name string, data map[string]interface{}) ([]byte, error) {
	_, err := e.InstallPythonLogicDependencies(name)
	if err != nil {
		return nil, err
	}

	inJson, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	execPath := path.Join(e.Root, "logic", name, "__init__.py")
	commandChunks := []string{e.PythonPath, execPath}

	if e.DryRun {
		commandChunks = append(commandChunks, "--dryrun")
	} else {
		commandChunks = append(commandChunks, "--no-dryrun")
	}

	cmd := exec.Command(commandChunks[0], commandChunks[1:]...)
	cmd.Stdin = bytes.NewReader(inJson)

	output, err := cmd.CombinedOutput()
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"dryrun": e.DryRun,
			"error":  err.Error(),
		}).Infof("Failed executing: " + strings.Join(commandChunks, " "))
	} else {
		logrus.WithFields(logrus.Fields{
			"dryrun": e.DryRun,
		}).Infof("Executed: " + strings.Join(commandChunks, " "))
	}

	return output, err
}

// ReadStack allows engine to read a particular stack defined in TOML file.
func (e *Engine) ReadStack(name string) (stack.Stack, error) {
	var stk stack.Stack

	stackPath := path.Join(e.Root, "stacks", name, name+".toml")
	if _, err := toml.DecodeFile(stackPath, &stk); err != nil {
		logrus.WithFields(logrus.Fields{
			"dryrun": e.DryRun,
			"error":  err.Error(),
		}).Errorf("Unable to decode %v", stackPath)

		return stk, err
	}

	return stk, nil
}

// ReadStackData allows engine to read stack data defined in JSON.
// This data will be passed to logics via STDIN.
func (e *Engine) ReadStackData(name string) (map[string]interface{}, error) {
	data := make(map[string]interface{})
	dataPath := path.Join(e.Root, "stacks", name, "data")

	// Skip if data directory does not exist.
	if _, err := os.Stat(dataPath); err != nil {
		if os.IsNotExist(err) {
			return data, nil
		}
	}

	jsonFiles, err := ioutil.ReadDir(dataPath)
	if err != nil {
		return nil, err
	}

	for _, jsonFile := range jsonFiles {
		if strings.HasSuffix(jsonFile.Name(), ".json") {
			fileContent, err := ioutil.ReadFile(path.Join(dataPath, jsonFile.Name()))
			if err != nil {
				return nil, err
			}

			var jsonData interface{}
			err = json.Unmarshal(fileContent, &jsonData)
			if err != nil {
				return nil, err
			}

			data[strings.Replace(jsonFile.Name(), ".json", "", -1)] = jsonData
		}
	}

	return data, nil
}

// RunStack allows engine to run a particular stack.
func (e *Engine) RunStack(name string, data map[string]interface{}) ([]byte, error) {
	logrus.WithFields(logrus.Fields{
		"dryrun": e.DryRun,
	}).Infof("Starting stack: %v", name)

	stk, err := e.ReadStack(name)
	if err != nil {
		return nil, err
	}

	// Create data if nil
	if data == nil {
		data = make(map[string]interface{})
	}

	newData, err := e.ReadStackData(name)
	if err != nil {
		return nil, err
	}

	for key, value := range newData {
		if _, ok := data[key]; !ok {
			data[key] = value
		}
	}

	for _, step := range stk.Steps {
		var output []byte
		var outputInterface interface{}

		if strings.HasPrefix(step, "stacks/") {
			stackName := strings.Replace(step, "stacks/", "", -1)

			output, err = e.RunStack(stackName, data)
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

			output, err = e.RunLogic(logicName, data)
			if err != nil {
				logrus.WithFields(logrus.Fields{
					"dryrun": e.DryRun,
					"error":  err.Error(),
					"output": string(output),
				}).Errorf("Unable to run logic: %v", logicName)

				return output, err
			}
		}

		// Capture previous output and pass it as part of data
		if output != nil && len(output) > 0 {
			err = json.Unmarshal(output, &outputInterface)
			if err != nil {
				return output, err
			}
			data["previous_step"] = outputInterface
		}
	}

	return make([]byte, 0), nil
}
