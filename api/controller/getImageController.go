package controller

import (
	"archive/zip"
	"bytes"
	"encoding/base64"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	util "imageStoreService/api/core/utilities"
	"io"
	"log"
	"net/http"
	"os"
	"reflect"
	"strconv"
	"strings"
)

// swagger:route GET /getAllImages getAllImages listImages
// Get all the images present in DB
// If no images are present - 404 is returned
// responses:
//  200: Success
//  400: BadRequest
// 404 : NotFound

//GetAllImages servers the API /getAllImages to get all the images present in DB
func GetAllImages(w http.ResponseWriter, r *http.Request) {

	allImages, err := util.GetAllImageFromDb()
	fmt.Println("length ", len(allImages))
	if err != nil {
		http.Error(w, http.StatusText(404), 404)
		fmt.Fprintf(w, "Could not get images - "+err.Error())
		return
	}
	flags := os.O_WRONLY | os.O_CREATE | os.O_TRUNC
	file, err := os.OpenFile("test.zip", flags, 0644)
	if err != nil {
		log.Fatalf("Failed to open zip for writing: %s", err)
	}
	defer file.Close()
	zipw := zip.NewWriter(file)
	defer zipw.Close()
	for i, data := range allImages {
		Image := data.ImageData
		caption := data.ImageCaption
		unbased, _ := base64.StdEncoding.DecodeString(Image)
		img, imageType, err := image.Decode(bytes.NewBuffer(unbased))
		if err != nil {
			fmt.Println("Error while decoding the image ", err)
		}
		fmt.Println("After ------")
		switch imageType {
		case "png":
			fmt.Println("Case PNG ")
			fileName := "./" + caption + "" + strconv.Itoa(i) + ".png"
			out, err := os.Create(fileName)
			if err != nil {
				fmt.Println("Error while creating image in path ", err)
			}
			fmt.Println("------------------ 2 before ----------- : ")
			err = png.Encode(out, img)
			fmt.Println("------------------ 2 After ---------------")
			if err != nil {
				fmt.Println("Error while png encoding ", err)
			}
			err = appendFiles(fileName, zipw)
			if err != nil {
				log.Fatalf("Failed to add file %s to zip: %s", fileName, err)
			}
		case "jpeg":
			fmt.Println("Case JPEG")
			fileName := "./" + caption + "" + strconv.Itoa(i) + ".jpeg"
			out, err := os.Create(fileName)
			if err != nil {
				fmt.Println("Error while creating image in path ", err)
			}
			fmt.Println("------------------ 2 before ----------- : ")
			err = jpeg.Encode(out, img, &jpeg.Options{Quality: 105})
			fmt.Println("------------------ 2 After ---------------")
			if err != nil {
				fmt.Println("Error while png encoding ", err)
			}
			err = appendFiles(fileName, zipw)
			if err != nil {
				log.Fatalf("Failed to add file %s to zip: %s", fileName, err)
			}
			// w.Header().Set("Content-Type", "image/jpeg") // <-- set the content-type header

		}

	}
	w.Header().Set("Content-Type", "application/zip")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s.zip\"", "test.zip"))
	//io.Copy(w, buf)
	//  w.Write(buf.Bytes())
	io.Copy(w, file)
	fmt.Fprintf(w, "getAllImages ")
	return

}

// swagger:route GET /getImage getImage listImage
// Get an Image for Album and Image Combination
//  If image and album combination not present in DB, 404 is returned
// responses:
//  200: Success
//  400: BadRequest
// 404 : NotFound

//GetImage servers the API /getImage to get an image present in DB with the input for query params imageCaption and albumName
func GetImage(w http.ResponseWriter, r *http.Request) {
	ImageCaption := r.URL.Query().Get("imageCaption")
	AlbumName := r.URL.Query().Get("albumName")
	if ImageCaption == "" || AlbumName == "" {
		http.Error(w, http.StatusText(404), 404)
		fmt.Fprintf(w, "Mandatory params 'imageCaption' or 'albumName' are missing")
		return
	}
	ImageCaption = strings.TrimSpace(ImageCaption)
	fmt.Println("ImageCaption Debug1", ImageCaption)
	err, Image := util.GetImageByAlbumAndCaption(AlbumName, ImageCaption)
	if err != nil && strings.Contains(err.Error(), "Album not present in DB") {
		fmt.Println("Error while getting image ", err)
		http.Error(w, http.StatusText(404), 404)
		fmt.Fprintf(w, "Album is not found in DB. Use POST to create!")
		return
	} else if err != nil && strings.Contains(err.Error(), "Image not present in DB") {
		fmt.Println("Error while getting image ", err)
		http.Error(w, http.StatusText(404), 404)
		fmt.Fprintf(w, "Image is not found in DB. Use POST to create!")
		return
	}
	fmt.Println(reflect.TypeOf(Image))
	//	b := []byte(Image)
	fmt.Println("Before ------")
	unbased, _ := base64.StdEncoding.DecodeString(Image)
	img, imageType, err := image.Decode(bytes.NewBuffer(unbased))
	if err != nil {
		fmt.Println("Error while decoding the image ", err)
	}
	fmt.Println("After ------")
	switch imageType {
	case "png":
		fmt.Println("Case PNG ")
		out, err := os.Create("./QRImg.png")
		if err != nil {
			fmt.Println("Error while creating image in path ", err)
		}
		fmt.Println("------------------ 2 before ----------- : ")
		err = png.Encode(out, img)
		fmt.Println("------------------ 2 After ---------------")
		if err != nil {
			fmt.Println("Error while png encoding ", err)
		}
	case "jpeg":
		fmt.Println("Case JPEG")
		out, err := os.Create("./QRImg.jpeg")
		if err != nil {
			fmt.Println("Error while creating image in path ", err)
		}
		fmt.Println("------------------ 2 before ----------- : ")
		err = jpeg.Encode(out, img, &jpeg.Options{Quality: 105})
		fmt.Println("------------------ 2 After ---------------")
		if err != nil {
			fmt.Println("Error while png encoding ", err)
		}
		w.Header().Set("Content-Type", "image/jpeg") // <-- set the content-type header
		io.Copy(w, out)

	}
	/*	out, err := os.Create("./QRImg.png")

			if err != nil {
		             fmt.Println("Error while creating image in path ",err)
			}
			fmt.Println("------------------ 2 before ----------- : ", img)
			err = png.Encode(out, img)
			fmt.Println("------------------ 2 After ---------------")
			if err != nil {
		           fmt.Println("Error while png encoding ",err)
			}*/
	fmt.Fprintf(w, "%+v", "Found Image ")

}
func appendFiles(filename string, zipw *zip.Writer) error {
	file, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("Failed to open %s: %s", filename, err)
	}
	defer file.Close()

	wr, err := zipw.Create(filename)
	if err != nil {
		msg := "Failed to create entry for %s in zip file: %s"
		return fmt.Errorf(msg, filename, err)
	}

	if _, err := io.Copy(wr, file); err != nil {
		return fmt.Errorf("Failed to write %s to zip: %s", filename, err)
	}

	return nil
}
