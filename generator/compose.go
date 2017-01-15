package generator

import (
	"fmt"
	. "github.com/meridor/perspective-installer/data"
	"github.com/meridor/perspective-installer/wizard"
	"os"
	"path"
	"io/ioutil"
	"gopkg.in/yaml.v2"
	"bytes"
)

const (
	fileMode = 0644
	rest = "rest"
	storage = "storage"
	log4jProperties = "log4j.properties"
	restProperties = "rest.properties"
	cloudsXml = "clouds.xml"
	workerProperties = "worker.properties"
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
	configDir string
	logsDir string
}

func (g DockerComposeGenerator) Name() string {
	return "docker-compose"
}

func (g DockerComposeGenerator) Config(config ClusterConfig) {
	fmt.Println("To use Docker Compose answer a bit more questions:")
	g.configDir = wizard.FreeInputQuestion("Specify configuration directory on Docker host machine:", "/etc/perspective")
	g.logsDir = wizard.FreeInputQuestion("Specify logs directory on Docker host machine:", "/var/log/perspective")
	if (wizard.YesNoQuestion("Will now create directories and configuration files. Proceed?", true)){
		g.createDirectory(g.logsDir)
		
		//Storage configs
		storageConfigDir := path.Join(g.configDir, storage)
		g.createDirectory(storageConfigDir)
		g.saveProperties(path.Join(storageConfigDir, log4jProperties), g.getLoggingProperties(storage))
		
		//Worker configs
		for cloudType, cloud := range config.Clouds {
			workerConfigDir := g.getWorkerConfigDir(cloudType)
			g.createDirectory(workerConfigDir)
			cloudsXmlPath := g.getCloudsXmlPath(cloudType)
			g.saveCloudsXml(cloudsXmlPath, cloud.XmlConfig)
			g.saveProperties(
				g.getWorkerPropertiesPath(cloudType),
				g.getStorageProperties(),
			)
			g.saveProperties(
				g.getWorkerLoggingPropertiesPath(cloudType),
				g.getLoggingProperties(prepareCloudType(cloudType)),
			)
		}
		
		//Rest configs
		restConfigDir := path.Join(g.configDir, rest)
		g.createDirectory(restConfigDir)
		g.saveProperties(path.Join(restConfigDir, restProperties), g.getStorageProperties())
		g.saveProperties(path.Join(restConfigDir, log4jProperties), g.getLoggingProperties(rest))
		
		
		//docker-compose.yml
		dockerComposeYml := g.createDockerCompose(config)
		g.createDirectory(g.Dir)
		g.saveDockerCompose(dockerComposeYml)
	}
	composeYmlPath := createComposeYmlPath(g.Dir)
	fmt.Printf(
		"Use the following command to start cluster: docker-compose -f ./%s pull && docker-compose -f %s up -d\n",
		composeYmlPath,
		composeYmlPath,
	)
	fmt.Printf(
		"To completely remove cluster type: docker-compose -f %s down && rm -Rf %s && rm -Rf %s\n",
		composeYmlPath,
		g.configDir,
		g.logsDir,
	)
}

func (g DockerComposeGenerator) createDirectory(path string) {
	fmt.Printf("Creating directory [%s]...\n", path)
	if (!g.DryRun) {
		err := os.MkdirAll(path, fileMode)
		if (err != nil) {
			fmt.Printf("Failed to create directory [%s]: %v. Exiting.\n", path, err)
			os.Exit(1)
		}
	}
}

func (g DockerComposeGenerator) getWorkerConfigDir(cloudType CloudType) string {
	return path.Join(g.configDir, prepareCloudType(cloudType))
}

func (g DockerComposeGenerator) getCloudsXmlPath(cloudType CloudType) string {
	cloudDir := g.getWorkerConfigDir(cloudType)
	return path.Join(cloudDir, cloudsXml)
}

func (g DockerComposeGenerator) getWorkerPropertiesPath(cloudType CloudType) string {
	cloudDir := g.getWorkerConfigDir(cloudType)
	return path.Join(cloudDir, workerProperties)
}

func (g DockerComposeGenerator) getWorkerLoggingPropertiesPath(cloudType CloudType) string {
	cloudDir := g.getWorkerConfigDir(cloudType)
	return path.Join(cloudDir, log4jProperties)
}

func (g DockerComposeGenerator) getStorageProperties() map[string] string {
	properties := make(map[string] string)
	properties["perspective.storage.hosts"] = "storage:5801"
	return properties
}

func (g DockerComposeGenerator) getLoggingProperties(serviceName string) map[string] string {
	properties := make(map[string] string)
	properties["log4j.rootLogger"] = "WARN, logfile"
	properties["log4j.appender.logfile"] = "org.apache.log4j.DailyRollingFileAppender"
	properties["log4j.appender.logfile.MaxFileSize"] = "2GB"
	properties["log4j.appender.logfile.File"] = path.Join(g.logsDir, "perspective-" + serviceName + ".log")
	properties["log4j.appender.logfile.bufferSize"] = "5242880"
	properties["log4j.appender.logfile.MaxBackupIndex"] = "7"
	properties["log4j.appender.logfile.layout"] = "org.apache.log4j.PatternLayout"
	properties["log4j.appender.logfile.layout.ConversionPattern"] = "%d [%25.25t] %-5p %-60.60c - %m%n"
	properties["log4j.logger.com.hazelcast "] = "INFO"
	properties["log4j.logger.org.meridor.perspective"] = "DEBUG"
	return properties
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

func (g DockerComposeGenerator) saveProperties(path string, properties map[string] string) {
	fmt.Printf("Saving [%s]...\n", path)
	var buffer bytes.Buffer
	for k, v := range properties {
		buffer.WriteString(fmt.Sprintf("%s=%s\n", k, v))
	}
	if (g.DryRun) {
		fmt.Println(buffer.String())
	} else {
		err := ioutil.WriteFile(path, buffer.Bytes(), fileMode)
		exitIfFailed(path, err)
	}
}

func (g DockerComposeGenerator) createDockerCompose(config ClusterConfig) DockerComposeYml {
	dockerComposeYml := DockerComposeYml{
		Version: "2.1",
		Services: make(map[string] ServiceConfig),
	}
	dockerComposeYml.Services[storage] = g.createStorageService(config)
	dockerComposeYml.Services[rest] = g.createRestService(config)
	for cloudType := range config.Clouds {
		workerServiceName := prepareCloudType(cloudType)
		dockerComposeYml.Services[workerServiceName] = g.createWorkerService(cloudType, config.Version)
	} 
	return dockerComposeYml
}

func (g DockerComposeGenerator) createEnvironment(serviceName string) map[string] interface{} {
	env := make(map[string] interface{})
	env["MISC_OPTS"] = fmt.Sprintf("-Xbootclasspath/a:/etc/perspective/%s", serviceName)
	env["LOGGING_OPTS"] = fmt.Sprintf("-Dlog4j.configuration=file:%s", path.Join(g.configDir, serviceName, log4jProperties))
	return env
}

func (g DockerComposeGenerator) createStorageService(config ClusterConfig) ServiceConfig {
	env := g.createEnvironment(storage)
	return ServiceConfig{
		ContainerName: "perspective-storage",
		Image: fmt.Sprintf("meridor/perspective-storage:%s", config.Version),
		Environment: env,
		Volumes: []string{volume(g.logsDir), readOnlyVolume(path.Join(g.configDir, storage))},
	}
}

func readOnlyVolume(dir string) string {
	return fmt.Sprintf("%s:ro", volume(dir))
}

func volume(dir string) string {
	return fmt.Sprintf("%s:%s", dir, dir)
}

func (g DockerComposeGenerator) createRestService(config ClusterConfig) ServiceConfig {
	return ServiceConfig{
		ContainerName: "perspective-rest",
		Image: fmt.Sprintf("meridor/perspective-rest:%s", config.Version),
		Ports: []string{fmt.Sprintf("8080:%d", config.ApiPort)},
		Environment: g.createEnvironment(rest),
		Links: []string{storage},
		DependsOn: []string{storage},
		Volumes: []string{volume(g.logsDir), readOnlyVolume(path.Join(g.configDir, rest))},
	}
}

func (g DockerComposeGenerator) createWorkerService(cloudType CloudType, version string) ServiceConfig {
	suffix := prepareCloudType(cloudType)
	volumes := []string{volume(g.logsDir), readOnlyVolume(path.Join(g.configDir, suffix))}
	if (cloudType == DOCKER) {
		volumes = append(volumes, volume("/var/run"))
	}
	return ServiceConfig{
		ContainerName: fmt.Sprintf("perspective-%s", suffix),
		Image: fmt.Sprintf("meridor/perspective-%s:%s", suffix, version),
		Environment: g.createEnvironment(suffix),
		Links: []string{storage},
		DependsOn: []string{rest},
		Volumes: volumes,
	}
}

func (g DockerComposeGenerator) saveDockerCompose(composeYml DockerComposeYml) {
	bts, err := yaml.Marshal(&composeYml)
	if (err != nil) {
		fmt.Printf("Failed to generate docker-compose.yml contents: %v\n", err)
	}
	ymlString := string(bts)
	ymlPath := createComposeYmlPath(g.Dir)
	fmt.Printf("Saving [%s]...\n", ymlPath)
	if (g.DryRun) {
		fmt.Println(ymlString)
	} else {
		err := ioutil.WriteFile(ymlPath, bts, fileMode)
		exitIfFailed(ymlPath, err)
	}
}

func createComposeYmlPath(dir string) string {
	return path.Join(dir, "docker-compose.yml")
}