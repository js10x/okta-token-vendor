package pkce

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"net/url"
	"strings"
)

// Parses and returns the formatted authorization URL used in an authorization grant flow
// to retrieve the session token.
func AuthURL(issuer string) string {
	iss, err := url.Parse(issuer)
	if err != nil {
		return ""
	}
	result := fmt.Sprintf("%v://%v/api/v1/authn", iss.Scheme, iss.Host)
	return result
}

// Parses and returns the formatted OAUTH URL used in an authorization grant flow with PKCE.
// E.g. "/authorize" and "/token".
func OAuth2URL(issuer string, endpointUri string) string {
	iss, err := url.Parse(issuer)
	if err != nil {
		return ""
	}
	indexOf := strings.LastIndex(iss.Path, "/") + 1
	result := fmt.Sprintf("%v://%v/oauth2/%v/v1/%v", iss.Scheme, iss.Host, iss.Path[indexOf:], endpointUri)
	return result
}

// Builds and returns the URL query parameters needed to get the authorization code.
// Also returns the generated code verifier used to compute the code challenge.
func AuthCodeQuery(clientID string, redirectUri string, sessionToken string) (string, string) {

	// According to RFC7636 [Section 4] [https://datatracker.ietf.org/doc/html/rfc7636#section-4]
	// The code verifier is a high-entropy cryptographic random URL-safe string with a recommended length of between 43 and 128 characters.
	code_verifier := base64UrlEncodedString(60)
	code_challenge := CodeChallenge(code_verifier)

	params := url.Values{}
	params.Add("client_id", clientID)
	params.Add("code_challenge_method", "S256")
	params.Add("code_challenge", code_challenge)
	params.Add("redirect_uri", redirectUri)
	params.Add("response_type", "code")
	params.Add("scope", "openid")
	params.Add("nonce", base64EncodedString(20))
	params.Add("state", base64EncodedString(20))
	params.Add("sessionToken", sessionToken)
	return code_verifier, "?" + params.Encode()
}

// Computes a code challenge based on PKCE standards, which dicates that the code challenge
// is a Base64 URL-encoded SHA-256 hash of the code verifier.
func CodeChallenge(verifier string) string {

	// Use the raw bytes of the verifier to compute the SHA256 hash of it
	sha256ByteArray := sha256.Sum256([]byte(verifier))

	// Using the SHA256 computed bytes, form the base64 URL encoded string we need.
	challenge := base64.URLEncoding.WithPadding(base64.NoPadding).EncodeToString(sha256ByteArray[:])
	return challenge
}

// Creates and returns a base64 string with standard encoding.
func base64EncodedString(size int) string {
	bytes := make([]byte, size)
	rand.Read(bytes)
	return base64.StdEncoding.EncodeToString(bytes)
}

// Creates and returns a base64 string with URL encoding.
func base64UrlEncodedString(size int) string {
	bytes := make([]byte, size)
	rand.Read(bytes)
	return base64.URLEncoding.WithPadding(base64.NoPadding).EncodeToString(bytes)
}
