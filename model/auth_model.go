package model

// parameters used by Cognito User Pool
const (
	ParamUsername             = "USERNAME"
	ParamPassword             = "PASSWORD"
	ParamNewPassword          = "NEW_PASSWORD"
	ParamRefreshToken         = "REFRESH_TOKEN"
	ParamSMSMFACode           = "SMS_MFA_CODE"
	ParamSoftwareTokenMFACode = "SOFTWARE_TOKEN_MFA_CODE"
)

type AuthToken struct {
	TokenType    string `json:"tokenType,omitempty"`
	AccessToken  string `json:"accessToken,omitempty"`
	IDToken      string `json:"idToken,omitempty"`
	RefreshToken string `json:"refreshToken,omitempty"`
	ExpiresIn    int64  `json:"expiresIn,omitempty"`
}

type SignInRequest struct {
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
}

type SignInResponse struct {
	AuthToken
	Session             string            `json:"session,omitempty"`
	ChallengeName       string            `json:"challengeName,omitempty"`
	ChallengeParameters map[string]string `json:"challengeParameters,omitempty"`
}

type RespondToAuthChallengeRequest struct {
	ChallengeName string `json:"challengeName,omitempty"`
	Session       string `json:"session,omitempty"`
	Username      string `json:"username,omitempty"`
	Password      string `json:"password,omitempty"`
	MFACode       string `json:"mfaCode,omitempty"`
}

type RespondToAuthChallengeResponse struct {
	AuthToken
	Session             string            `json:"session,omitempty"`
	ChallengeName       string            `json:"challengeName,omitempty"`
	ChallengeParameters map[string]string `json:"challengeParameters,omitempty"`
}
