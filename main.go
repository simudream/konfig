package main

import (
	"bufio"
	"bytes"
	"errors"
	"flag"
	"os/exec"

	"github.com/Sirupsen/logrus"
	"github.com/resourced/configurator/engine"
)

func main() {
	cmdInput := flag.String("cmd", "run", "Command")
	rootInput := flag.String("root", "", "Project root directory")
	pythonInput := flag.String("python", "", "Path to python executable")
	pipInput := flag.String("pip", "", "Path to pip executable")
	rubyInput := flag.String("ruby", "", "Path to ruby executable")
	bundleInput := flag.String("bundle", "", "Path to bundle executable")
	dryRunInput := flag.Bool("dryrun", true, "Dry run mode, default to true")

	flag.Parse()

	if *rootInput == "" {
		logrus.Fatal(errors.New("root directory must be specified."))
	}

	if *cmdInput == "run" {
		engine, err := engine.New(*rootInput)
		if err != nil {
			logrus.Fatal(err)
		}

		engine.DryRun = *dryRunInput

		if *pythonInput != "" {
			engine.PythonPath = *pythonInput
		}
		if *pipInput != "" {
			engine.PipPath = *pipInput
		}
		if *rubyInput != "" {
			engine.RubyPath = *rubyInput
		}
		if *bundleInput != "" {
			engine.BundlePath = *bundleInput
		}

		output, err := engine.RunRoles()
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

	if *cmdInput == "new" {
		// Download blank template folder to *rootInput.
		// svn export http://github.com/resourced/configurator/trunk/blank target/project
		output, err := exec.Command("svn", "export", "https://github.com/resourced/configurator/trunk/blank", *rootInput).CombinedOutput()
		if err != nil {
			logrus.Fatal(err)
		}

		logrus.Info(string(output))
	}
}
