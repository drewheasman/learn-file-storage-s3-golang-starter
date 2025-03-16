package main

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/bootdotdev/learn-file-storage-s3-golang-starter/internal/auth"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerUploadThumbnail(w http.ResponseWriter, r *http.Request) {
	videoIDString := r.PathValue("videoID")
	videoID, err := uuid.Parse(videoIDString)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid ID", err)
		return
	}

	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't find JWT", err)
		return
	}

	userID, err := auth.ValidateJWT(token, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't validate JWT", err)
		return
	}

	fmt.Println("uploading thumbnail for video", videoID, "by user", userID)

	const maxMemory = 10 << 20

	err = r.ParseMultipartForm(maxMemory)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Couldn't get file data", err)
		return
	}

	fileData, fileHeader, err := r.FormFile("thumbnail")
	defer fileData.Close()
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Couldn't get file data", err)
		return
	}

	mediaType, _, err := mime.ParseMediaType(fileHeader.Header.Get("Content-Type"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Error parsing media type", err)
		return
	}
	if mediaType != "image/jpeg" && mediaType != "image/png" {
		respondWithError(w, http.StatusBadRequest, "Invalid media type", err)
		return
	}

	metadata, err := cfg.db.GetVideo(videoID)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Error looking up video", err)
		return
	}
	if metadata.UserID != userID {
		respondWithError(w, http.StatusUnauthorized, "User doesn't own video", err)
		return
	}

	newKey := make([]byte, 32)
	_, err = rand.Read(newKey)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error generating random videoId", err)
		return
	}
	newKeyString := base64.RawURLEncoding.EncodeToString(newKey)
	filename := newKeyString + "." + strings.Split(mediaType, "/")[1]
	filepath := filepath.Join(cfg.assetsRoot, "/", filename)
	fmt.Println("filepath", filepath)

	file, err := os.Create(filepath)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error creating file", err)
		return
	}

	writtenBytes, err := io.Copy(file, fileData)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error writing to file", err)
		return
	}
	fmt.Printf("%v bytes written to %v\n", writtenBytes, filename)

	newURL := "http://localhost:8091/assets/" + filename
	metadata.ThumbnailURL = &newURL

	err = cfg.db.UpdateVideo(metadata)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error updating video record", err)
		return
	}

	respondWithJSON(w, http.StatusOK, metadata)
}
