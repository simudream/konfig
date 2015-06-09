package main

import (
	"bufio"
	"bytes"
	"flag"
	"log"

	"github.com/kardianos/osext"
	"github.com/resourced/configurator/engine"
)

func main() {
	var root string
	var err error

	rootInput := flag.String("root", "", "Project root directory")
	pythonInput := flag.String("python", "", "Path to python executable")
	pipInput := flag.String("pip", "", "Path to pip executable")
	rubyInput := flag.String("ruby", "", "Path to ruby executable")
	bundleInput := flag.String("bundle", "", "Path to bundle executable")
	dryRunInput := flag.Bool("dryrun", true, "Dry run mode, default to true")

	flag.Parse()

	if *rootInput == "" {
		root, err = osext.ExecutableFolder()
		if err != nil {
			log.Fatal(err)
		}

	} else {
		root = *rootInput
	}

	engine, err := engine.New(root)
	if err != nil {
		log.Fatal(err)
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
		log.Fatal(err)
	}

	scanner := bufio.NewScanner(bytes.NewReader(output))
	for scanner.Scan() {
		if scanner.Text() != "" {
			log.Print(scanner.Text())
		}
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}
