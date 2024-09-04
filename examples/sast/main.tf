terraform {
  required_providers {
    snyk = {
      source = "github.com/xulei-nl/snyk"
    }
  }
}

provider "snyk" {
  api_token = "6a3edd59-aff1-4eb0-8bba-ce8c40d9c8b7"
  endpoint  = "https://api.eu.snyk.io/rest"
}

resource "snyk_sast" "example" {
  data = {
    attributes = {
      sast_enabled = true,
    },
    id = "da69edc7-1a7d-4318-8658-8339540ce0b5",
  }
}
