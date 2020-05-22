package service

import (
	"context"
	"eks-demo/model"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws/awserr"
	"github.com/aws/aws-sdk-go-v2/aws/external"
	cip "github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
)

func NewAuthService(cfg *AuthServiceConfig) *AuthService {

	awsCfg, err := external.LoadDefaultAWSConfig()
	if err != nil {
		panic("failed to load config, " + err.Error())
	}

	client := cip.New(awsCfg)

	return &AuthService{
		cfg:    cfg,
		client: client,
	}
}

type AuthServiceConfig struct {
	UserPoolClientID string
}

type AuthService struct {
	cfg    *AuthServiceConfig
	client *cip.Client
}

func (svc *AuthService) SignIn(param *model.SignInRequest) (*model.SignInResponse, error) {

	authParams := map[string]string{
		model.ParamUsername: param.Username,
		model.ParamPassword: param.Password,
	}

	input := &cip.InitiateAuthInput{
		AuthFlow:       cip.AuthFlowTypeUserPasswordAuth,
		ClientId:       &svc.cfg.UserPoolClientID,
		AuthParameters: authParams,
	}

	req := svc.client.InitiateAuthRequest(input)

	resp, err := req.Send(context.Background())
	if err != nil {
		return nil, err
	}

	var session string
	if resp.Session != nil {
		session = *resp.Session
	}

	var challengeName string
	if resp.ChallengeName != "" {
		cn, _ := resp.ChallengeName.MarshalValue()
		challengeName = cn
	}

	return &model.SignInResponse{
		AuthToken:     toAuthToken(resp.AuthenticationResult),
		Session:       session,
		ChallengeName: challengeName,
	}, nil
}

func (svc *AuthService) RespondToAuthChallenge(param *model.RespondToAuthChallengeRequest) (*model.RespondToAuthChallengeResponse, error) {

	var challengeName cip.ChallengeNameType

	challengeResponses := map[string]string{}

	if cn, _ := cip.ChallengeNameTypeSmsMfa.MarshalValue(); cn == param.ChallengeName {

		challengeName = cip.ChallengeNameTypeSmsMfa
		challengeResponses[model.ParamUsername] = param.Username
		challengeResponses[model.ParamSMSMFACode] = param.MFACode

	} else if cn, _ := cip.ChallengeNameTypeSoftwareTokenMfa.MarshalValue(); cn == param.ChallengeName {

		challengeName = cip.ChallengeNameTypeSoftwareTokenMfa
		challengeResponses[model.ParamUsername] = param.Username
		challengeResponses[model.ParamSoftwareTokenMFACode] = param.MFACode

	} else if cn, _ := cip.ChallengeNameTypeNewPasswordRequired.MarshalValue(); cn == param.ChallengeName {

		challengeName = cip.ChallengeNameTypeNewPasswordRequired
		challengeResponses[model.ParamUsername] = param.Username
		challengeResponses[model.ParamNewPassword] = param.Password

	} else {
		return nil, awserr.New(cip.ErrCodeInvalidParameterException, fmt.Sprintf("challenge name %s is invalid", param.ChallengeName), nil)
	}

	input := &cip.RespondToAuthChallengeInput{
		ClientId:           &svc.cfg.UserPoolClientID,
		Session:            &param.Session,
		ChallengeName:      challengeName,
		ChallengeResponses: challengeResponses,
	}

	req := svc.client.RespondToAuthChallengeRequest(input)

	resp, err := req.Send(context.Background())
	if err != nil {
		return nil, err
	}

	var session string
	if resp.Session != nil {
		session = *resp.Session
	}

	var nextChallengeName string
	if resp.ChallengeName != "" {
		cn, _ := resp.ChallengeName.MarshalValue()
		nextChallengeName = cn
	}

	return &model.RespondToAuthChallengeResponse{
		AuthToken:     toAuthToken(resp.AuthenticationResult),
		Session:       session,
		ChallengeName: nextChallengeName,
	}, nil
}

func toAuthToken(result *cip.AuthenticationResultType) model.AuthToken {

	token := model.AuthToken{}

	if result != nil {
		if result.TokenType != nil {
			token.TokenType = *result.TokenType
		}
		if result.AccessToken != nil {
			token.AccessToken = *result.AccessToken
		}
		if result.IdToken != nil {
			token.IDToken = *result.IdToken
		}
		if result.RefreshToken != nil {
			token.RefreshToken = *result.RefreshToken
		}
		if result.ExpiresIn != nil {
			token.ExpiresIn = *result.ExpiresIn
		}
	}
	return token
}
