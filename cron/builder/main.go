package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/aws/aws-lambda-go/lambda"
	codeship "github.com/codeship/codeship-go"
)

// EnvKeyCSOrganization is the environment variable for
// the Codeship organization
const EnvKeyCSOrganization = "CS_ORGANIZATION"

// EnvKeyCSPassword is the environment variable for
// the Codeship password
const EnvKeyCSPassword = "CS_PASSWORD"

// EnvKeyCSUsername is the environment variable for
// the Codeship username
const EnvKeyCSUsername = "CS_USERNAME"

// EnvKeyCSProjectUUID is the environment variable for
// the uuid of the target Codeship project
const EnvKeyCSProjectUUID = "CS_PROJECT_UUID"

// EnvKeyCSBuildReference is the environment variable for
// the desired reference for the build of the Codeship project
// e.g. "20.04"  (which will have "heads/" tacked onto the beginning, if it isn't there)
const EnvKeyCSBuildReference = "CS_BUILD_REFERENCE"

func getRequiredString(envKey string, configEntry *string) error {
	if *configEntry != "" {
		return nil
	}

	value := os.Getenv(envKey)
	if value == "" {
		return fmt.Errorf("required value missing for environment variable %s", envKey)
	}
	*configEntry = value

	return nil
}

func (c *BuilderConfig) init() error {
	if err := getRequiredString(EnvKeyCSOrganization, &c.CSOrganization); err != nil {
		return err
	}

	if err := getRequiredString(EnvKeyCSPassword, &c.CSPassword); err != nil {
		return err
	}

	if err := getRequiredString(EnvKeyCSUsername, &c.CSUsername); err != nil {
		return err
	}

	if err := getRequiredString(EnvKeyCSProjectUUID, &c.CSProjectUUID); err != nil {
		return err
	}

	if err := getRequiredString(EnvKeyCSBuildReference, &c.CSBuildReference); err != nil {
		return err
	}

	return nil
}

type BuilderConfig struct {
	CSOrganization string `json:"CSOrganization"`
	CSPassword     string `json:"CSPassword"`
	CSUsername     string `json:"CSUsername"`

	CSBuildReference string `json:"CSBuildReference"`
	CSProjectUUID    string `json:"CSProjectUUID"`
}

func getProject(ctx context.Context, config BuilderConfig, org *codeship.Organization) (codeship.Project, error) {
	projList, _, err := org.ListProjects(ctx)
	if err != nil {
		return codeship.Project{}, errors.New("error getting project list: " + err.Error())
	}

	projCount := len(projList.Projects)
	fmt.Printf("Project count: %d for org: %s\n", projCount, org.Name)

	if projCount < 1 {
		return codeship.Project{}, errors.New("no projects found for org " + config.CSOrganization)
	}

	for _, p := range projList.Projects {
		if p.UUID == config.CSProjectUUID {
			return p, nil
		}
	}

	return codeship.Project{}, errors.New("failed to find Codeship project with uuid: " + config.CSProjectUUID)
}

func triggerBuild(ctx context.Context, config BuilderConfig, org *codeship.Organization) (codeship.Response, error) {
	headsPrefix := "heads/"
	buildRef := config.CSBuildReference
	if !strings.HasPrefix(buildRef, headsPrefix) {
		buildRef = headsPrefix + buildRef
	}
	success, resp, err := org.CreateBuild(ctx, config.CSProjectUUID, buildRef, "")

	if err != nil {
		return codeship.Response{},
			fmt.Errorf("error triggering build on project with uuid: %s. %s",
				config.CSProjectUUID, err)
	}

	if !success {
		return codeship.Response{},
			errors.New("failed to trigger build on project with uuid: " + config.CSProjectUUID)
	}
	return resp, nil
}

func handler(config BuilderConfig) error {
	auth := codeship.NewBasicAuth(config.CSUsername, config.CSPassword)
	client, err := codeship.New(auth)
	if err != nil {
		return errors.New("error creating api client: " + err.Error())
	}

	ctx := context.Background()

	org, err := client.Organization(ctx, config.CSOrganization)
	if err != nil {
		return errors.New("error scoping api client to organization: " + err.Error())
	}

	log.Print("Succeeded in authenticating with Codeship")

	resp, err := triggerBuild(ctx, config, org)
	if err != nil {
		return err
	}
	fmt.Printf("Response from codeship build call: %d %s\n", resp.StatusCode, resp.Status)

	return nil

}

func main() {
	lambda.Start(handler)
}
