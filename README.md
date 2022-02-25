# scheduled-codeship-build

This lambda is intended to trigger a project build on Codeship.
It makes use of the *codeship-go* API client.

** Make sure to give permissions to the Codeship user to trigger builds
on the targeted projects

## Required Environment Variables (for the lambda config)

* CS_ORGANIZATION // the Codeship organization
* CS_PASSWORD // the Codeship password
* CS_USERNAME // the Codeship username for authentication


* CS1_PROJECT_UUID // the uuid of the first target Codeship project
* CS1_BUILD_REFERENCE // the desired reference for the build, e.g. "20.04" or "develop"
(note: "heads/" will be prepended to this value, if it is not already there)


* CS2_PROJECT_UUID // the uuid of the second target Codeship project
* CS2_BUILD_REFERENCE // the desired reference for the build, e.g. "20.04" or "develop"
(note: "heads/" will be prepended to this value, if it is not already there)

## Environment variables for deploying via Codeship
See the `aws.env.example` file for the full set of environment variables.

You will need to create an `aws.env` file with the proper values
and then an `aws.env.encrypted` file as follows.

Save the AES key from the codeship project in a `codeship.aes` file and 
then run the following command.

`jet encrypt aws.env aws.env.encrypted --key-path ./codeship.aes`

## Testing against Codeship

Add the real environment variable entries to your local.env file.
Then do `make debug`, comment out the `t.Skip` line in app/cron/builder/main_test.go in the `TestHandler` function. 
Then, in that directory, run `go test ./...`

Note: that will make one of the other tests fail because of the presence of all the environment variables.
You can just ignore that.

## Credential Rotation

### AWS Serverless User

1. Copy the aes key from Codeship
2. Paste it in a new file `codeship.aes`
3. Run `jet decrypt aws.env.encrypted aws.env`
4. (Optional) Compare the key in `aws.env` with the key in the most recent Terraform Cloud output
5. Use the Terraform CLI to taint the old access key
6. Run a new plan on Terraform Cloud
7. Review the new plan and apply if it is correct
8. Copy the new key and secret from the Terraform output into the aws.env file, overwriting the old values
9. Run `jet encrypt aws.env aws.env.encrypted`
10. Commit the new `aws.env.encrypted` file on the `develop` branch and push it to Github
11. Submit a PR to release the change to the `main` branch

### Codeship

(TBD)
