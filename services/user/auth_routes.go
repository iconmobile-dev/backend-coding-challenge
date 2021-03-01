package user

import (
	"encoding/json"
	"net/http"

	"github.com/iconmobile-dev/backend-coding-challenge/lib/handlers"
)

type loginRequest struct {
	Email    string `validate:"required"`
	Password string `validate:"required"`
}

type loginResponse struct {
	Token string
}

// @Summary v1/Login
// @Description Validates user `email`, `password` and creates a Token with a default TTL of 10 days.
// @Tags Auth 📘
// @Accept  json
// @Produce json
// @Param Authorization header string true "Example: Bearer token"
// @Param Accept-Language header string true "Example: en-US" default(en-US)
// @Param data body loginRequest true "request JSON params"
// @Success 200 {object} loginResponse
// @Failure 400 {object} handlers.JSONMsgStr "Invalid request JSON"
// @Failure 403 {object} handlers.JSONMsgStr "Forbidden"
// @Failure 422 {object} handlers.JSONMsgStr "Params validation error"
// @Failure 500 {object} handlers.JSONMsgStr "Internal server error"
// @Router /auth/v1/Login [post]
func (s *Server) loginRoute(w http.ResponseWriter, r *http.Request) {
	// parse request JSON
	var req loginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Error("Could not decode JSON to login")
		handlers.JSONMsg(w, r, 400, "Invalid request JSON")
		return
	}

	token := ""

	handlers.JSONMsg(w, r, 200, loginResponse{Token: token})
}

// @Summary v1/Logout
// @Description Invalidates token present in the request Authorization header.
// @Tags Auth 📘
// @Accept  json
// @Produce json
// @Param Authorization header string true "Example: Bearer token"
// @Param Accept-Language header string true "Example: en-US" default(en-US)
// @Success 200 {object} interface{} "OK"
// @Failure 400 {object} handlers.JSONMsgStr "Invalid request JSON"
// @Failure 403 {object} handlers.JSONMsgStr "Forbidden"
// @Failure 422 {object} handlers.JSONMsgStr "Params validation error"
// @Failure 500 {object} handlers.JSONMsgStr "Internal server error"
// @Router /auth/v1/Logout [post]
func (s *Server) logoutRoute(w http.ResponseWriter, r *http.Request) {
	handlers.JSONMsg(w, r, 200, struct{}{})
}
