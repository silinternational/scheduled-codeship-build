terraform {
  cloud {
    organization = "gtis"

    workspaces {
      name = "scheduled-codeship-build"
    }
  }
}

