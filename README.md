# TerraConf

## Overview

[![travis-ci](https://travis-ci.org/jzbruno/terraconf.svg)](https://travis-ci.org/jzbruno/terraconf) [![codeclimate maintainability](https://api.codeclimate.com/v1/badges/a8355a1720309e1c63c2/maintainability)](https://codeclimate.com/github/jzbruno/terraconf) [![codeclimate test coverage](https://api.codeclimate.com/v1/badges/a8355a1720309e1c63c2/test_coverage)](https://codeclimate.com/github/jzbruno/terraconf) [![go report card](https://goreportcard.com/badge/github.com/jzbruno/terraconf)](https://goreportcard.com/report/github.com/jzbruno/terraconf) [![go doc](https://img.shields.io/badge/godoc-reference-blue.svg)](https://godoc.org/github.com/jzbruno/terraconf/pkg/terraconf)

Go package and CLI for generating Terraform config from a Terraform state file. 

Terraform supports importing resources into it's managed state but it does not create the config
files required to manage that state. The *terraconf* CLI tool generates the matching Terraform config
for a Terraform state file.

## Install

1. You can either build the CLI from source or download it from a GitHub release.

    * To build from source

        ```bash
        git clone git@github.com:jzbruno/terraconf.git
        cd terraconf/
        go get
        go install
        ```
        &NewLine;

    * To download from GitHub

        ```bash
        curl -sL https://github.com/jzbruno/terraconf/releases/download/v0.5.0/terraconf -o terraconf
        ```
        &NewLine;

## Usage

After *terraconf* is installed run the command with the following syntax. The Terraform config will 
be output to standard out which can be redirected to a file.

```bash
terraconf /path/to/state/file > main.tf
```
&NewLine;

Note that at this time purely computed attributes are included in the output. These need to be 
removed or Terraform plan and apply will complain about un-supported attributes. In a future release 
it may be possible to have computed attributes removed automatically by referencing the reosurce 
schema.

## Example

The following is more thorough example of using *terraconf* to bring existing resources under
Terraform control. This example will use the AWS provider but *terraconf* should work with any
provider.

1. Create an instance that will be used during the example. **If you already have existing resources 
you can skip this step.**

    ```bash
    instanceID="$(aws --profile jzbruno-terraform --region us-east-1 ec2 run-instances --instance-type t3.nano --image-id ami-0ff8a91507f77f867 --tag-specifications 'ResourceType=instance,Tags=[{Key=Name,Value=test-terraconf-instance}]' | jq -r '.Instances[0].InstanceId')"
    ```
    &NewLine;

2. Import the existing resources into Terraform state. **WARNING: Importing state without config using
the *-allow-missing-config* flag is dangerous because it is easy to accidently delete resources.**

    ```bash
    echo 'provider "aws" { region = "us-east-1" profile = "jzbruno-terraform" }' > main.tf
    terraform init
    terraform import -allow-missing-config aws_instance.test_terraconf_instance $instanceID
    ```
    &NewLine;

3. Geenrate the Terraform config from the Terraform state file.

    ```bash
    terraconf terraform.tfstate > main.tf
    ```
    &NewLine;

4. Update the Terraform config to remove computed attributes. In a future release it may be possible
to have computed attributes removed automatically by referencing the reosurce schema.

    ```bash
    terraform plan
    ```
    &NewLine;

    For any errors like the following remove that attribute from *main.tf*

    ```
    Error: aws_instance.test_terraconf_instance: "arn": this field cannot be set
    Error: aws_instance.test_terraconf_instance: : invalid or unknown key: id
    ```
    &NewLine;

    Some further modifications may be required to avoid diff issues depending on the resource type.
