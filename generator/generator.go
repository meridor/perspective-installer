package generator

import (
	"fmt"
	. "github.com/meridor/perspective-installer/data"
	"encoding/xml"
	"os"
	"strings"
)

var (
	generators = make(map[string]Generator)
)

type Generator interface {
	Name() string
	Config(config ClusterConfig)
}

type BaseGenerator struct {
	Dir string
	DryRun bool
}

func InitGenerators(dir string, dryRun bool) {
	dirAware := BaseGenerator{dir, dryRun}
	addGenerator(DockerComposeGenerator{BaseGenerator: dirAware})
}

func addGenerator(generator Generator) {
	generators[generator.Name()] = generator
}

func marshalCloudsXml(cloudsXml CloudsXml) []byte {
	bytes, err := xml.MarshalIndent(cloudsXml, "", "    ")
	if (err != nil) {
		fmt.Printf("Failed to output clouds.xml: %v\n", err)
		os.Exit(1)
	}
	return bytes
}

func exitIfFailed(path string, err error) {
	if (err != nil) {
		fmt.Printf("Failed to save [%s]: %v. Exiting.", path, err)
		os.Exit(1)
	}
}

func prepareCloudType(cloudType CloudType) string {
	return strings.Replace(strings.ToLower(cloudType.String()), "_", "-", -1)
}

func GetNames() []string {
	generatorNames := make([]string, 0, len(generators))
	for k := range generators {
		generatorNames = append(generatorNames, k)
	}
	return generatorNames
}

func RunGenerators(config ClusterConfig, generatorNames []string) {
	for _, generatorName := range generatorNames {
		if gen, ok := generators[generatorName]; ok {
			gen.Config(config)
		} else {
			fmt.Printf("Skipping unsupported generator: %s\n", generatorName)
		}
	}
}
