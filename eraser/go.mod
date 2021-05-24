module github.com/Jimon-s/slackhub/eraser

go 1.16

require (
	github.com/Jimon-s/slackhub/auth v0.0.0
	github.com/Jimon-s/slackhub/tool v0.0.0
	github.com/aws/aws-lambda-go v1.24.0
	github.com/aws/aws-sdk-go v1.38.45 // indirect
	github.com/cenkalti/backoff v2.2.1+incompatible // indirect
	github.com/gofrs/uuid v4.0.0+incompatible // indirect
	github.com/guregu/dynamo v1.10.4 // indirect
	github.com/slack-go/slack v0.9.1
	golang.org/x/net v0.0.0-20210521195947-fe42d452be8f // indirect
)

replace github.com/Jimon-s/slackhub/auth => ../auth

replace github.com/Jimon-s/slackhub/tool => ../tool
