package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var expectedRepositories = []string{
	"busybox",
	"nginx:stable",
	"mesosphere/marathon-lb~/^v1/",
	"quay.io/coreos/awscli=master,latest,edge",
	"gcr.io/google-containers/hyperkube~/^v1\\.(9|10)\\./",
}

func TestLoadYAMLFile(t *testing.T) {
	assert := assert.New(t)

	yc, err := LoadYAMLFile("../fixtures/config/config.yaml")

	assert.NotNil(yc, "should load config from valid config file")

	assert.Nil(err, "should NOT give an error while loading valid config file")

	if yc != nil {
		assert.Equal(expectedRepositories, yc.Repositories)
	}
}

func TestLoadYAMLFile_Shared(t *testing.T) {
	assert := assert.New(t)

	yc, err := LoadYAMLFile("../fixtures/config/config.yaml.shared")

	assert.NotNil(yc, "should load config from valid config file shared with others")

	assert.Nil(err, "should NOT give an error while loading valid shared config file")

	if yc != nil {
		assert.Equal(expectedRepositories, yc.Repositories)
	}
}

func TestLoadYAMLFile_Invalid(t *testing.T) {
	assert := assert.New(t)

	yc, err := LoadYAMLFile("../fixtures/config/config.yaml.invalid")

	assert.Nil(yc, "should NOT load config from invalid config file")

	assert.NotNil(err, "should give an error while trying to load invalid config file")
}

func TestLoadYAMLFile_Irrelevant(t *testing.T) {
	assert := assert.New(t)

	yc, err := LoadYAMLFile("../fixtures/config/config.yaml.irrelevant")

	assert.Nil(yc, "should NOT load config from irrelevant config file")

	assert.NotNil(err, "should give an error while trying to load irrelevant config file")
}

func TestLoadYAMLFile_NonExisting(t *testing.T) {
	assert := assert.New(t)

	yc, err := LoadYAMLFile("/i/do/not/exist/sorry")

	assert.Nil(yc, "should NOT load config from non-existing config file")

	assert.NotNil(err, "should give an error while trying to load non-existing config file")
}
