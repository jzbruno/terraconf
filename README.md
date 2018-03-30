# TerraConf [![travis-ci](https://travis-ci.org/jzbruno/terraconf.svg)](https://travis-ci.org/jzbruno/terraconf) [![codeclimate maintainability](https://api.codeclimate.com/v1/badges/a8355a1720309e1c63c2/maintainability)](https://codeclimate.com/github/jzbruno/terraconf) [![codeclimate test coverage](https://api.codeclimate.com/v1/badges/a8355a1720309e1c63c2/test_coverage)](https://codeclimate.com/github/jzbruno/terraconf) [![go report card](https://goreportcard.com/badge/github.com/jzbruno/terraconf)](https://goreportcard.com/report/github.com/jzbruno/terraconf) [![go doc](https://img.shields.io/badge/godoc-reference-blue.svg)](https://godoc.org/github.com/jzbruno/terraconf/pkg/terraconf)

Go package and cli for generating Terraform config from a Terraform state file.

# Usage

    cd cmd/terraconf
    go install
    terraconf <state-file>
