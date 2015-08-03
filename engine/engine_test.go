package engine

import (
	"os"
	"strings"
	"testing"
)

func TestConstructor(t *testing.T) {
	engine, err := New(os.ExpandEnv("$GOPATH/src/github.com/resourced/resourced-stacks/tests/project"), "")
	if err != nil {
		t.Fatalf("Creating new engine should not fail. Error: %v", err)
	}
	if engine.PythonPath == "" {
		t.Fatalf("Engine Python path should not be empty.")
	}
	if engine.DryRun != true {
		t.Fatalf("DryRun should be true bydefault.")
	}
}

func TestIsGitRepo(t *testing.T) {
	engine := &Engine{Root: os.ExpandEnv("$GOPATH/src/github.com/resourced/resourced-stacks")}
	if !engine.IsGitRepo() {
		t.Fatalf(os.ExpandEnv("$GOPATH/src/github.com/resourced/resourced-stacks should be a git repo."))
	}
}

func TestInstallPythonLogicDependencies(t *testing.T) {
	engine, err := New(os.ExpandEnv("$GOPATH/src/github.com/resourced/resourced-stacks/tests/project"), "")
	if err != nil {
		t.Fatalf("Creating new engine should not fail. Error: %v", err)
	}

	_, err = engine.InstallPythonLogicDependencies("helloworld")
	if err != nil {
		t.Fatalf("InstallPythonLogicDependencies should not fail in dry run mode. Error: %v", err)
	}

	engine.DryRun = false
	output, _ := engine.InstallPythonLogicDependencies("helloworld")
	if !strings.Contains(string(output), "Installing collected packages: requests") && !strings.Contains(string(output), "Requirement already satisfied (use --upgrade to upgrade): requests") {
		t.Fatalf("InstallPythonLogicDependencies should contain 'Installing collected packages: requests' in live mode. Output: %v", string(output))
	}
}

func TestRunPythonLogic(t *testing.T) {
	engine, err := New(os.ExpandEnv("$GOPATH/src/github.com/resourced/resourced-stacks/tests/project"), "")
	if err != nil {
		t.Fatalf("Creating new engine should not fail. Error: %v", err)
	}

	output, err := engine.RunPythonLogic("helloworld", nil)
	if err != nil {
		t.Fatalf("RunPythonLogic should not fail in dry run mode. Error: %v, Output: %v", err, string(output))
	}

	engine.DryRun = false
	output, err = engine.RunPythonLogic("helloworld", nil)
	if err != nil {
		t.Fatalf("RunPythonLogic should not fail in live mode. Error: %v, Output: %v", err, output)
	}
	if !strings.Contains(strings.ToLower(string(output)), "hello world") {
		t.Fatalf("RunPythonLogic should contain hello world in live mode. Output: %v", string(output))
	}
}

func TestReadStack(t *testing.T) {
	engine, err := New(os.ExpandEnv("$GOPATH/src/github.com/resourced/resourced-stacks/tests/project"), "")
	if err != nil {
		t.Fatalf("Creating new engine should not fail. Error: %v", err)
	}

	stk, err := engine.ReadStack("helloworld")
	if err != nil {
		t.Fatalf("ReadStack should not fail. Error: %v", err)
	}
	if len(stk.Steps) != 1 {
		t.Fatalf("stack steps should == 1. Length: %v", len(stk.Steps))
	}
}

func TestRunStack(t *testing.T) {
	engine, err := New(os.ExpandEnv("$GOPATH/src/github.com/resourced/resourced-stacks/tests/project"), "")
	if err != nil {
		t.Fatalf("Creating new engine should not fail. Error: %v", err)
	}

	_, err = engine.RunStack("helloworld", nil)
	if err != nil {
		t.Fatalf("RunStack should not fail. Error: %v", err)
	}
}
