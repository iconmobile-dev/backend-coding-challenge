package user

import (
	"encoding/json"
	"net/http"

	"github.com/iconmobile-dev/go-coding-challenge/lib/handlers"
	"github.com/iconmobile-dev/go-coding-challenge/lib/userlib"
)

type userResponse struct {
	User userlib.User
}

type userCreateRequest struct {
	Email     string `validate:"required"`
	Password  string `validate:"required" minLength:"8"`
	FirstName string
	LastName  string
}

// @Summary v1/UserCreate
// @Description Creates an User
// @Tags User 📘
// @Accept  json
// @Produce json
// @Param data body userCreateRequest true "request JSON params"
// @Success 200 {object} userResponse
// @Failure 400 {object} handlers.JSONMsgStr "Invalid request JSON"
// @Failure 403 {object} handlers.JSONMsgStr "Forbidden"
// @Failure 422 {object} handlers.JSONMsgStr "Params validation error"
// @Failure 500 {object} handlers.JSONMsgStr "Internal server error"
// @Router /users/v1/UserCreate [post]
func (s *Server) userCreateRoute(w http.ResponseWriter, r *http.Request) {
	var req userCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Error("Could not decode JSON to create User")
		handlers.JSONMsg(w, r, 400, "Invalid request JSON")
		return
	}

	// create User
	user := userlib.User{
		Email:     req.Email,
		Password:  req.Password,
		FirstName: req.FirstName,
		LastName:  req.LastName,
	}

	err := user.Insert(s.db, s.cache)
	if err != nil {
		log.Error("user.Insert err:", err)
		handlers.JSONMsgErr(w, r, err, "Could not create User")
		return
	}

	// remove sensitive data
	user = removeSensitiveDataFromUser(user)

	handlers.JSONMsg(w, r, 200, userResponse{
		User: user,
	})
}

type userGetRequest struct {
	ID int `validate:"required"`
}

// @Summary v1/UserGet
// @Description Gets an User
// @Tags User 📘
// @Accept  json
// @Produce json
// @Param Authorization header string true "Example: Bearer token"
// @Param Accept-Language header string true "Example: en-US" default(en-US)
// @Param data body userGetRequest true "request JSON params"
// @Success 200 {object} userResponse
// @Failure 400 {object} handlers.JSONMsgStr "Invalid request JSON"
// @Failure 403 {object} handlers.JSONMsgStr "Forbidden"
// @Failure 422 {object} handlers.JSONMsgStr "Params validation error"
// @Failure 500 {object} handlers.JSONMsgStr "Internal server error"
// @Router /users/v1/UserGet [post]
func (s *Server) userGetRoute(w http.ResponseWriter, r *http.Request) {
	// parse request JSON
	var req userGetRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Error("Could not decode JSON to get User")
		handlers.JSONMsg(w, r, 400, "Invalid request JSON")
		return
	}

	// load the User
	user, err := userlib.UserByID(req.ID, s.db)
	if err != nil {
		log.Error("userlib.UserByID err:", err)
		handlers.JSONMsgErr(w, r, err, "Could not update User")
		return
	}

	// remove sensitive data
	user = removeSensitiveDataFromUser(user)

	log.Verbose("Got User with ID", user.ID)
	handlers.JSONMsg(w, r, 200, userResponse{
		User: user,
	})
}

// removeSensitiveDataFromUser removes the Password from the userlib.User
// removes the Email from the userlib.User, if the role of the authenticated User is less than Support and it is not the same User
func removeSensitiveDataFromUser(user userlib.User) userlib.User {
	user.Password = ""
	user.Email = ""

	return user
}

type userDeleteRequest struct {
	ID int `validate:"required"`
}

// @Summary v1/UserDelete
// @Description Deletes an User
// @Tags User 📘
// @Accept  json
// @Produce json
// @Param Authorization header string true "Example: Bearer token"
// @Param Accept-Language header string true "Example: en-US" default(en-US)
// @Param data body userDeleteRequest true "request JSON params"
// @Success 200 {object} interface{} "OK"
// @Failure 400 {object} handlers.JSONMsgStr "Invalid request JSON"
// @Failure 403 {object} handlers.JSONMsgStr "Forbidden"
// @Failure 422 {object} handlers.JSONMsgStr "Params validation error"
// @Failure 500 {object} handlers.JSONMsgStr "Internal server error"
// @Router /users/v1/UserDelete [post]
func (s *Server) userDeleteRoute(w http.ResponseWriter, r *http.Request) {
	// parse request JSON
	var req userDeleteRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Error("Could not decode JSON to delete Users")
		handlers.JSONMsg(w, r, 400, "Invalid request JSON")
		return
	}

	// load the User
	user, err := userlib.UserByID(req.ID, s.db)
	if err != nil {
		log.Error("userlib.UserById err:", err)
		handlers.JSONMsgErr(w, r, err, "Could not delete User")
		return
	}

	// delete the User
	err = user.Delete(s.db)
	if err != nil {
		log.Error("userlib.Delete err:", err)
		handlers.JSONMsgErr(w, r, err, "Could not delete User")
		return
	}

	handlers.JSONMsg(w, r, 200, map[string]string{})
}
