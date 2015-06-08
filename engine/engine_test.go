package engine

import (
	"os"
	"strings"
	"testing"
)

func TestConstructor(t *testing.T) {
	engine, err := New(os.ExpandEnv("$GOPATH/src/github.com/resourced/configurator/blank"))
	if err != nil {
		t.Fatalf("Creating new engine should not fail. Error: %v", err)
	}
	if engine.PythonPath == "" {
		t.Fatalf("Engine Python path should not be empty.")
	}
	if engine.RubyPath == "" {
		t.Fatalf("Engine Ruby path should not be empty.")
	}
	if engine.DryRun != true {
		t.Fatalf("DryRun should be true bydefault.")
	}
}

func TestInstallPythonLogicDependencies(t *testing.T) {
	engine, err := New(os.ExpandEnv("$GOPATH/src/github.com/resourced/configurator/blank"))
	if err != nil {
		t.Fatalf("Creating new engine should not fail. Error: %v", err)
	}

	output, err := engine.InstallPythonLogicDependencies("helloworld-py")
	if err != nil {
		t.Fatalf("InstallPythonLogicDependencies should not fail in dry run mode. Error: %v", err)
	}
	if !strings.Contains(string(output), "/pip") {
		t.Fatalf("InstallPythonLogicDependencies should contain python command in dry run mode. Output: %v", string(output))
	}
	if !strings.HasSuffix(string(output), "requirements.txt") {
		t.Fatalf("InstallPythonLogicDependencies should contain requirements.txt in dry run mode. Output: %v", string(output))
	}

	engine.DryRun = false
	output, _ = engine.InstallPythonLogicDependencies("helloworld-py")
	if !strings.Contains(string(output), "Installing collected packages: requests") {
		t.Fatalf("InstallPythonLogicDependencies should contain 'Installing collected packages: requests' in live mode. Output: %v", string(output))
	}
}

func TestRunPythonLogic(t *testing.T) {
	engine, err := New(os.ExpandEnv("$GOPATH/src/github.com/resourced/configurator/blank"))
	if err != nil {
		t.Fatalf("Creating new engine should not fail. Error: %v", err)
	}

	output, err := engine.RunPythonLogic("helloworld-py")
	if err != nil {
		t.Fatalf("RunPythonLogic should not fail in dry run mode. Error: %v", err)
	}
	if !strings.Contains(string(output), "python") {
		t.Fatalf("RunPythonLogic should contain python command in dry run mode. Output: %v", string(output))
	}
	if !strings.HasSuffix(string(output), "__init__.py") {
		t.Fatalf("RunPythonLogic should contain __init__.py in dry run mode. Output: %v", string(output))
	}

	engine.DryRun = false
	output, err = engine.RunPythonLogic("helloworld-py")
	if err != nil {
		t.Fatalf("RunPythonLogic should not fail in live mode. Error: %v", err)
	}
	if !strings.Contains(string(output), "Hello World") {
		t.Fatalf("RunPythonLogic should contain Hello World in live mode. Output: %v", string(output))
	}
}

func TestInstallRubyLogicDependencies(t *testing.T) {
	engine, err := New(os.ExpandEnv("$GOPATH/src/github.com/resourced/configurator/blank"))
	if err != nil {
		t.Fatalf("Creating new engine should not fail. Error: %v", err)
	}

	output, err := engine.InstallRubyLogicDependencies("helloworld-rb")
	if err != nil {
		t.Fatalf("InstallRubyLogicDependencies should not fail in dry run mode. Error: %v", err)
	}
	if !strings.Contains(string(output), "bundle") {
		t.Fatalf("InstallRubyLogicDependencies should contain bundle command in dry run mode. Output: %v", string(output))
	}
}

func TestRunRubyLogic(t *testing.T) {
	engine, err := New(os.ExpandEnv("$GOPATH/src/github.com/resourced/configurator/blank"))
	if err != nil {
		t.Fatalf("Creating new engine should not fail. Error: %v", err)
	}

	output, err := engine.RunRubyLogic("helloworld-rb")
	if err != nil {
		t.Fatalf("RunRubyLogic should not fail in dry run mode. Error: %v", err)
	}
	if !strings.Contains(string(output), "ruby") {
		t.Fatalf("RunRubyLogic should contain ruby command in dry run mode. Output: %v", string(output))
	}
	if !strings.HasSuffix(string(output), "helloworld-rb.rb") {
		t.Fatalf("RunRubyLogic should contain helloworld-rb.rb in dry run mode. Output: %v", string(output))
	}
}
