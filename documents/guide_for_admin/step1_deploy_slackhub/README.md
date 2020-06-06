# STEP1 Deploy SlackHub to AWS account
Hi! In this step, we deploy SlackHub to your AWS account.

First, you should clone repository to your local environment and go to the project root. 

## Preparation
To set up SlackHub, you need 

- Go Runtime
- Java & Maven
- AWS CLI
- AWS CDK

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
<summary>If you are new of Java & Maven</summary>

**Java**

Please visit [OpenJDK](https://openjdk.java.net/) and download & install Java.

If you use macOS, you can install Java by using home brew.

```
$ brew install openjdk
```

**Maven**

Please visit [Apache Maven Project](https://maven.apache.org/download.cgi) and download & install Maven.

If you use macOS, you can install Maven by using home brew.

```
$ brew install maven
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
<summary>If you are new of AWS CDK</summary>

Please install cdk by npm command.

```
$ npm install -g aws-cdk
```

</details>

## Set properties
Default region is ap-northeast-1. If you want to use other regions, please open `config.properies`.

You can set AWS region to use in this file.

```
# region
# NOTE: If you want to use other regions, please replace here!
region = ap-northeast-1
```

## Build and Deploy
You can build all source code and deploy them to aws account at a time. Please the following command.

```
& make deploy
```

If you use aws-cli's profile, you can use the following command instead.

```
& make deploy OPT="--profile YOU_PROFILE_HERE"
```

After that, you'll get a message like below. You should memo these endpoint urls.

```
Outputs:
SlackHubStack.SlackHubInteractiveEndpointXXXXXXXX = https://XXXXXXXXXX.execute-api.your-region.amazonaws.com/prod/
SlackHubStack.SlackHubInitialEndpointXXXXXXXXXX = https://XXXXXXXXXX.execute-api.your-region.amazonaws.com/prod/
```

NOTE: If you forget these urls, you should re-execute `make deploy`. (Of course, you can also see them in the AWS console.)
