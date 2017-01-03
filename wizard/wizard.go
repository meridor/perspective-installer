package wizard

import (
	"bufio"
	"fmt"
	. "github.com/meridor/perspective-installer/data"
	"os"
	"strings"
	"strconv"
)

var (
	reader  Reader = defaultReader{}
	wizards        = make(map[CloudType]Wizard)
)

type Reader interface {
	Read() string
}

type defaultReader struct {
}

func (r defaultReader) Read() string {
	reader := bufio.NewReader(os.Stdin)
	text, _ := reader.ReadString('\n')
	return strings.Replace(text, "\n", "", -1)
}

type Wizard interface {
	Run() Cloud
}

const (
	defaultApiPort = "8080"
)

func RunWizards() (ClusterConfig, map[CloudType]Cloud) {
	initWizards()
	printWelcomeMessages()
	latestRelease := loadLatestRelease()
	clusterConfig := ClusterConfig{}
	clusterConfig.Version = FreeInputQuestion("Enter desired Perspective version:", latestRelease)
	clusterConfig.ApiPort = enterApiPort()
	clouds := make(map[CloudType]Cloud)
	for cloudType, wizard := range wizards {
		cloudName := cloudType.String()
		if YesNoQuestion(fmt.Sprintf("Are you going to use %s cloud?", cloudName), false) {
			cloud := wizard.Run()
			clouds[cloudType] = cloud
			fmt.Printf("Setting up %s cloud done.\n", cloudName)
			fmt.Println()
			fmt.Println()
		}
	}
	return clusterConfig, clouds
}

func printWelcomeMessages() {
	fmt.Println("Welcome to Perspective Installer!")
	fmt.Println("I will now ask you a set of questions. Default answers are shown in [square brackets].")
	fmt.Println("To abort this wizard type Ctrl+C.")
}

func loadLatestRelease() string {
	fmt.Println("Loading latest Perspective version...")
	latestRelease := GetLatestRelease()
	if (latestRelease == "") {
		fmt.Println("Failed to load latest version.")
	} else {
		fmt.Printf("Latest version is %s.\n", latestRelease)
	}
	return latestRelease
}

func enterApiPort() int {
	port := FreeInputQuestion("Enter API listen port:", defaultApiPort)
	intPort, err := strconv.Atoi(port)
	if (err != nil) {
		fmt.Printf("Not a number: %s. Using default port - %s.\n", port, defaultApiPort)
		return 8080
	}
	return intPort

}

func initWizards() {
	wizards[DIGITAL_OCEAN] = DigitalOceanWizard{}
	wizards[DOCKER] = DockerWizard{}
	wizards[OPENSTACK] = OpenstackWizard{}
}

func YesNoQuestion(message string, defaultAnswer bool) bool {
	printMessageWithDefaultAnswer(message, boolToString(defaultAnswer))
	answer := waitForAnswer(
		func() string {
			return boolToString(YesNoQuestion(message, defaultAnswer))
		},
		boolToString(defaultAnswer),
	)
	return isYesAnswer(answer)
}

func FreeInputQuestion(message string, defaultAnswer string) string {
	printMessageWithDefaultAnswer(message, defaultAnswer)
	return waitForAnswer(
		func() string {
			return FreeInputQuestion(message, defaultAnswer)
		},
		defaultAnswer,
	)
}

func createCloudsXml() CloudsXml {
	cloudsXml := CloudsXml{
		XmlNS: "urn:config.perspective.meridor.org",
	}
	return cloudsXml
}

func printMessageWithDefaultAnswer(message string, defaultAnswer string) {
	if defaultAnswer != "" {
		fmt.Printf("%s [%s]\n", message, defaultAnswer)
	} else {
		fmt.Println(message)
	}
}

func waitForAnswer(retryAction func() string, defaultAnswer string) string {
	answer := reader.Read()
	if answer == "" {
		if defaultAnswer == "" {
			return retryAction()
		}
		return defaultAnswer
	}
	return answer
}

func boolToString(b bool) string {
	if b {
		return "y"
	}
	return "n"
}

func isYesAnswer(input string) bool {
	return "y" == input || "Y" == input
}
