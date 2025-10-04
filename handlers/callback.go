package handlers

import (
	"fmt"
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
	tokenData, err := exchangeCodeForToken(code)
	if err != nil {
		utils.WriteJSONResponse(w, http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Error:   "Failed to exchange code for token: " + err.Error(),
		})
		return
	}

	// Fetch user info using the access token
	userInfo, err := fetchUserInfo(tokenData.AccessToken)
	if err != nil {
		// Log error but don't fail the entire request
		// User can still get token and fetch user info separately
		fmt.Printf("Warning: Failed to fetch user info: %v\n", err)
		userInfo = &models.UserInfo{} // Empty user info
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
	resp, err := client.PostForm(config.TokenURL, formData)
	if err != nil {
		return nil, fmt.Errorf("token request failed: %w", err)
	}

	// Parse response
	var tokenResp models.TokenResponse
	if err := utils.ReadJSONResponse(resp, &tokenResp); err != nil {
		return nil, fmt.Errorf("failed to parse token response: %w", err)
	}

	// Check for API errors
	if tokenResp.ErrorCode != 0 {
		return nil, fmt.Errorf("TikTok API error: %s", tokenResp.Description)
	}

	return &tokenResp.Data, nil
}

// fetchUserInfo fetches user information from TikTok API
func fetchUserInfo(accessToken string) (*models.UserInfo, error) {
	// Create HTTP client
	client := utils.NewHTTPClient("https://open.tiktokapis.com")

	// Create request
	req, err := http.NewRequest("GET", "/v2/user/info/?fields=open_id,union_id,avatar_url,avatar_url_100,avatar_large_url,display_name,bio_description,profile_deep_link,is_verified,username,follower_count,following_count,likes_count,video_count", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set authorization header
	req.Header.Set("Authorization", "Bearer "+accessToken)

	// Make request
	resp, err := client.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}

	// Parse response
	var userResp models.UserInfoResponse
	if err := utils.ReadJSONResponse(resp, &userResp); err != nil {
		return nil, fmt.Errorf("failed to parse user info response: %w", err)
	}

	// Check for API errors
	if userResp.ErrorCode != 0 {
		return nil, fmt.Errorf("TikTok API error: %s", userResp.Description)
	}

	return &userResp.Data, nil
}
