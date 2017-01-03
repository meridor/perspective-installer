package wizard

import (
	"fmt"
	. "github.com/meridor/perspective-installer/data"
)

type DigitalOceanWizard struct {
}

func (w DigitalOceanWizard) Run() Cloud {
	fmt.Println("Setting up Digital Ocean worker.")
	cloud := Cloud{}
	xmlConfig := createCloudsXml()
	clouds := []CloudsXmlCloud{}
	clouds = append(clouds, addDigitalOceanCloud())
	for YesNoQuestion("Add one more account?", false) {
		clouds = append(clouds, addDigitalOceanCloud())
	}
	xmlConfig.Clouds = clouds
	cloud.XmlConfig = xmlConfig
	return cloud
}

func addDigitalOceanCloud() CloudsXmlCloud {
	cloud := CloudsXmlCloud{
		Endpoint: "unused",
		Identity: "unused",
		Enabled: true,
	}
	cloud.Credential = FreeInputQuestion("Enter Digital Ocean token:", "")
	cloud.Name = FreeInputQuestion(
		"Specify project name to use in Perspective:",
		fmt.Sprintf("digital_ocean_%s", cloud.Credential[:2]), //Using first symbols of token
	)
	cloud.Id = cloud.Name
	return cloud
}
