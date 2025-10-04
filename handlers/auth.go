package handlers

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net/http"
	"net/url"
	"tiktok-oauth2/config"
	"tiktok-oauth2/models"
	"tiktok-oauth2/utils"
)

// AuthHandler handles the initial OAuth authorization request
func AuthHandler(w http.ResponseWriter, r *http.Request) {
	// Generate random state for CSRF protection
	state, err := generateRandomState()
	if err != nil {
		config.DebugLog("❌ Failed to generate state parameter: %v", err)
		utils.WriteJSONResponse(w, http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Error:   "Failed to generate state parameter",
		})
		return
	}

	// Store state (in production, use Redis or database)
	// For now, we'll just validate it in callback
	_ = state

	// Build authorization URL
	authURL := buildAuthURL(state)
	config.DebugLog("🔗 Generated auth URL: %s", authURL)
	config.DebugLog("🛡️ Generated state: %s", state)

	// Redirect to TikTok OAuth page
	config.DebugLog("↗️ Redirecting to TikTok OAuth page")
	http.Redirect(w, r, authURL, http.StatusFound)
}

// RefreshTokenHandler handles token refresh requests
func RefreshTokenHandler(w http.ResponseWriter, r *http.Request) {
	// Get refresh token from request body
	var req struct {
		RefreshToken string `json:"refresh_token"`
	}

	if err := utils.ReadJSONResponse(&http.Response{Body: r.Body}, &req); err != nil {
		utils.WriteJSONResponse(w, http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error:   "Invalid request body",
		})
		return
	}

	if req.RefreshToken == "" {
		utils.WriteJSONResponse(w, http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error:   "Refresh token is required",
		})
		return
	}

	// Create HTTP client
	client := utils.NewHTTPClient("")

	// Prepare form data for token refresh
	formData := url.Values{}
	formData.Add("client_key", config.ClientKey)
	formData.Add("client_secret", config.ClientSecret)
	formData.Add("grant_type", "refresh_token")
	formData.Add("refresh_token", req.RefreshToken)

	// Make request to TikTok token endpoint
	resp, err := client.PostForm(config.TokenURL, formData)
	if err != nil {
		utils.WriteJSONResponse(w, http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Error:   "Failed to refresh token: " + err.Error(),
		})
		return
	}

	// Parse response
	var tokenResp models.TokenResponse
	if err := utils.ReadJSONResponse(resp, &tokenResp); err != nil {
		utils.WriteJSONResponse(w, http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Error:   "Failed to parse token response: " + err.Error(),
		})
		return
	}

	// Check if we got a valid access token
	if tokenResp.AccessToken == "" {
		utils.WriteJSONResponse(w, http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error:   "No access token received",
		})
		return
	}

	// Convert to TokenResponseData format
	tokenData := models.TokenResponseData{
		AccessToken:      tokenResp.AccessToken,
		ExpiresIn:        tokenResp.ExpiresIn,
		OpenID:           tokenResp.OpenID,
		RefreshToken:     tokenResp.RefreshToken,
		RefreshExpiresIn: tokenResp.RefreshExpiresIn,
	}

	// Return success response
	utils.WriteJSONResponse(w, http.StatusOK, models.APIResponse{
		Success: true,
		Message: "Token refreshed successfully",
		Data:    tokenData,
	})
}

// generateRandomState generates a random state string for CSRF protection
func generateRandomState() (string, error) {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

// buildAuthURL constructs the TikTok OAuth authorization URL
func buildAuthURL(state string) string {
	params := url.Values{}
	params.Add("client_key", config.ClientKey)
	params.Add("redirect_uri", config.RedirectURI)
	params.Add("response_type", "code")
	params.Add("scope", "user.info.basic,user.info.profile,user.info.stats,video.list,video.upload,video.publish")
	params.Add("state", state)

	authURL := fmt.Sprintf("%s?%s", config.AuthURL, params.Encode())

	config.DebugLog("🔧 Building auth URL with params:")
	config.DebugLog("  - client_key: %s", config.ClientKey)
	config.DebugLog("  - redirect_uri: %s", config.RedirectURI)
	config.DebugLog("  - response_type: code")
	config.DebugLog("  - scope: user.info.basic,user.info.profile,user.info.stats,video.list,video.upload,video.publish")
	config.DebugLog("  - state: %s", state)
	config.DebugLog("  - Final URL: %s", authURL)

	return authURL
}
