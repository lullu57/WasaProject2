package api

import (
	"io/ioutil"
	"net/http"
	"time"

	"encoding/json"

	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/api/reqcontext"
	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/database"
	"github.com/gofrs/uuid"
	"github.com/julienschmidt/httprouter"
)

func handleUploadPhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	// Extract username from context
	if ctx.User == nil {
		w.WriteHeader(http.StatusUnauthorized) // Sets the status code only
		return
	}
	userId := ctx.User.ID
	ctx.Logger.Info("Called successfully")
	// Read image data from the request body
	// Parse the multipart form
	err := r.ParseMultipartForm(10 << 20) // For example, max 10 MB file size
	if err != nil {
		w.WriteHeader(http.StatusBadRequest) // Sets the status code only
		return
	}

	// Retrieve the file from form data
	file, _, err := r.FormFile("image") // "image" should be the name of your file input field
	if err != nil {
		w.WriteHeader(http.StatusBadRequest) // Sets the status code only
		return
	}
	defer file.Close()

	// Read the file data
	ImageData, err := ioutil.ReadAll(file)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest) // Sets the status code only
		return
	}
	defer r.Body.Close()

	ctx.Logger.Info("Received image data length: ", len(ImageData))
	// Set current time as Timestamp
	Timestamp := time.Now()

	// Create a Photo struct
	photo := database.Photo{
		ID:        uuid.Must(uuid.NewV4()).String(),
		UserID:    userId,
		ImageData: ImageData,
		Timestamp: Timestamp,
		Likes:     []database.Like{},
		Comments:  []database.Comment{},
	}
	ctx.Logger.Info("Photo created " + photo.Timestamp.String())
	// Call AddPhoto method to insert the photo into the database
	err = ctx.Database.AddPhoto(photo)
	if err != nil {
		ctx.Logger.WithError(err).Error("Failed to add photo to the database")
		ctx.Logger.Info("Failed to add photo to the database")
		w.WriteHeader(http.StatusInternalServerError) // Sets the status code only
		return
	}
	ctx.Logger.Info("Photo added to the database")
	// Respond with success message
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusCreated)
	_, err = w.Write([]byte("Photo uploaded successfully"))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError) // Sets the status code only
		return
	}
}

func handleGetPhotos(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	// Retrieve all photos from the database
	photos, err := ctx.Database.GetPhotos()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError) // Sets the status code only
		return
	}
	// Respond with the list of photos
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(photos)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError) // Sets the status code only
		return
	}
}

func handleGetMyStream(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	if ctx.User == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	photos, err := ctx.Database.GetMyStream(ctx.User.ID)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	ctx.Logger.Info("My stream fetched")
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(photos)
}

func handleDeletePhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	photoID := ps.ByName("photoId")
	if photoID == "" {
		http.Error(w, "Invalid photo ID", http.StatusBadRequest)
		return
	}

	err := ctx.Database.DeletePhoto(photoID)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	ctx.Logger.Infof("Photo %s deleted by %s", photoID, ctx.User.Username)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Photo deleted successfully"))
}
