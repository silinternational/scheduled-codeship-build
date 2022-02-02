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
