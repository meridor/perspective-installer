package generator

import (
	"fmt"
	. "github.com/meridor/perspective-installer/data"
	"github.com/meridor/perspective-installer/wizard"
	"os"
	"path"
	"io/ioutil"
	"gopkg.in/yaml.v2"
)

const (
	fileMode = 0644
) 

type DockerComposeYml struct {
	Version  string                   `yaml:"version,omitempty"`
	Services map[string]ServiceConfig `yaml:"services,omitempty"`
	Networks map[string]interface{}   `yaml:"networks,omitempty"`
}

type ServiceConfig struct {
	ContainerName string                 `yaml:"container_name,omitempty"`
	DependsOn     []string               `yaml:"depends_on,omitempty"`
	Environment   map[string]interface{} `yaml:"environment,omitempty"`
	Expose        []string               `yaml:"expose,omitempty"`
	Image         string                 `yaml:"image,omitempty"`
	Labels        []string               `yaml:"labels,omitempty"`
	Links         []string               `yaml:"links,omitempty"`
	Ports         []string               `yaml:"ports,omitempty"`
	Privileged    bool                   `yaml:"privileged,omitempty"`
	Volumes       []string               `yaml:"volumes,omitempty"`
}

type DockerComposeGenerator struct {
	BaseGenerator
}

func (g DockerComposeGenerator) Name() string {
	return "docker-compose"
}

func (g DockerComposeGenerator) Config(config ClusterConfig) {
	fmt.Println("To use Docker Compose answer a bit more questions:")
	configDir := wizard.FreeInputQuestion("Specify configuration directory on Docker host machine:", "/etc/perspective")
	logsDir := wizard.FreeInputQuestion("Specify logs directory on Docker host machine:", "/var/log/perspective")
	if (wizard.YesNoQuestion("Will now create directories and configuration files. Proceed?", true)){
		g.createDirectory(logsDir)
		for cloudType, cloud := range config.Clouds {
			cloudsXmlPath := g.getCloudsXmlPath(configDir, cloudType)
			g.saveCloudsXml(cloudsXmlPath, cloud.XmlConfig)
		}
		dockerComposeYml := createDockerCompose(config)
		g.saveDockerCompose(dockerComposeYml)
	}
	composeYmlPath := createComposeYmlPath(g.Dir)
	fmt.Printf(
		"Use the following command to start cluster: docker-compose -f %s up\n",
		composeYmlPath,
	)
	fmt.Printf(
		"To completely remove cluster type: docker-compose -f %s down && rm -Rf %s && rm -Rf %s\n",
		composeYmlPath,
		configDir,
		logsDir,
	)
}

func (g DockerComposeGenerator) createDirectory(path string) {
	fmt.Printf("Creating directory [%s]...\n", path)
	if (!g.DryRun) {
		err := os.MkdirAll(path, fileMode)
		if (err != nil) {
			fmt.Printf("Failed to create directory [%s]: %v. Exiting.", path, err)
			os.Exit(1)
		}
	}
}

func (g DockerComposeGenerator) getCloudsXmlPath(configDir string, cloudType CloudType) string {
	dirName := prepareCloudType(cloudType)
	cloudDir := path.Join(configDir, dirName)
	g.createDirectory(cloudDir)
	return path.Join(cloudDir, "clouds.xml")
}

func (g DockerComposeGenerator) saveCloudsXml(path string, cloudsXml CloudsXml) {
	cloudsXmlContents := marshalCloudsXml(cloudsXml)
	fmt.Printf("Saving [%s]...\n", path)
	if (g.DryRun) {
		fmt.Println(string(cloudsXmlContents))
	} else {
		err := ioutil.WriteFile(path, cloudsXmlContents, fileMode)
		exitIfFailed(path, err)
	}
}

func createDockerCompose(config ClusterConfig) DockerComposeYml {
	dockerComposeYml := DockerComposeYml{
		Version: "2.1",
		Services: make(map[string] ServiceConfig),
	}
	dockerComposeYml.Services["storage"] = createStorageService(config)
	dockerComposeYml.Services["rest"] = createRestService(config)
	for cloudType := range config.Clouds {
		workerServiceName := prepareCloudType(cloudType)
		dockerComposeYml.Services[workerServiceName] = createWorkerService(cloudType, config.Version)
	} 
	return dockerComposeYml
}

func createEnvironment() map[string] interface{} {
	env := make(map[string] interface{})
	env["MISC_PROPERTIES"] = "-Dperspective.storage.hosts=storage:5801"
	return env
}

func createStorageService(config ClusterConfig) ServiceConfig {
	return ServiceConfig{
		ContainerName: "perspective-storage",
		Image: fmt.Sprintf("meridor/perspective-storage:%s", config.Version),
	}
}

func createRestService(config ClusterConfig) ServiceConfig {
	return ServiceConfig{
		ContainerName: "perspective-rest",
		Image: fmt.Sprintf("meridor/perspective-rest:%s", config.Version),
		Ports: []string{fmt.Sprintf("8080:%d", config.ApiPort)},
		Environment: createEnvironment(),
		Links: []string{"storage"},
		DependsOn: []string{"storage"},
	}
}

func createWorkerService(cloudType CloudType, version string) ServiceConfig {
	suffix := prepareCloudType(cloudType)
	return ServiceConfig{
		ContainerName: fmt.Sprintf("perspective-%s", suffix),
		Image: fmt.Sprintf("meridor/%s:%s", suffix, version),
		Environment: createEnvironment(),
		Links: []string{"storage"},
		DependsOn: []string{"rest"},
	}
}

func (g DockerComposeGenerator) saveDockerCompose(composeYml DockerComposeYml) {
	bytes, err := yaml.Marshal(&composeYml)
	if (err != nil) {
		fmt.Printf("Failed to generate docker-compose.yml contents: %v\n", err)
	}
	ymlString := string(bytes)
	ymlPath := createComposeYmlPath(g.Dir)
	fmt.Printf("Saving [%s]...\n", ymlPath)
	if (g.DryRun) {
		fmt.Println(ymlString)
	} else {
		g.createDirectory(g.Dir)
		err := ioutil.WriteFile(ymlPath, bytes, fileMode)
		exitIfFailed(ymlPath, err)
	}
}

func createComposeYmlPath(dir string) string {
	return path.Join(dir, "docker-compose.yml")
}