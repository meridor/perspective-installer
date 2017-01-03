package wizard

import (
	"fmt"
	. "github.com/meridor/perspective-installer/data"
	"github.com/pborman/uuid"
)

type DockerWizard struct {
}

func (w DockerWizard) Run() Cloud {
	fmt.Println("Setting up Docker worker.")
	cloud := Cloud{}
	xmlConfig := createCloudsXml()
	clouds := []CloudsXmlCloud{}
	clouds = append(clouds, addDockerCloud())
	for YesNoQuestion("Add one more API?", false) {
		clouds = append(clouds, addDockerCloud())
	}
	xmlConfig.Clouds = clouds
	cloud.XmlConfig = xmlConfig
	return cloud
}

func addDockerCloud() CloudsXmlCloud {
	cloud := CloudsXmlCloud{
		Identity: "unused",
		Credential: "unused",
		Enabled: true,
	}
	//TODO: add credentials support 
	cloud.Endpoint = FreeInputQuestion("Enter Docker API URL:", "unix:///var/run/docker.sock")
	cloud.Name = FreeInputQuestion(
		"Specify project name to use in Perspective:", 
		fmt.Sprintf("docker_%s", uuid.New()[:5]),
	)
	cloud.Id = cloud.Name
	return cloud
}
