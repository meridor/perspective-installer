package main

import (
	"fmt"
	"strings"
)

type Generator interface {
	Config(answers Answers) string
	Command(filename string) string
}

func Act(generatorNames []string) {
	answers := RunWizard()
	for _, generatorName := range generatorNames {
		if generator, ok := generators[generatorName]; ok {
			config := generator.Config(answers)
			if (file == "") {
				fmt.Println(generatorName)
				fmt.Println(strings.Repeat("-", len(generatorName)))
			}
			fmt.Print(config)
			fmt.Println()
			fmt.Printf("Use the following command to start cluster: %s\n", generator.Command(file))
		} else {
			fmt.Printf("Skipping unsupported generator: %s\n", generatorName)
		}
	}
}