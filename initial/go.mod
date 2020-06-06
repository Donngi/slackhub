module github.com/nicoJN/chatops-slackhub/initial

go 1.14

require (
	github.com/aws/aws-lambda-go v1.17.0
	github.com/aws/aws-sdk-go v1.31.8
	github.com/cenkalti/backoff v2.2.1+incompatible // indirect
	github.com/gofrs/uuid v3.3.0+incompatible // indirect
	github.com/guregu/dynamo v1.8.0
	github.com/nicoJN/chatops-slackhub/auth v0.0.0
	github.com/nicoJN/chatops-slackhub/tool v0.0.0
	github.com/slack-go/slack v0.6.5
	golang.org/x/net v0.0.0-20200602114024-627f9648deb9 // indirect
)

replace github.com/nicoJN/chatops-slackhub/auth => ../auth

replace github.com/nicoJN/chatops-slackhub/tool => ../tool
