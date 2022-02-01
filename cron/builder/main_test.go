package main

import (
	"testing"
)

func TestBuilderConfigMissingValue(t *testing.T) {
	config := BuilderConfig{
		CSOrganization: "OurOrg",
		CSUsername:     "MyName",
	}

	err := config.init()
	if err == nil {
		t.Error("Expected an error, but did not got one")
		return
	}

	want := "required value missing for environment variable " + EnvKeyCSPassword
	if err.Error() != want {
		t.Errorf(`incorrect error. Expected "%s"", but got "%s"`, want, err.Error())
		return
	}
}

func TestBuilderConfigOK(t *testing.T) {
	config := BuilderConfig{
		CSOrganization:   "OurOrg",
		CSPassword:       "MyPass",
		CSUsername:       "MyName",
		CSBuildReference: "develop",
		CSProjectUUID:    "abcd1234-abcd-1234-dcba-abcd1234dcba",
	}

	if err := config.init(); err != nil {
		t.Errorf("Did not expect an error, but got one: %v", err.Error())
		return
	}
}

// To run this test locally, first ensure you have all the required
// environment variables set and then comment out the t.Skip line
func TestHandler(t *testing.T) {
	t.Skip("Only run this in local development")

	// Just initialize config from .env file
	config := BuilderConfig{}

	if err := config.init(); err != nil {
		t.Errorf("failed initializing config for test: %v", err)
		return
	}

	err := handler(config)
	if err != nil {
		t.Errorf("error getting results: %v", err)
		return
	}
}
