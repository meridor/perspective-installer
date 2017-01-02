package generator

import (
	. "github.com/meridor/perspective-installer/data"
	"fmt"
	"strings"
)

var (
	generators = make(map[string] Generator)
)

type Generator interface {
	Name() string
	Config(clouds map[CloudType] Cloud) string
	Command(dir string) string
}

func init() {
	addGenerator(DockerComposeGenerator{})
}

func addGenerator(generator Generator) {
	generators[generator.Name()] = generator
}

func GetNames() []string {
	generatorNames := make([]string, 0, len(generators))
	for k := range generators {
		generatorNames = append(generatorNames, k)
	}
	return generatorNames
}

func RunGenerators(dir string, clouds map[CloudType] Cloud, generatorNames []string) {
	for _, generatorName := range generatorNames {
		if gen, ok := generators[generatorName]; ok {
			config := gen.Config(clouds)
			if (dir == "") {
				fmt.Println(generatorName)
				fmt.Println(strings.Repeat("-", len(generatorName)))
			}
			fmt.Print(config)
			fmt.Println()
			fmt.Printf("Use the following command to start cluster: %s\n", gen.Command(dir))
		} else {
			fmt.Printf("Skipping unsupported generator: %s\n", generatorName)
		}
	}
}