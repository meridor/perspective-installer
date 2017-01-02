package wizard

import (
	"bufio"
	"os"
	"fmt"
	. "github.com/meridor/perspective-installer/data"
	"strings"
)

var (
	reader Reader = defaultReader{}
	wizards = make(map[CloudType] Wizard)
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

func RunWizards() map[CloudType] Cloud {
	initWizards()
	fmt.Println("Welcome to Perspective installer!")
	fmt.Println("I will now ask you a set of questions. Default answers are shown in [square brackets].")
	fmt.Println("To abort this wizard type Ctrl+C.")
	clouds := make(map[CloudType]Cloud)
	for cloudType, wizard := range wizards {
		cloudName := cloudType.String()
		if YesNoQuestion(fmt.Sprintf("Are you going to use %s cloud?", cloudName), false) {
			cloud := wizard.Run()
			clouds[cloudType] = cloud
			fmt.Printf("Setting up %s cloud done.\n", cloudName)
		}
	}
	return clouds
}

func initWizards() {
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
		XmlNS:    "urn:config.perspective.meridor.org",
	}
	return cloudsXml
}

func printMessageWithDefaultAnswer(message string, defaultAnswer string) {
	if (defaultAnswer != "") {
		fmt.Printf("%s [%s]\n", message, defaultAnswer)
	} else {
		fmt.Println(message)
	}
}

func waitForAnswer(retryAction func() string, defaultAnswer string) string {
	answer := reader.Read()
	if (answer == "") {
		if (defaultAnswer == "") {
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
	return "y" == input || "Y" == input;
}

