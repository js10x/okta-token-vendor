package vendor

import (
	"fmt"
	"time"
)

type SessionTokenRequest struct {
	Username                  string `json:"username"`
	Password                  string `json:"password"`
	MultiOptionalFactorEnroll bool   `json:"multiOptionalFactorEnroll"`
	WarnBeforePasswordExpired bool   `json:"warnBeforePasswordExpired"`
}

type SessionTokenResponse struct {
	ExpiresAt time.Time `json:"expiresAt"`
	Status    string    `json:"status"`
	Token     string    `json:"sessionToken"`
}

type AuthorizationCodeResponse struct {
	CodeVerifier string
	Code         string
}

type AccessTokenResponse struct {
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
	AccessToken string `json:"access_token"`
	Scope       string `json:"scope"`
	IDToken     string `json:"id_token"`
}

func (t *AccessTokenResponse) ToString() string {
	return fmt.Sprintf("\nAccess Token: \n\nType: %v \n\nExpires In: %v \n\nAccess Token: %v \n\nScope: %v \n\n",
		t.TokenType, t.ExpiresIn, t.AccessToken, t.Scope)
}
