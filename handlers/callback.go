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

	// Return success response with token data
	utils.WriteJSONResponse(w, http.StatusOK, models.APIResponse{
		Success: true,
		Message: "Authentication successful",
		Data:    tokenData,
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
