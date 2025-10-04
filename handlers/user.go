package handlers

import (
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
	userInfo, err := FetchUserInfo(token)
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

// extractBearerToken extracts token from "Bearer TOKEN" format
func extractBearerToken(authHeader string) string {
	const bearerPrefix = "Bearer "
	if len(authHeader) > len(bearerPrefix) && authHeader[:len(bearerPrefix)] == bearerPrefix {
		return authHeader[len(bearerPrefix):]
	}
	return ""
}
