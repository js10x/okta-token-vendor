package vendor

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/js10x/okta-token-vendor/pkce"
)

type TokenVendor struct {
	Ops Options
}

func NewTokenVendor(options []Option) *TokenVendor {
	ops := GetDefaultOptions()
	for _, op := range options {
		if op != nil {
			op(&ops)
		}
	}
	return &TokenVendor{Ops: ops}
}

// 1.) Get the session token
func (t *TokenVendor) GetSessionToken(username string, password string) (*SessionTokenResponse, error) {

	postConfig := &SessionTokenRequest{
		Username:                  username,
		Password:                  password,
		MultiOptionalFactorEnroll: true,
		WarnBeforePasswordExpired: true,
	}
	byteContent, err := json.Marshal(postConfig)
	if err != nil {
		return nil, err
	}

	request, err := http.NewRequest(http.MethodPost, pkce.AuthURL(t.Ops.Issuer), bytes.NewBuffer(byteContent))
	if err != nil {
		return nil, err
	}
	request.Header.Add("Content-Type", "application/json; charset=utf-8")
	request.Header.Add("Content-Length", strconv.Itoa(len(byteContent)))

	response, err := t.Ops.Client.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	oe := checkResponseFromOkta(response)
	if oe != nil {
		return nil, oe
	}

	var tokenResponse SessionTokenResponse
	if err := json.NewDecoder(response.Body).Decode(&tokenResponse); err != nil {
		return nil, err
	}

	switch {
	case tokenResponse.Status == "LOCKED_OUT":
		return nil, fmt.Errorf("OKTA issuer is reporting LOCKED_OUT")

	case len(strings.TrimSpace(tokenResponse.Token)) == 0:
		return nil, fmt.Errorf("failed to retrieve the SESSION TOKEN")
	}
	return &tokenResponse, nil
}

// 2.) Get the authorization code using the session token
func (t *TokenVendor) GetAuthorizationCode(sessionToken string) (*AuthorizationCodeResponse, error) {

	var codeVerifier string
	codeVerifier, encodedParameters := pkce.AuthCodeQuery(t.Ops.ClientID, t.Ops.RedirectURI, sessionToken)
	authorizeUrl := pkce.OAuth2URL(t.Ops.Issuer, "authorize") + encodedParameters

	request, err := http.NewRequest(http.MethodGet, authorizeUrl, nil)
	if err != nil {
		return nil, err
	}

	response, err := t.Ops.Client.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	oktaErr := checkResponseFromOkta(response)
	if oktaErr != nil {
		return nil, oktaErr
	}

	var authorizationCode string
	switch response.StatusCode {

	// Redirect (302)
	case http.StatusFound:
		location := response.Header.Get("location")
		redirect, err := url.Parse(location)
		if err != nil {
			return nil, err
		}
		authorizationCode = redirect.Query().Get("code")

	// Status OK (200)
	case http.StatusOK:
		if response.Request != nil {
			authorizationCode = response.Request.URL.Query().Get("code")
		}

	default:
		return nil, fmt.Errorf("something unexpected occurred. Status Code [%v]", response.StatusCode)
	}

	if len(strings.TrimSpace(authorizationCode)) == 0 {
		return nil, fmt.Errorf("failed to retrieve the AUTHORIZATION CODE")
	}

	codeResponse := &AuthorizationCodeResponse{
		CodeVerifier: codeVerifier,
		Code:         authorizationCode,
	}
	return codeResponse, nil
}

// 3.) Get the access token using the authorization code and the code verifier generated in step 2.
func (t *TokenVendor) GetAccessToken(codeVerifier string, authorizationCode string) (*AccessTokenResponse, error) {

	payload := url.Values{}
	payload.Set("client_id", t.Ops.ClientID)
	payload.Set("redirect_uri", t.Ops.RedirectURI)
	payload.Set("code_verifier", codeVerifier)
	payload.Set("code", authorizationCode)
	payload.Set("grant_type", "authorization_code")

	request, err := http.NewRequest(http.MethodPost, pkce.OAuth2URL(t.Ops.Issuer, "token"), strings.NewReader(payload.Encode()))
	if err != nil {
		return nil, err
	}

	request.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	response, err := t.Ops.Client.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	oktaErr := checkResponseFromOkta(response)
	if oktaErr != nil {
		return nil, oktaErr
	}

	var tokenResponse AccessTokenResponse
	if err := json.NewDecoder(response.Body).Decode(&tokenResponse); err != nil {
		return nil, err
	}

	if len(strings.TrimSpace(tokenResponse.AccessToken)) == 0 {
		return nil, fmt.Errorf("failed to retrieve the ACCESS TOKEN")
	}

	if t.Ops.OnTokenReceived != nil {
		t.Ops.OnTokenReceived(tokenResponse.AccessToken)
	}
	return &tokenResponse, nil
}

// Checks for a special error sent from Okta and returns it, if present in the response.
func checkResponseFromOkta(response *http.Response) *OktaError {

	var oktaErr OktaError
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return &OktaError{
			ErrorSummary: fmt.Sprintf("failed to read response from Okta server [%v]", err.Error()),
		}
	}
	json.Unmarshal(body, &oktaErr)

	if len(strings.TrimSpace(oktaErr.ErrorCode)) != 0 {
		return &oktaErr
	}

	// Restore the buffer of the response body.
	response.Body = ioutil.NopCloser(bytes.NewBuffer(body))
	return nil
}
