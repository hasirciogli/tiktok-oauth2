package handlers

import (
	"fmt"
	"net/http"
	"tiktok-oauth2/models"
	"tiktok-oauth2/utils"
)

// UserInfoHandler handles user info requests
func UserInfoHandler(w http.ResponseWriter, r *http.Request) {
	// Get access token from Authorization header
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		utils.WriteJSONResponse(w, http.StatusUnauthorized, models.APIResponse{
			Success: false,
			Error:   "Authorization header required",
		})
		return
	}

	// Extract token from "Bearer TOKEN" format
	token := extractBearerToken(authHeader)
	if token == "" {
		utils.WriteJSONResponse(w, http.StatusUnauthorized, models.APIResponse{
			Success: false,
			Error:   "Invalid authorization format. Use 'Bearer TOKEN'",
		})
		return
	}

	// Fetch user info from TikTok API
	userInfo, err := fetchUserInfoFromTikTok(token)
	if err != nil {
		utils.WriteJSONResponse(w, http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Error:   "Failed to fetch user info: " + err.Error(),
		})
		return
	}

	// Return success response
	utils.WriteJSONResponse(w, http.StatusOK, models.APIResponse{
		Success: true,
		Message: "User info retrieved successfully",
		Data:    userInfo,
	})
}

// fetchUserInfoFromTikTok fetches user information from TikTok API
func fetchUserInfoFromTikTok(accessToken string) (*models.UserInfo, error) {
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

// extractBearerToken extracts token from "Bearer TOKEN" format
func extractBearerToken(authHeader string) string {
	const bearerPrefix = "Bearer "
	if len(authHeader) > len(bearerPrefix) && authHeader[:len(bearerPrefix)] == bearerPrefix {
		return authHeader[len(bearerPrefix):]
	}
	return ""
}
