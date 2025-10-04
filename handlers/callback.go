package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"tiktok-oauth2/config"
	"tiktok-oauth2/models"
	"tiktok-oauth2/utils"
)

// CallbackHandler handles the OAuth callback from TikTok
func CallbackHandler(w http.ResponseWriter, r *http.Request) {
	// Get query parameters
	query := r.URL.Query()
	code := query.Get("code")
	state := query.Get("state")
	errorParam := query.Get("error")
	errorDescription := query.Get("error_description")

	// Debug: Log all query parameters
	config.DebugLog("🔍 Callback received - Full URL: %s", r.URL.String())
	config.DebugLog("📝 Query parameters: %+v", query)
	config.DebugLog("🔑 Code: %s", code)
	config.DebugLog("🛡️ State: %s", state)
	config.DebugLog("❌ Error: %s", errorParam)
	config.DebugLog("📄 Error Description: %s", errorDescription)

	// Check for OAuth errors
	if errorParam != "" {
		utils.WriteJSONResponse(w, http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error:   fmt.Sprintf("OAuth error: %s - %s", errorParam, errorDescription),
		})
		return
	}

	// Validate required parameters
	if code == "" {
		utils.WriteJSONResponse(w, http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error:   "Authorization code not found",
		})
		return
	}

	// In production, validate state parameter against stored value
	// For now, we'll just check if it exists
	if state == "" {
		utils.WriteJSONResponse(w, http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error:   "State parameter missing",
		})
		return
	}

	// Exchange authorization code for access token
	config.DebugLog("🔄 Starting token exchange process...")
	tokenData, err := exchangeCodeForToken(code)
	if err != nil {
		config.DebugLog("❌ Token exchange error: %v", err)
		utils.WriteJSONResponse(w, http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Error:   "Failed to exchange code for token: " + err.Error(),
		})
		return
	}

	// Debug: Log token data
	config.DebugLog("✅ Token data received: %+v", tokenData)

	// Fetch user info using the access token
	config.DebugLog("👤 Fetching user info with access token: %s", tokenData.AccessToken)
	userInfo, err := FetchUserInfo(tokenData.AccessToken)
	if err != nil {
		// Log error but don't fail the entire request
		// User can still get token and fetch user info separately
		config.DebugLog("⚠️ Warning: Failed to fetch user info: %v", err)
		userInfo = &models.UserInfo{} // Empty user info
	} else {
		config.DebugLog("✅ User info received: %+v", userInfo)
	}

	// Create combined response
	authResponse := models.AuthResponse{
		Token:    *tokenData,
		UserInfo: *userInfo,
	}

	// Return success response with token and user data
	utils.WriteJSONResponse(w, http.StatusOK, models.APIResponse{
		Success: true,
		Message: "Authentication successful",
		Data:    authResponse,
	})
}

// exchangeCodeForToken exchanges authorization code for access token
func exchangeCodeForToken(code string) (*models.TokenResponseData, error) {
	// Create HTTP client
	client := utils.NewHTTPClient("")

	// Prepare form data
	formData := url.Values{}
	formData.Add("client_key", config.ClientKey)
	formData.Add("client_secret", config.ClientSecret)
	formData.Add("grant_type", "authorization_code")
	formData.Add("code", code)
	formData.Add("redirect_uri", config.RedirectURI)

	// Make request to TikTok token endpoint
	config.DebugLog("🔄 Making token request to: %s", config.TokenURL)
	config.DebugLog("📝 Form data: %+v", formData)
	config.DebugLog("🔑 Client Key: %s", config.ClientKey)
	config.DebugLog("🔐 Client Secret: %s", config.ClientSecret)
	config.DebugLog("🌐 Redirect URI: %s", config.RedirectURI)

	resp, err := client.PostForm(config.TokenURL, formData)
	if err != nil {
		config.DebugLog("❌ Token request failed: %v", err)
		return nil, fmt.Errorf("token request failed: %w", err)
	}

	// Debug: Log response status and body
	config.DebugLog("📊 Response status: %d", resp.StatusCode)
	config.DebugLog("📋 Response headers: %+v", resp.Header)

	// Read raw response body for debugging
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		config.DebugLog("❌ Failed to read response body: %v", err)
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// Log raw response
	config.DebugLog("📄 Raw response body: %s", string(bodyBytes))

	// Parse response
	var tokenResp models.TokenResponse
	if err := json.Unmarshal(bodyBytes, &tokenResp); err != nil {
		config.DebugLog("❌ Failed to parse token response: %v", err)
		return nil, fmt.Errorf("failed to parse token response: %w", err)
	}

	// Debug: Log parsed response
	config.DebugLog("📦 Parsed token response: %+v", tokenResp)

	// Check for API errors
	if tokenResp.ErrorCode != 0 {
		config.DebugLog("❌ TikTok API error: %d - %s", tokenResp.ErrorCode, tokenResp.Description)
		return nil, fmt.Errorf("TikTok API error: %s", tokenResp.Description)
	}

	return &tokenResp.Data, nil
}

// FetchUserInfo fetches user information from TikTok API
func FetchUserInfo(accessToken string) (*models.UserInfo, error) {
	// Create HTTP client
	client := utils.NewHTTPClient("https://open.tiktokapis.com")

	// Create request
	userInfoURL := "https://open.tiktokapis.com/v2/user/info/?fields=open_id,union_id,avatar_url,avatar_url_100,avatar_large_url,display_name,bio_description,profile_deep_link,is_verified,username,follower_count,following_count,likes_count,video_count"
	config.DebugLog("👤 Fetching user info from: %s", userInfoURL)

	req, err := http.NewRequest("GET", userInfoURL, nil)
	if err != nil {
		config.DebugLog("❌ Failed to create user info request: %v", err)
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set authorization header
	req.Header.Set("Authorization", "Bearer "+accessToken)
	config.DebugLog("🔑 Authorization header: Bearer %s", accessToken)

	// Make request
	resp, err := client.Client.Do(req)
	if err != nil {
		config.DebugLog("❌ User info request failed: %v", err)
		return nil, fmt.Errorf("request failed: %w", err)
	}

	// Debug: Log response
	config.DebugLog("📊 User info response status: %d", resp.StatusCode)
	config.DebugLog("📋 User info response headers: %+v", resp.Header)

	// Parse response
	var userResp models.UserInfoResponse
	if err := utils.ReadJSONResponse(resp, &userResp); err != nil {
		config.DebugLog("❌ Failed to parse user info response: %v", err)
		return nil, fmt.Errorf("failed to parse user info response: %w", err)
	}

	// Debug: Log parsed response
	config.DebugLog("📦 Parsed user info response: %+v", userResp)

	// Check for API errors
	if userResp.Error.Code != "ok" {
		config.DebugLog("❌ TikTok User Info API error: %s - %s", userResp.Error.Code, userResp.Error.Message)
		return nil, fmt.Errorf("TikTok API error: %s", userResp.Error.Message)
	}

	return &userResp.Data.User, nil
}
