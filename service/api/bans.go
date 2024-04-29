package api

import (
	"encoding/json"
	"net/http"

	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/api/reqcontext"
	"github.com/julienschmidt/httprouter"
)

func handleBanUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	username := ps.ByName("username")
	if username == "" {
		http.Error(w, "Invalid username", http.StatusBadRequest)
		return
	}

	// Assuming you have a method to fetch user ID by username
	bannedUser, err := ctx.Database.GetUserByUsername(username)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	bannedBy := ctx.User.ID
	err = ctx.Database.BanUser(bannedBy, bannedUser.ID)
	if err != nil {
		if err.Error() == "user is already banned" {
			http.Error(w, "User is already banned", http.StatusConflict)
			return
		}
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	ctx.Logger.Infof("User %s banned by %s", bannedUser.Username, ctx.User.Username)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("User successfully banned"))
}

// Handler for unbanning a user
func handleUnbanUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	bannedUsername := ps.ByName("username")

	if bannedUsername == "" {
		ctx.Logger.Infof("Invalid parameters")
		http.Error(w, "Invalid parameters", http.StatusBadRequest)
		return
	}

	// Fetch user IDs by username
	bannedUser, err := ctx.Database.GetUserByUsername(bannedUsername)
	if err != nil {
		http.Error(w, "Banned user not found", http.StatusNotFound)
		return
	}
	bannerUser := ctx.User.ID

	err = ctx.Database.UnbanUser(bannerUser, bannedUser.ID)
	if err != nil {
		ctx.Logger.Infof("Internal server error " + err.Error())
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	ctx.Logger.Infof("User %s unbanned by %s", bannedUser.Username, ctx.User.Username)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("User successfully unbanned"))
}

// Handler for getting all banned users
func handleGetBannedUsers(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	bannedUsers, err := ctx.Database.GetBans()
	if err != nil {
		ctx.Logger.Infof("Internal server error1 " + err.Error())
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(bannedUsers)
	ctx.Logger.Infof("Banned users fetched")
	if err != nil {
		ctx.Logger.Infof("Internal server error " + err.Error())
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}
