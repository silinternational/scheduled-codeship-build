# scheduled-codeship-build

This lambda is intended to trigger a project build on Codeship.
It makes use of the *codeship-go* API client.

** Make sure to give permissions to the Codeship user to trigger builds
on the targeted projects

# Environment Variables

## Configuration for Lambda function

* CS_ORGANIZATION // the Codeship organization
* CS_PASSWORD // the Codeship password
* CS_USERNAME // the Codeship username for authentication


* CS1_PROJECT_UUID // the uuid of the first target Codeship project
* CS1_BUILD_REFERENCE // the desired reference for the build, e.g. "20.04" or "develop"
(note: "heads/" will be prepended to this value, if it is not already there)


* CS2_PROJECT_UUID // the uuid of the second target Codeship project
* CS2_BUILD_REFERENCE // the desired reference for the build, e.g. "20.04" or "develop"
(note: "heads/" will be prepended to this value, if it is not already there)

## Deploy to AWS

The included Terraform config generates a "Serverless User" with sufficient permission to
deploy the Lambda function. The access key and secret are outputs of this config.

* AWS_ACCESS_KEY_ID
* AWS_SECRET_ACCESS_KEY

## 1Password Secrets Automation

Secrets and other configuration are stored in a 1Password vault. During deployment, the
1Password CLI (`op`) pulls the values into the correct environment variables. This requires
two environment variables:

* OP_CONNECT_HOST
* OP_CONNECT_TOKEN

To configure these variables:

1. Create a file `op.env` with these the OP variables above and their correct values.
2. Save the AES key from the codeship project in a `cloudbees.aes` file.
3. Run the following command and commit the `op.env.encrypted` file to the git repo.

```sh
jet encrypt op.env op.env.encrypted
```

## Manual testing

1. Run `make debug` to start a container in Docker Compose
2. At the container shell prompt, run `eval $(op signin)`
3. Enter your 1Password credentials at the prompts. Your email address and secret key can be found at https://my.1password.com/profile. You may need to replace the "my" in that URL with your own organization name.
4. You can now run any serverless command prefixed with `op run` and 1Password will insert the correct credentials from your vault. For example: `op run sls info` will show the status of the Lambda on your AWS account.
 
To run the Go tests, comment out the `t.Skip` line in app/cron/builder/main_test.go in the `TestHandler` function. 
Then, in that directory, run `go test ./...`

Note: that will make one of the other tests fail because of the presence of all the environment variables.
You can just ignore that.

## Credential Rotation

### 1Password

1. Copy the aes key from Codeship
2. Paste it in a new file `cloudbees.aes`
3. Run `jet decrypt op.env.encrypted op.env`
4. Request a new token from a 1Password user who is an Owner or Administrator. (The URL for this is https://my.1password.com/integrations/active)
5. Paste the revised 1Password token into the OP_CONNECT_TOKEN line in op.env.
6. Run `jet encrypt op.env op.env.encrypted`
7. Commit the new `op.env.encrypted` file on the `develop` branch and push it to Github
8. Submit a PR to release the change to the `main` branch

### AWS Serverless User

1. Use the Terraform CLI to taint the old access key
2. Run a new plan on Terraform Cloud
3. Review and apply the new plan if it is correct
4. Copy the new key and secret from the Terraform output into 1Password. (Reference the codeship-services.yml file for the name of the vault and item to use.)

### Codeship

Codeship uses an ordinary username and password with HTTP Basic Authentication for API access. When the username
or password changes, simply update the correct fields in 1Password. Reference the codeship-services.yml file to find the name of the vault, item, and field to use.
