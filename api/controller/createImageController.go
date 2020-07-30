package controller

import (
	"encoding/base64"
	"fmt"
	util "imageStoreService/api/core/utilities"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

// swagger:route POST /createImage createImage postImage
// Creates a new Image
// If Image is already present in DB, Error Conflict (409) will be returned.
// responses:
// 200: Success
// 400: BadRequest
// 409 : DuplicateAlbum
// 404: NotFound

//CreateImage servers the API /createImage to create a new Image for an already existing album
func CreateImage(w http.ResponseWriter, r *http.Request) {
	// swagger:operation POST /createImage?albumName=&imageCaption=
	// ---
	// summary: Creates a new Image
	// description: If Image is already present in DB for the Album, Error Conflict (409) will be returned.
	// parameters:
	// - name: albumName
	//   in: path
	//   description: Name of the Album
	//   type: string
	//   required: true
	// - name: albumName
	//   in: path
	//   description: Name of the Image
	//   type: string
	//   required: true
	// responses:
	//   "200":
	//	 description: Successfully created
	//   "409":
	//	   description: Duplicate Image for the Album
	//   "404":
	//	   description: Album is not found
	ImageCaption := r.URL.Query().Get("imageCaption")
	ImageCaption = strings.TrimSpace(ImageCaption)
	AlbumName := r.URL.Query().Get("albumName")
	AlbumName = strings.TrimSpace(AlbumName)
	if ImageCaption == "" || AlbumName == "" {
		http.Error(w, http.StatusText(404), 404)
		fmt.Fprintf(w, "Mandatory params 'imageCaption' or 'albumName' are missing")
		return
	}
	imageFile := make([]byte, 0)

	Image, _, err := r.FormFile("image")
	if err != nil {
		// No image
		log.Print("Error parsing image", err)
		http.Error(w, http.StatusText(400), 400)
		fmt.Fprintf(w, "Image not attached to Form")
		return
	}
	//if !strings.HasPrefix(fileHeader.Header["Content-Type"][0], "image") {
	// Something wrong with file type
	//}
	imageFile, err = ioutil.ReadAll(Image)
	//		var buff bytes.Buffer
	//		png.Encode(&buff, Image)
	if err != nil {
		// Error reading uploaded image from stream
		log.Print("Error reading image")
		http.Error(w, http.StatusText(400), 400)
		fmt.Fprintf(w, "Error Reading image")
		return
	}
	encodedImage := base64.StdEncoding.EncodeToString(imageFile)
	err = util.UploadImage(encodedImage, AlbumName, ImageCaption)
	if err != nil && strings.Contains(err.Error(), "no rows in result set") {
		fmt.Println("Error while Inserting the image", err)
		http.Error(w, http.StatusText(404), 404)
		fmt.Fprintf(w, "Album is not found in DB. Use POST to create!")
		return
	} else if err != nil && strings.Contains(err.Error(), "Duplicate entry") {
		fmt.Println("Error while Inserting the image", err)
		http.Error(w, http.StatusText(409), 409)
		fmt.Fprintf(w, "Duplicate Entry")
		return
	}
	fmt.Println("Done")

	fmt.Println("Debug  1")
	s := "{Message : Image \"" + string(ImageCaption) + "\" created successfully." + "}"
	util.PostDataToKafka(s)

	fmt.Fprintf(w, "%+v", string(ImageCaption)+" created succesfully.")
	return
}
