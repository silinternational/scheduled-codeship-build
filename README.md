# scheduled-codeship-build

This lambda is intended to trigger a project build on Codeship.
It makes use of the *codeship-go* API client.

## Required Environment Variables (for the lambda config)

* CS_ORGANIZATION // the Codeship organization
* CS_PASSWORD // the Codeship password
* CS_USERNAME // the Codeship username for authentication
* CS_PROJECT_UUID // the uuid of the target Codeship project
* CS_BUILD_REFERENCE // the desired reference for the build, e.g. "20.04" or "develop"
(note: "heads/" will be prepended to this value, if it is not already there)

## Testing against Codeship

Add the real environment variable entries to your local.env file.
Then do `make debug`, comment out the `skip` line in app/cron/builder/main_test.go in the `TestHandler` function. 
Then, in that directory, run `go test ./...`

Note: that will make one of the other tests fail because of the presence of all the environment variables.
You can just ignore that.
