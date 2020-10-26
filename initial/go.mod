module github.com/nicoJN/slackhub/initial

go 1.14

require (
	github.com/aws/aws-lambda-go v1.19.1
	github.com/aws/aws-sdk-go v1.35.15
	github.com/cenkalti/backoff v2.2.1+incompatible // indirect
	github.com/gofrs/uuid v3.3.0+incompatible // indirect
	github.com/guregu/dynamo v1.9.1
	github.com/nicoJN/slackhub/auth v0.0.0
	github.com/nicoJN/slackhub/tool v0.0.0
	github.com/slack-go/slack v0.7.2
	golang.org/x/net v0.0.0-20201026091529-146b70c837a4 // indirect
)

replace github.com/nicoJN/slackhub/auth => ../auth

replace github.com/nicoJN/slackhub/tool => ../tool
