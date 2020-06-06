package auth

import (
	"net/http"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ssm/ssmiface"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ssm"
	"github.com/slack-go/slack"
)

// Auth is struct of auth client.
type Auth struct {
	client ssmiface.SSMAPI
}

// New is a constractor of Client.
func New(region string) *Auth {
	s := new(Auth)
	s.client = ssm.New(session.New(), &aws.Config{
		Region: aws.String(region),
	})
	return s
}

// getParam returns a parameter which is loaded from SSM Parameter Store.
func (a *Auth) getParam(key string) (string, error) {

	param, err := a.client.GetParameter(&ssm.GetParameterInput{
		Name:           aws.String(key),
		WithDecryption: aws.Bool(true),
	})
	if err != nil {
		return "", err
	}

	return *param.Parameter.Value, nil
}

// Authorize conducts a series of slack authentication.
func (a *Auth) Authorize(request events.APIGatewayProxyRequest, secretKey string) error {

	sc, err := a.getParam(secretKey)
	if err != nil {
		return err
	}

	if err := verify(request, sc); err != nil {
		return err
	}

	return nil
}

// verify returns the result of slack signing secret verification.
func verify(request events.APIGatewayProxyRequest, sc string) error {
	body := request.Body
	header := http.Header{}
	for k, v := range request.Headers {
		header.Set(k, v)
	}

	sv, err := slack.NewSecretsVerifier(header, sc)
	if err != nil {
		return err
	}

	sv.Write([]byte(body))
	return sv.Ensure()
}

// Client returns the slack-go client
func (a *Auth) Client(tokenKey string) (*slack.Client, error) {
	token, err := a.getParam(tokenKey)
	if err != nil {
		return nil, err
	}
	return slack.New(token), nil
}
