package utils

import (
	"fmt"
	"net/url"

	"github.com/google/uuid"
)

// GetMicrosoftOAuthURL builds the Microsoft OAuth login URL
func GetMicrosoftOAuthURL(clientID, redirectURI, tenantID string) string {
	params := url.Values{}
	params.Add("client_id", clientID)
	params.Add("response_type", "code")
	params.Add("redirect_uri", redirectURI)
	params.Add("response_mode", "query")
	params.Add("scope", "User.Read openid email profile")
	params.Add("state", "xyz") // You should generate a real state value

	return fmt.Sprintf("https://login.microsoftonline.com/%s/oauth2/v2.0/authorize?%s", tenantID, params.Encode())
}

// GenerateSessionID generates a unique session ID
func GenerateSessionID() string {
	return uuid.New().String()
}
