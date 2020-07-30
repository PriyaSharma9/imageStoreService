package controller

import (
	"fmt"
	util "imageStoreService/api/core/utilities"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

// swagger:route POST /createAlbum createAlbum postAlbum
// Creates a new Album
//  If album is already present in DB, Error Conflict (409) will be returned.
// responses:
//  200: Success
//  400: BadRequest
// 409 : DuplicateAlbum

//CreateAlbum servers the API /createAlbum to create a new Album
func CreateAlbum(w http.ResponseWriter, r *http.Request) {
	// swagger:operation POST /createAlbum/{albumName}
	// ---
	// summary: Creates a new Album
	// description: If album is already present in DB, Error Conflict (409) will be returned.
	// parameters:
	// - name: albumName
	//   in: path
	//   description: Name of the Album
	//   type: string
	//   required: true
	// responses:
	//   "200":
	//	 description: Successfully created
	//   "409":
	//	   description: Duplicate Album
	vars := mux.Vars(r)
	albumName := vars["albumName"]
	albumName = strings.TrimSpace(albumName)
	err := util.CreateAlbumInDb(albumName)
	if err != nil {
		http.Error(w, http.StatusText(409), 409)
		fmt.Fprintf(w, "Error while creating album - "+err.Error())
		return
	}
	s := "{Message : Album \"" + string(albumName) + "\" created successfully." + "}"
	util.PostDataToKafka(s)
	fmt.Fprintf(w, "Album Name : "+albumName+" created successfully.")
	return

}
