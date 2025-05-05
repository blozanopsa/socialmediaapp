package controllers

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"

	"backend/configs"
	"backend/services"
	"backend/utils"
)

type UserController struct {
	UserService *services.UserService
}

func NewUserController(userService *services.UserService) *UserController {
	return &UserController{UserService: userService}
}

// MicrosoftLogin redirects the user to the Microsoft OAuth login page
func (uc *UserController) MicrosoftLogin(w http.ResponseWriter, r *http.Request) {
	cfg := configs.LoadConfig()
	oauthURL := utils.GetMicrosoftOAuthURL(cfg.OAuthClientID, cfg.OAuthRedirectURL, cfg.OAuthMicrosoftTenantID)
	http.Redirect(w, r, oauthURL, http.StatusFound)
}

// MicrosoftCallback handles the OAuth callback, exchanges code for tokens, stores user, and returns credentials
func (uc *UserController) MicrosoftCallback(w http.ResponseWriter, r *http.Request) {
	cfg := configs.LoadConfig()
	code := r.URL.Query().Get("code")
	if code == "" {
		http.Error(w, "Missing code in callback", http.StatusBadRequest)
		return
	}

	tokenEndpoint := fmt.Sprintf("https://login.microsoftonline.com/%s/oauth2/v2.0/token", cfg.OAuthMicrosoftTenantID)

	data := url.Values{}
	data.Set("client_id", cfg.OAuthClientID)
	data.Set("scope", "User.Read openid email profile")
	data.Set("code", code)
	data.Set("redirect_uri", cfg.OAuthRedirectURL)
	data.Set("grant_type", "authorization_code")
	data.Set("client_secret", cfg.OAuthClientSecret)

	resp, err := http.PostForm(tokenEndpoint, data)
	if err != nil {
		log.Println("MicrosoftCallback: Error exchanging code for token:", err)
		http.Error(w, "Failed to exchange code for token", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		http.Error(w, string(body), http.StatusUnauthorized)
		return
	}

	// Parse token response
	tokenBody, err := io.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, "Failed to read token response", http.StatusInternalServerError)
		return
	}

	var tokenResp map[string]interface{}
	if err := json.Unmarshal(tokenBody, &tokenResp); err != nil {
		http.Error(w, "Failed to parse token response", http.StatusInternalServerError)
		return
	}

	accessToken, ok := tokenResp["access_token"].(string)
	if !ok || accessToken == "" {
		http.Error(w, "No access token in response", http.StatusInternalServerError)
		return
	}

	// Use access token to get user info from Microsoft Graph
	graphReq, err := http.NewRequest("GET", "https://graph.microsoft.com/v1.0/me", nil)
	if err != nil {
		http.Error(w, "Failed to create graph request", http.StatusInternalServerError)
		return
	}
	graphReq.Header.Set("Authorization", "Bearer "+accessToken)

	graphResp, err := http.DefaultClient.Do(graphReq)
	if err != nil {
		log.Println("MicrosoftCallback: Error fetching user info:", err)
		http.Error(w, "Failed to get user info", http.StatusInternalServerError)
		return
	}
	defer graphResp.Body.Close()

	if graphResp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(graphResp.Body)
		http.Error(w, string(body), http.StatusUnauthorized)
		return
	}

	userBody, err := io.ReadAll(graphResp.Body)
	if err != nil {
		http.Error(w, "Failed to read user info", http.StatusInternalServerError)
		return
	}

	var msUser struct {
		ID                string `json:"id"`
		DisplayName       string `json:"displayName"`
		Mail              string `json:"mail"`
		UserPrincipalName string `json:"userPrincipalName"`
	}
	if err := json.Unmarshal(userBody, &msUser); err != nil {
		http.Error(w, "Failed to parse user info", http.StatusInternalServerError)
		return
	}
	email := msUser.Mail
	if email == "" {
		email = msUser.UserPrincipalName
	}

	// Store or update user in DB
	log.Println("MicrosoftCallback: Received code:", code)
	log.Println("MicrosoftCallback: Exchanging code for token")
	user, err := uc.UserService.FindOrCreateByEmail(msUser.DisplayName, email, "microsoft")
	if err != nil {
		log.Println("MicrosoftCallback: Error storing user in database:", err)
		http.Error(w, "Failed to save user", http.StatusInternalServerError)
		return
	}
	log.Println("MicrosoftCallback: User stored successfully")

	// Generate a session ID and set it in the response header
	log.Println("MicrosoftCallback: Generating session ID")
	sessionID := utils.GenerateSessionID()
	if err := uc.UserService.SaveSessionID(user.ID, sessionID); err != nil {
		log.Println("MicrosoftCallback: Error saving session ID:", err)
		http.Error(w, "Failed to save session ID", http.StatusInternalServerError)
		return
	}
	log.Println("MicrosoftCallback: Session ID saved successfully")
	// Update Set-Cookie header to include SameSite=None and Secure attributes
	w.Header().Set("Set-Cookie", "sessionID="+sessionID+"; Path=/; HttpOnly; SameSite=None; Secure")

	// Redirect to the frontend's auth-callback route
	frontendRedirect := "http://localhost:5173/auth/microsoft/callback"
	http.Redirect(w, r, frontendRedirect, http.StatusFound)
}

// MicrosoftLogout handles logout by redirecting to Microsoft logout endpoint and then back to frontend
func (uc *UserController) MicrosoftLogout(w http.ResponseWriter, r *http.Request) {
	// Microsoft logout endpoint
	// Docs: https://learn.microsoft.com/en-us/azure/active-directory/develop/v2-protocols-oidc#single-sign-out
	// Example: https://login.microsoftonline.com/common/oauth2/v2.0/logout?post_logout_redirect_uri=http://localhost:5173/
	frontendRedirect := "http://localhost:5173/"
	logoutURL := fmt.Sprintf("https://login.microsoftonline.com/common/oauth2/v2.0/logout?post_logout_redirect_uri=%s", url.QueryEscape(frontendRedirect))
	http.Redirect(w, r, logoutURL, http.StatusFound)

	// Clear the sessionID cookie on logout
	http.SetCookie(w, &http.Cookie{
		Name:   "sessionID",
		Value:  "",
		Path:   "/",
		MaxAge: -1, // Expire immediately
	})
}

// FacebookLogin redirects the user to the Facebook OAuth login page
func (uc *UserController) FacebookLogin(w http.ResponseWriter, r *http.Request) {
	cfg := configs.LoadConfig()
	clientID := cfg.FacebookAppID
	redirectURI := "http://localhost:8080/auth/facebook/callback"
	fbLoginURL := fmt.Sprintf(
		"https://www.facebook.com/v15.0/dialog/oauth?client_id=%s&redirect_uri=%s&scope=email,public_profile",
		clientID,
		url.QueryEscape(redirectURI),
	)
	http.Redirect(w, r, fbLoginURL, http.StatusFound)
}

// FacebookCallback handles the OAuth callback, exchanges code for tokens, and fetches user info

func (uc *UserController) FacebookCallback(w http.ResponseWriter, r *http.Request) {
	cfg := configs.LoadConfig()
	clientID := cfg.FacebookAppID
	clientSecret := cfg.FacebookAppSecret
	redirectURI := "http://localhost:8080/auth/facebook/callback"

	code := r.URL.Query().Get("code")
	if code == "" {
		http.Error(w, "Authorization code not provided", http.StatusBadRequest)
		return
	}

	// Exchange the authorization code for an access token
	tokenURL := fmt.Sprintf(
		"https://graph.facebook.com/v15.0/oauth/access_token?client_id=%s&redirect_uri=%s&client_secret=%s&code=%s",
		clientID,
		url.QueryEscape(redirectURI),
		clientSecret,
		code,
	)
	resp, err := http.Get(tokenURL)
	if err != nil || resp.StatusCode != http.StatusOK {
		http.Error(w, "Failed to exchange code for token", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	var tokenResp struct {
		AccessToken string `json:"access_token"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		http.Error(w, "Failed to parse token response", http.StatusInternalServerError)
		return
	}

	// Fetch user info from Facebook Graph API
	userInfoURL := fmt.Sprintf("https://graph.facebook.com/me?fields=id,name,email&access_token=%s", tokenResp.AccessToken)
	userResp, err := http.Get(userInfoURL)
	if err != nil || userResp.StatusCode != http.StatusOK {
		http.Error(w, "Failed to fetch user info", http.StatusInternalServerError)
		return
	}
	defer userResp.Body.Close()

	var fbUser struct {
		ID    string `json:"id"`
		Name  string `json:"name"`
		Email string `json:"email"`
	}
	if err := json.NewDecoder(userResp.Body).Decode(&fbUser); err != nil {
		http.Error(w, "Failed to parse user info", http.StatusInternalServerError)
		return
	}

	// Store or update user in DB
	user, err := uc.UserService.FindOrCreateByEmail(fbUser.Name, fbUser.Email, "facebook")
	if err != nil {
		http.Error(w, "Failed to save user", http.StatusInternalServerError)
		return
	}

	// Generate a session ID and set it in the response header
	sessionID := utils.GenerateSessionID()
	uc.UserService.SaveSessionID(user.ID, sessionID)
	// Update Set-Cookie header to include SameSite=None and Secure attributes
	w.Header().Set("Set-Cookie", "sessionID="+sessionID+"; Path=/; HttpOnly; SameSite=None; Secure")

	// Redirect to the frontend's auth-callback route (same as Microsoft)
	frontendRedirect := "http://localhost:5173/auth/microsoft/callback"
	http.Redirect(w, r, frontendRedirect, http.StatusFound)
}

// Add a new method to fetch user data based on session ID
func (uc *UserController) GetUserData(w http.ResponseWriter, r *http.Request) {
	// Update to retrieve session ID from cookies
	sessionID, err := r.Cookie("sessionID")
	if err != nil {
		http.Error(w, "Missing session ID", http.StatusUnauthorized)
		return
	}

	user, err := uc.UserService.GetUserBySessionID(sessionID.Value)
	if err != nil {
		http.Error(w, "Invalid session ID", http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	// Return full user info (including ID) for frontend
	response := map[string]interface{}{
		"id":    user.ID,
		"name":  user.Name,
		"email": user.Email,
	}
	json.NewEncoder(w).Encode(response)
}

// Add a new endpoint to handle logout and clear the sessionID cookie
func (uc *UserController) Logout(w http.ResponseWriter, r *http.Request) {
	// Clear the sessionID cookie
	http.SetCookie(w, &http.Cookie{
		Name:   "sessionID",
		Value:  "",
		Path:   "/",
		MaxAge: -1, // Expire immediately
	})

	// Respond with a success message
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Logged out successfully"))
}
