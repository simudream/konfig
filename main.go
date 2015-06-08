package main

import (
	"flag"
	"log"

	"github.com/kardianos/osext"
	"github.com/resourced/configurator/engine"
)

func main() {
	var root string
	var err error

	rootInput := flag.String("root", "", "Project root directory")
	pythonInput := flag.String("python", "", "Path to Python executable")
	rubyInput := flag.String("ruby", "", "Path to Ruby executable")

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

	if *pythonInput != "" {
		engine.PythonPath = *pythonInput
	}
	if *rubyInput != "" {
		engine.RubyPath = *rubyInput
	}

	for _, f := range engine.Logic {
		println(f.Name())
	}
}
