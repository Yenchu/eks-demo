package service_test

import (
	"eks-demo/model"
	"eks-demo/service"
	"fmt"
	"os"
	"testing"
)

func newAuthService() *service.AuthService {

	cfg := &service.AuthServiceConfig{
		UserPoolClientID: os.Getenv("USER_POOL_CLIENT_ID"),
	}

	return service.NewAuthService(cfg)
}

func TestSignIn(t *testing.T) {

	username := "@gmail.com"
	password := ""

	req := &model.SignInRequest{
		Username: username,
		Password: password,
	}

	svc := newAuthService()

	resp, err := svc.SignIn(req)
	if err != nil {
		t.Fatalf("SignIn failed: %v", err)
	}

	if resp.AccessToken != "" {
		fmt.Printf("idToken:\n%s\n\n", resp.IDToken)
		fmt.Printf("accessToken:\n%s\n\n", resp.AccessToken)
		fmt.Printf("refreshToken:\n%s\\nn", resp.RefreshToken)
		fmt.Printf("expiresIn:\n%d\n\n", resp.ExpiresIn)
	} else {
		fmt.Printf("challengeName:\n%s\n\n", resp.ChallengeName)
		fmt.Printf("session:\n%s\n\n", resp.Session)
	}
}

func TestRespondToAuthChallenge(t *testing.T) {

	challengeName := ""
	session := ""
	username := ""
	password := ""

	req := &model.RespondToAuthChallengeRequest{
		ChallengeName: challengeName,
		Session: session,
		Username: username,
		Password: password,
	}

	svc := newAuthService()

	resp, err := svc.RespondToAuthChallenge(req)
	if err != nil {
		t.Fatalf("RespondToAuthChallenge failed: %v", err)
	}
	t.Logf("RespondToAuthChallenge result:\n%+v\n\n", resp)
}