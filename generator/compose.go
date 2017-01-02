package generator

import (
	"fmt"
	. "github.com/meridor/perspective-installer/data"
	"github.com/meridor/perspective-installer/wizard"
)

type Config struct {
	Version  string                 `yaml:"version,omitempty"`
	Services map[string] ServiceConfig          `yaml:"services,omitempty"`
	Networks map[string]interface{} `yaml:"networks,omitempty"`
}

type ServiceConfig struct {
	ContainerName string               `yaml:"container_name,omitempty"`
	DependsOn     []string             `yaml:"depends_on,omitempty"`
	Environment   map[string] interface{} `yaml:"environment,omitempty"`
	Expose        []string             `yaml:"expose,omitempty"`
	Image         string               `yaml:"image,omitempty"`
	Labels        []string      `yaml:"labels,omitempty"`
	Links         []string `yaml:"links,omitempty"`
	Ports         []string             `yaml:"ports,omitempty"`
	Privileged    bool                 `yaml:"privileged,omitempty"`
	Volumes       []string        `yaml:"volumes,omitempty"`
}

type DockerComposeGenerator struct {
	
}

func (g DockerComposeGenerator) Name() string {
	return "docker-compose"
}

func (g DockerComposeGenerator) Config(clouds map[CloudType]Cloud) string {
	fmt.Println("To use Docker Compose answer a bit more questions:")
	configDir := wizard.FreeInputQuestion("Specify configuration directory on host machine:", "/etc/perspective")
	logsDir := wizard.FreeInputQuestion("Specify logs directory on host machine:", "/var/log/perspective")
	panic("Not implemented " + configDir + " " + logsDir)
}

func (g DockerComposeGenerator) Command(dir string) string {
	return fmt.Sprintf("docker-compose -f %s/docker-compose.yml up", dir);
}