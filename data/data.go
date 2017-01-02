//go:generate stringer -type=CloudType enums.go
package data

import "encoding/xml"

// Supported clouds
type CloudType int

const (
	DIGITAL_OCEAN CloudType = iota
	DOCKER
	OPENSTACK
)

type Cloud struct {
	XmlConfig CloudsXml
	Properties map[string] interface{}
}

type CloudsXml struct {
	XMLName  xml.Name     `xml:"clouds"`
	XmlNS    string       `xml:"xmlns,attr"`
	Clouds []CloudsXmlCloud `xml:"cloud"`
}

type CloudsXmlCloud struct {
	Id string `xml:"id"`
	Name string `xml:"name"`
	Endpoint string `xml:"endpoint"`
	Identity string `xml:"identity"`
	Credential string `xml:"credential"`
	Enabled bool `xml:"enabled"`
}