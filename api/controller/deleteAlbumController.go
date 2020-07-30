package controller

import (
	"fmt"
	util "imageStoreService/api/core/utilities"
	"net/http"

	"github.com/gorilla/mux"
)

// swagger:route DELETE /deleteAlbum deleteAlbum removeAlbum
// Deletes an Album
// If album not present in DB, 404 is returned
// responses:
// 200: Success
// 400: BadRequest
// 404: NotFound

//DeleteAlbum servers the API /deleteAlbum to delete a new Album
func DeleteAlbum(w http.ResponseWriter, r *http.Request) {
	// swagger:operation DELETE /deleteAlbum/{albumName}
	// ---
	// summary: Deletes an Album
	// description: If album not present in DB, 404 is returned
	// parameters:
	// - name: albumName
	//   in: path
	//   description: Name of the Album
	//   type: string
	//   required: true
	// responses:
	//   "200":
	//	 description: Successfully deleted
	//   "404":
	//	   description: Not Found
	vars := mux.Vars(r)
	albumName := vars["albumName"]
	err := util.DeleteAlbumInDb(albumName)
	if err != nil {
		http.Error(w, http.StatusText(404), 404)
		fmt.Fprintf(w, "Error while deleting album - "+err.Error())
		return
	}
	s := "{Message : Album \"" + string(albumName) + "\" deleted successfully." + "}"
	util.PostDataToKafka(s)
	fmt.Fprintf(w, "Album Name : "+albumName+" deleted successfully.")
	return
}
