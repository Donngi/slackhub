module github.com/nicoJN/slackhub/interactive

go 1.14

require (
	github.com/aws/aws-lambda-go v1.19.1
	github.com/aws/aws-sdk-go v1.34.10
	github.com/cenkalti/backoff v2.2.1+incompatible // indirect
	github.com/gofrs/uuid v3.3.0+incompatible // indirect
	github.com/guregu/dynamo v1.9.1 // indirect
	github.com/nicoJN/slackhub/auth v0.0.0
	github.com/nicoJN/slackhub/tool v0.0.0
	github.com/slack-go/slack v0.6.6
	golang.org/x/net v0.0.0-20200822124328-c89045814202 // indirect
)

replace github.com/nicoJN/slackhub/auth => ../auth

replace github.com/nicoJN/slackhub/tool => ../tool
