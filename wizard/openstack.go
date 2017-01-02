package wizard

import (
	. "github.com/meridor/perspective-installer/data"
	"fmt"
)

type OpenstackWizard struct {
	
}

func (w OpenstackWizard) Run() Cloud {
	fmt.Println("Setting up Openstack worker.")
	cloud := Cloud{}
	xmlConfig := createCloudsXml()
	clouds := []CloudsXmlCloud{}
	clouds = append(clouds, addCloud())
	for YesNoQuestion("Add one more project?", false) {
		clouds = append(clouds, addCloud())
	}
	xmlConfig.Clouds = clouds
	cloud.XmlConfig = xmlConfig
	return cloud
}

func addCloud() CloudsXmlCloud {
	cloud := CloudsXmlCloud{Enabled: true}
	cloud.Endpoint = FreeInputQuestion("Enter Openstack API URL:", "")
	projectName := FreeInputQuestion("Enter Openstack project name:", "")
	username := FreeInputQuestion("Enter Openstack user name:", "")
	cloud.Identity = fmt.Sprintf("%s:%s", projectName, username)
	cloud.Credential = FreeInputQuestion("Enter Openstack password:", "")
	cloud.Name = FreeInputQuestion("Specify project name to use in Perspective:", projectName)
	cloud.Id = cloud.Name
	return cloud
}