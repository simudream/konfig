package engine

import (
	"os"
	"strings"
	"testing"
)

func TestConstructor(t *testing.T) {
	engine, err := New(os.ExpandEnv("$GOPATH/src/github.com/resourced/configurator/tests/project"))
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

func TestIsGitRepo(t *testing.T) {
	engine := &Engine{Root: os.ExpandEnv("$GOPATH/src/github.com/resourced/configurator")}
	if !engine.IsGitRepo() {
		t.Fatalf(os.ExpandEnv("$GOPATH/src/github.com/resourced/configurator should be a git repo."))
	}
}

func TestInstallPythonLogicDependencies(t *testing.T) {
	engine, err := New(os.ExpandEnv("$GOPATH/src/github.com/resourced/configurator/tests/project"))
	if err != nil {
		t.Fatalf("Creating new engine should not fail. Error: %v", err)
	}

	_, err = engine.InstallPythonLogicDependencies("helloworld-py")
	if err != nil {
		t.Fatalf("InstallPythonLogicDependencies should not fail in dry run mode. Error: %v", err)
	}

	engine.DryRun = false
	output, _ := engine.InstallPythonLogicDependencies("helloworld-py")
	if !strings.Contains(string(output), "Installing collected packages: requests") {
		t.Fatalf("InstallPythonLogicDependencies should contain 'Installing collected packages: requests' in live mode. Output: %v", string(output))
	}
}

func TestRunPythonLogic(t *testing.T) {
	engine, err := New(os.ExpandEnv("$GOPATH/src/github.com/resourced/configurator/tests/project"))
	if err != nil {
		t.Fatalf("Creating new engine should not fail. Error: %v", err)
	}

	_, err = engine.RunPythonLogic("helloworld-py")
	if err != nil {
		t.Fatalf("RunPythonLogic should not fail in dry run mode. Error: %v", err)
	}

	engine.DryRun = false
	output, err := engine.RunPythonLogic("helloworld-py")
	if err != nil {
		t.Fatalf("RunPythonLogic should not fail in live mode. Error: %v, Output: %v", err, output)
	}
	if !strings.Contains(string(output), "Hello World") {
		t.Fatalf("RunPythonLogic should contain Hello World in live mode. Output: %v", string(output))
	}
}

func TestInstallRubyLogicDependencies(t *testing.T) {
	engine, err := New(os.ExpandEnv("$GOPATH/src/github.com/resourced/configurator/tests/project"))
	if err != nil {
		t.Fatalf("Creating new engine should not fail. Error: %v", err)
	}

	_, err = engine.InstallRubyLogicDependencies("helloworld-rb")
	if err != nil {
		t.Fatalf("InstallRubyLogicDependencies should not fail in dry run mode. Error: %v", err)
	}
}

func TestRunRubyLogic(t *testing.T) {
	engine, err := New(os.ExpandEnv("$GOPATH/src/github.com/resourced/configurator/tests/project"))
	if err != nil {
		t.Fatalf("Creating new engine should not fail. Error: %v", err)
	}

	_, err = engine.RunRubyLogic("helloworld-rb")
	if err != nil {
		t.Fatalf("RunRubyLogic should not fail in dry run mode. Error: %v", err)
	}
}

func TestReadStack(t *testing.T) {
	engine, err := New(os.ExpandEnv("$GOPATH/src/github.com/resourced/configurator/tests/project"))
	if err != nil {
		t.Fatalf("Creating new engine should not fail. Error: %v", err)
	}

	stk, err := engine.ReadStack("helloworld")
	if err != nil {
		t.Fatalf("ReadStack should not fail. Error: %v", err)
	}
	if len(stk.Steps) != 2 {
		t.Fatalf("stack steps should == 2. Length: %v", len(stk.Steps))
	}
}

func TestRunStack(t *testing.T) {
	engine, err := New(os.ExpandEnv("$GOPATH/src/github.com/resourced/configurator/tests/project"))
	if err != nil {
		t.Fatalf("Creating new engine should not fail. Error: %v", err)
	}

	_, err = engine.RunStack("helloworld")
	if err != nil {
		t.Fatalf("RunStack should not fail. Error: %v", err)
	}
}

func TestReadRole(t *testing.T) {
	engine, err := New(os.ExpandEnv("$GOPATH/src/github.com/resourced/configurator/tests/project"))
	if err != nil {
		t.Fatalf("Creating new engine should not fail. Error: %v", err)
	}

	rl, err := engine.ReadRole("helloworld-staging.toml")
	if err != nil {
		t.Fatalf("ReadRole should not fail. Error: %v", err)
	}
	if len(rl.Steps) != 2 {
		t.Fatalf("role steps should == 2. Length: %v", len(rl.Steps))
	}
	if rl.Steps[0] != "stacks/helloworld.toml" {
		t.Fatalf("role steps[0] should == stacks/helloworld.toml. Step: %v", rl.Steps[0])
	}

	if len(rl.Matchers.Hostname) == 0 {
		t.Fatalf("role hostname matcher should not be empty. Length: %v", len(rl.Matchers.Hostname))
	}
	if rl.Matchers.Hostname[0] != "=" {
		t.Fatalf("role hostname matcher operator should be =. Operator: %v", rl.Matchers.Hostname[0])
	}
	if rl.Matchers.Hostname[1] != "$HOSTNAME" {
		t.Fatalf("role hostname matcher value should be $HOSTNAME. Value: %v", rl.Matchers.Hostname[1])
	}
}

func TestRunRole(t *testing.T) {
	engine, err := New(os.ExpandEnv("$GOPATH/src/github.com/resourced/configurator/tests/project"))
	if err != nil {
		t.Fatalf("Creating new engine should not fail. Error: %v", err)
	}

	_, err = engine.RunRole("helloworld-staging.toml")
	if err != nil {
		t.Fatalf("RunRole should not fail. Error: %v", err)
	}
}

func TestCheckIfRoleMatchedByHostname(t *testing.T) {
	engine, err := New(os.ExpandEnv("$GOPATH/src/github.com/resourced/configurator/tests/project"))
	if err != nil {
		t.Fatalf("Creating new engine should not fail. Error: %v", err)
	}

	rl, err := engine.ReadRole("helloworld-staging.toml")
	if err != nil {
		t.Fatalf("ReadRole should not fail. Error: %v", err)
	}

	isMatched := engine.checkIfRoleMatchedByHostname(rl)
	if isMatched != true {
		t.Fatalf(`["=", "$HOSTNAME"] should always match any hostname. IsMatched: %v`, isMatched)
	}
}
