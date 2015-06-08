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

func TestRunLogic(t *testing.T) {
	engine, err := New(os.ExpandEnv("$GOPATH/src/github.com/resourced/configurator/blank"))
	if err != nil {
		t.Fatalf("Creating new engine should not fail. Error: %v", err)
	}

	output, err := engine.RunLogic("ignore-me")
	if err != nil {
		t.Fatalf("RunLogic should not fail in dry run mode. Error: %v", err)
	}
	if !strings.Contains(string(output), "/python") {
		t.Fatalf("RunLogic should contain python command in dry run mode. Output: %v", string(output))
	}
	if !strings.HasSuffix(string(output), "__init__.py") {
		t.Fatalf("RunLogic should contain __init__.py in dry run mode. Output: %v", string(output))
	}
}
