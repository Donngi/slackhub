# STEP1 Deploy SlackHub to AWS account
Hi! In this step, we deploy SlackHub to your AWS account.

First, you should clone repository to your local environment and go to the project root. 

## Preparation
To set up SlackHub, you need 

- Go Runtime
- AWS CLI
- Terraform

If you don't have these tools. Please install them according to below.

<details>
<summary>If you are new of Go</summary>

Please visit [Go Official Page](https://golang.org/dl/) and download & install Go.

If you use macOS, you can install Go by using home brew.

```
$ brew install go
```

</details>

<details>
<summary>If you are new of AWS CLI</summary>

Please visit [how to page](https://docs.aws.amazon.com/cli/latest/userguide/install-cliv2.html) and install cli.

If you use macOS, you can install Go by using home brew.

```
$ brew install awscli
```

</details>

<details>
<summary>If you are new of Terraform</summary>

Please visit [how to page](https://learn.hashicorp.com/tutorials/terraform/install-cli) and install terraform.

If you use macOS, you can install Go by using home brew.

```
$ brew install terraform
```

</details>

## Set config
Default region is ap-northeast-1. If you want to use other regions, please open `terraform/env/example/aws.tf`.

You can set AWS region to use in this file.

```
provider "aws" {
  region = "ap-northeast-1"
}
```

## Build and Deploy
You can build all source code and deploy them to aws account at a time. Please execute a following command.

```
$ make deploy
```

After that, all resources will be provisioned. Then you should execute a following command to get endpoint url of API Gateway. Please memo a result.

```
$ aws apigatewayv2 get-apis --query "Items[?Name=='SlackHubAPI'].ApiEndpoint"
```
