package main

import (
	"bufio"
	"os"
	"fmt"
)

var (
	answers = make(Answers)
	reader = &defaultReader{}
	supportedClouds = make(map[int] string)
)

type Answers map[int] interface{}

// Answer keys
const (
	CLOUDS = iota
)

func init() {
	for index, name := range []string{DIGITAL_OCEAN.String(), DOCKER.String(), OPENSTACK.String()} {
		supportedClouds[index] = name
	}
}

type Reader interface {
	Read() string
}

type defaultReader struct {
	
}

type AccessInfo struct {
	Api string
	Identity string
	Credential string
}

func (r *defaultReader) Read() string {
	reader := bufio.NewReader(os.Stdin)
	text, _ := reader.ReadString('\n')
	return text
}

func RunWizard() Answers {
	selectDesiredClouds()
	return answers
}

func selectDesiredClouds(){
	clouds := []int{}
	for key, name := range supportedClouds {
		if yesNoQuestion(fmt.Sprintf("Will you use %s cloud?", name), false) {
			fmt.Printf("Using %s.\n", name)
			clouds = append(clouds, key)
		}
	}
	answers[CLOUDS] = clouds
}

//func specifyAccessInfo(accessInfo *AccessInfo) {
//	if () {
//		
//	}
//}

func yesNoQuestion(message string, defaultAnswer bool) bool {
	fmt.Printf("%s [%s]\n", message, boolToString(defaultAnswer))
	answer := reader.Read()
	return isYesAnswer(answer)
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