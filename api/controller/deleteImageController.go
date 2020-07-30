package controller

import (
	"fmt"
	util "imageStoreService/api/core/utilities"
	"net/http"
	"strings"
)

// swagger:route DELETE /deleteImage deleteImage removeImage
// Deletes an Image
// If image and album combination not present in DB, 404 is returned
// responses:
// 200: Success
// 400: BadRequest
// 404: NotFound

//DeleteImage servers the API /deleteImage to delete a new Album
func DeleteImage(w http.ResponseWriter, r *http.Request) {
	ImageCaption := r.URL.Query().Get("imageCaption")
	AlbumName := r.URL.Query().Get("albumName")
	ImageCaption = strings.TrimSpace(ImageCaption)
	AlbumName = strings.TrimSpace(AlbumName)
	var err error
	if AlbumName == "" {
		err = util.DeleteImageInDb(ImageCaption, "")
	} else {

		err = util.DeleteImageInDb(ImageCaption, AlbumName)

	}
	if err != nil {
		http.Error(w, http.StatusText(404), 404)
		fmt.Fprintf(w, "Error while deleting image - "+err.Error())
		return
	}
	s := "{Message : Image \"" + string(ImageCaption) + "\" deleted successfully." + "}"
	util.PostDataToKafka(s)
	fmt.Fprintf(w, "Image with caption : "+ImageCaption+" deleted successfully.")
	return
}
