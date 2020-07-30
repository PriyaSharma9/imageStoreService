// Package main
//
// The purpose of this application is to provide an application
// that will expose REST API's to store,retrieve,delete and get Images and Albums.
// Terms Of Service:
//
// there are no TOS at this moment, use at your own risk we take no responsibility
//
//     Schemes: http
//     Host: localhost
//     BasePath: /
//     Version: 0.0.1
//     Contact: Priya Sharma<priya.sharma6693@gmail.com>
//
//     Consumes:
//     - application/image
//     - application/json
//
//     Produces:
//     - application/json
//     - application/zip
//     - application/image
// swagger:meta
package main

import (
	"imageStoreService/api/router"
)

func main() {

	router.HandleRouter()

}
