package engine

import (
	"os"
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
}
