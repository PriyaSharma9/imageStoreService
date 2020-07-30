package router

import (
	"log"
	"net/http"

	"imageStoreService/api/controller"
	util "imageStoreService/api/core/utilities"

	"github.com/gorilla/mux"
)

//HandleRouter handles exposing all the API's for Image Store Service
func HandleRouter() {

	util.Writer = util.NewKafkaWriter(util.KafkaURL, util.Topic)
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/createAlbum/{albumName}", controller.CreateAlbum).Methods("POST")
	router.HandleFunc("/createImage", controller.CreateImage).Methods("POST")
	router.HandleFunc("/getImage", controller.GetImage).Methods("GET")
	router.HandleFunc("/deleteAlbum/{albumName}", controller.DeleteAlbum).Methods("DELETE")
	router.HandleFunc("/deleteImage", controller.DeleteImage).Methods("DELETE")
	router.HandleFunc("/getAllImages", controller.GetAllImages).Methods("GET")
	log.Fatal(http.ListenAndServe(":8193", router))
	defer util.Writer.Close()

}
