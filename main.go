package main

import (
	"bufio"
	"bytes"
	"errors"
	"flag"
	"os"

	"github.com/Sirupsen/logrus"
	"github.com/resourced/resourced-stacks/engine"
)

func main() {
	cmdInput := flag.String("cmd", "run", "Command")
	rootInput := flag.String("root", "", "Project root directory")
	stackInput := flag.String("stack", "", "Stack to run")
	conditionsInput := flag.String("conditions", "true", "Conditions to match before running the command")
	pythonInput := flag.String("python", "python", "Path to python executable")
	pipInput := flag.String("pip", "pip", "Path to pip executable")
	dryRunInput := flag.Bool("dryrun", true, "Dry run mode")

	// git related options
	gitInput := flag.String("git", "", "HTTPS URL to git repo")
	gitBranchInput := flag.String("git-branch", "master", "Checkout a specific branch")

	flag.Parse()

	if *rootInput == "" {
		logrus.Fatal(errors.New("root directory must be specified."))
	}
	err := os.MkdirAll(*rootInput, 0755)
	if err != nil {
		logrus.Fatal(err)
	}

	engine, err := engine.New(*rootInput, *conditionsInput)
	if err != nil {
		logrus.Fatal(err)
	}

	engine.DryRun = *dryRunInput
	engine.Git.HTTPS = *gitInput
	engine.Git.Branch = *gitBranchInput

	if *cmdInput == "run" {
		if *stackInput == "" {
			logrus.Fatal(errors.New("stack name must be specified."))
		}

		if *pythonInput != "" {
			engine.PythonPath = *pythonInput
		}
		if *pipInput != "" {
			engine.PipPath = *pipInput
		}

		conditionOutput, err := engine.EvalConditions()
		if err != nil {
			logrus.Fatal(err)
		}
		if !conditionOutput {
			logrus.Info("Conditions are not met")
			os.Exit(0)
		}

		output, err := engine.RunStack(*stackInput, nil)
		if err != nil {
			scanner := bufio.NewScanner(bytes.NewReader(output))
			for scanner.Scan() {
				if scanner.Text() != "" {
					logrus.Error(scanner.Text())
				}
			}
			if err := scanner.Err(); err != nil {
				logrus.Fatal(err)
			}

			logrus.Fatal(err)
		}
	}

	if *cmdInput == "pull" {
		err := engine.GitPull()
		if err != nil {
			logrus.Fatal(err)
		}
	}

	if *cmdInput == "new" {
		err := engine.NewProject()
		if err != nil {
			logrus.Fatal(err)
		}
	}
}
