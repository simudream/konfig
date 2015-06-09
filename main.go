package main

import (
	"bufio"
	"bytes"
	"errors"
	"flag"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/resourced/configurator/engine"
)

func main() {
	cmdInput := flag.String("cmd", "run", "Configurator command")
	rootInput := flag.String("root", "", "Project root directory")
	pythonInput := flag.String("python", "", "Path to python executable")
	pipInput := flag.String("pip", "", "Path to pip executable")
	rubyInput := flag.String("ruby", "", "Path to ruby executable")
	bundleInput := flag.String("bundle", "", "Path to bundle executable")
	dryRunInput := flag.Bool("dryrun", true, "Dry run mode")

	flag.Parse()

	if *rootInput == "" {
		logrus.Fatal(errors.New("root directory must be specified."))
	}
	err := os.MkdirAll(*rootInput, 0755)
	if err != nil {
		logrus.Fatal(err)
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
		// 1. Create tmp directory.
		dir, err := ioutil.TempDir(os.TempDir(), "configurator")
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"error": err.Error(),
			}).Fatal(err)
		}
		defer os.RemoveAll(dir)

		// 2. git clone to /tmp directory.
		output, err := exec.Command("git", "clone", "https://github.com/resourced/configurator.git", dir).CombinedOutput()
		if err != nil {
			os.RemoveAll(dir)

			logrus.WithFields(logrus.Fields{
				"error": err.Error(),
			}).Fatal(string(output))
		}

		logrus.Info(string(output))

		// 3. mv blank template folder to *rootInput
		logrus.Infof("Moving %v to %v...", path.Join(dir, "blank"), *rootInput)
		err = os.Rename(path.Join(dir, "blank"), *rootInput)
		if err != nil {
			os.RemoveAll(dir)

			if !strings.Contains(err.Error(), "directory not empty") {
				logrus.WithFields(logrus.Fields{
					"error": err.Error(),
				}).Fatal(string(output))
			}
		}
	}
}
