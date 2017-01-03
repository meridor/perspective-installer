package generator

import (
	"testing"
	. "github.com/aandryashin/matchers"
	. "github.com/meridor/perspective-installer/data"
)

func TestPrepareCloudType(t *testing.T) {
	AssertThat(t, prepareCloudType(DOCKER), EqualTo{"docker"})
	AssertThat(t, prepareCloudType(DIGITAL_OCEAN), EqualTo{"digital-ocean"})
	AssertThat(t, prepareCloudType(OPENSTACK), EqualTo{"openstack"})
} 
