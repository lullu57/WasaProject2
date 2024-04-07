package api

import (
	"encoding/json"
	"net/http"

	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/api/reqcontext"
	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/database"
	"github.com/julienschmidt/httprouter"
)

// HandleAddUser processes the request to add a new user
func HandleAddUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	var user database.User
	db := ctx.Database

	// Decode the incoming JSON payload into the user struct
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Add the user to the database
	r.Body.Close()
	ctx.Logger.Info("Adding user to the database")
	err := db.AddUser(&user)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError) // Sets the status code only
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "text/plain")
	ctx.Logger.Info("User created successfully")
}

func HandleSetUsername(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	currentUsername := ps.ByName("username") // Assuming username is the URL parameter

	var user database.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	r.Body.Close()

	ctx.Logger.Info("Setting new username")
	err := ctx.Database.SetUsername(currentUsername, user.Username)
	if err != nil {
		ctx.Logger.Error("Failed to update username: ", err)
		http.Error(w, "Failed to update username", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Username updated successfully"})
}

func HandleGetUserProfile(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	username := ps.ByName("username") // Assuming username is the URL parameter

	ctx.Logger.Info("Retrieving user profile for username: ", username)
	user, err := ctx.Database.GetUserProfile(username)
	if err != nil {
		ctx.Logger.Error("User not found: ", err)
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func HandleFollowUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	followedUsername := ps.ByName("username") // Extracting username from URL parameter

	followerID := "extractedFollowerID" // This should be obtained from your auth system, placeholder here
	followedID, err := ctx.Database.GetUserIDByUsername(followedUsername)
	if err != nil {
		http.Error(w, "Failed to find user: "+followedUsername, http.StatusNotFound)
		return
	}

	err = ctx.Database.FollowUser(followerID, followedID)
	if err != nil {
		ctx.Logger.Errorf("Error following user: %v", err)
		http.Error(w, "Failed to follow user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "User Followed successfully"})
}

func HandleUnfollowUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	followedUsername := ps.ByName("username") // Extracting username from URL parameter

	followerID := "extractedFollowerID" // This should be obtained from your auth system, placeholder here
	followedID, err := ctx.Database.GetUserIDByUsername(followedUsername)
	if err != nil {
		http.Error(w, "Failed to find user: "+followedUsername, http.StatusNotFound)
		return
	}

	err = ctx.Database.UnfollowUser(followerID, followedID)
	if err != nil {
		ctx.Logger.Errorf("Error unfollowing user: %v", err)
		http.Error(w, "Failed to unfollow user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "User Unfollowed successfully"})
}
