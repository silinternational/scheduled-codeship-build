package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
)

// EnvKeyCSPassword is the environment variable for
// the Codeship password
const EnvKeyCSPassword = "CS_PASSWORD"

// EnvKeyCSUsername is the environment variable for
// the Codeship username
const EnvKeyCSUsername = "CS_USERNAME"

// ParamOrganization is the parameter name for the
// the Codeship organization
const ParamOrganization = "organization"

// ParamProjects is a list of Codeship project UUIDs and build references, as a list of json objects
// e.g.: [{"uuid":"26e97136-8265-4172-867d-3392c7b3c322","ref":"20.04"}]
// Each reference will be prepended with "heads/" if it is not already included.
const ParamProjects = "projects"

type ProjectConfig struct {
	UUID string `json:"uuid"`
	Ref  string `json:"ref"`
}

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
	if err := getRequiredString(EnvKeyCSPassword, &c.CSPassword); err != nil {
		return err
	}

	if err := getRequiredString(EnvKeyCSUsername, &c.CSUsername); err != nil {
		return err
	}

	cfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		return err
	}

	ssmClient := ssm.NewFromConfig(cfg)
	input := ssm.GetParametersByPathInput{
		Path:           aws.String("/scheduled-codeship-build"),
		WithDecryption: aws.Bool(true),
	}
	path, err := ssmClient.GetParametersByPath(context.Background(), &input)
	if err != nil {
		return err
	}
	params := path.Parameters
	if c.appConfigParams == nil {
		c.appConfigParams = map[string]string{}
	}
	for _, p := range params {
		if p.Name == nil || p.Value == nil {
			continue
		}

		s := strings.Split(*p.Name, "/")
		if len(s) != 3 {
			continue
		}
		name := s[2]

		c.appConfigParams[name] = *p.Value
	}

	requiredParams := []string{
		ParamOrganization,
		ParamProjects,
	}

	for _, p := range requiredParams {
		if c.appConfigParams[p] == "" {
			return fmt.Errorf("required value missing for parameter '%s'", p)
		}
	}
	return nil
}

type BuilderConfig struct {
	CSPassword string `json:"CSPassword"`
	CSUsername string `json:"CSUsername"`

	appConfigParams map[string]string
}

func triggerBuild(ctx context.Context, project ProjectConfig, org *codeship.Organization) error {
	log.Printf("Building project %s on reference %s\n", project.UUID, project.Ref)

	headsPrefix := "heads/"
	buildRef := project.Ref
	if !strings.HasPrefix(buildRef, headsPrefix) {
		buildRef = headsPrefix + buildRef
	}
	success, resp, err := org.CreateBuild(ctx, project.UUID, buildRef, "")
	if err != nil {
		return fmt.Errorf("error triggering build on project with uuid: %s. %s",
			project.UUID, err)
	}

	if !success {
		return errors.New("failed to trigger build on project with uuid: " + project.UUID)
	}

	log.Printf("Response from codeship build call: %d %s\n", resp.StatusCode, resp.Status)
	return nil
}

func triggerBuilds(ctx context.Context, projects []ProjectConfig, org *codeship.Organization) error {
	for i := range projects {
		if err := triggerBuild(ctx, projects[i], org); err != nil {
			return err
		}
	}
	return nil
}

func handler() error {
	var config BuilderConfig
	if err := config.init(); err != nil {
		return err
	}

	auth := codeship.NewBasicAuth(config.CSUsername, config.CSPassword)
	client, err := codeship.New(auth)
	if err != nil {
		return errors.New("error creating api client: " + err.Error())
	}

	ctx := context.Background()

	organization := config.appConfigParams[ParamOrganization]
	org, err := client.Organization(ctx, organization)
	if err != nil {
		return errors.New("error scoping api client to organization: " + err.Error())
	}

	log.Print("Succeeded in authenticating with Codeship")

	projects := config.appConfigParams[ParamProjects]
	projectList, err := unmarshalProjectList(projects)
	if err != nil {
		return fmt.Errorf("unable to parse project configuration: %w", err)
	}

	if err = triggerBuilds(ctx, projectList, org); err != nil {
		return err
	}

	return nil
}

func main() {
	lambda.Start(handler)
}

func unmarshalProjectList(jsonList string) (config []ProjectConfig, err error) {
	err = json.Unmarshal([]byte(jsonList), &config)
	return
}
