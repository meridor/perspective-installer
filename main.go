package main

import (
	"flag"
	"fmt"
	"os"
)

var (
	file string
	generators = make(map[string] Generator)
)

func init() {
	
	generators["docker-compose"] = DockerComposeGenerator{}
	
	flag.StringVar(&file, "-o", "", "File to output generated config")
	flag.Parse()
}

func main() {
	if (len(flag.Args()) == 0) {
		fmt.Println("Usage: perspective-installer [-o filename] generator1 generator2 ...")
		generatorNames := make([]string, 0, len(generators))
		for k := range generators {
			generatorNames = append(generatorNames, k)
		}
		fmt.Println("The following generators are supported:")
		for _, gn := range generatorNames {
			fmt.Printf("* %s\n", gn)
		}
		os.Exit(1)
	}
	Act(flag.Args())
	os.Exit(0)
}