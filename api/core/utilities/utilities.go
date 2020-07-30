package utilities

import (
	"archive/zip"
	"database/sql"
	"errors"
	"fmt"
	"imageStoreService/api/core/models"
	"io"
	"log"
	"os"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

const (
	DB_HOST = "tcp(test-mysql:3306)"
	DB_NAME = "test"
	DB_USER = /*"root"*/ "root"
	DB_PASS = /*""*/ "my-secret-pw"
)

var db *sql.DB
var dsn = DB_USER + ":" + DB_PASS + "@" + DB_HOST + "/"

func init() {
	//    dsn := DB_USER + ":" + DB_PASS + "@" + DB_HOST + "/"
	var err error
	db, err = sql.Open("mysql", dsn)
	// if there is an error opening the connection, handle it
	if err != nil {
		log.Print(err.Error())
	} else {
		log.Print("DB connected successfully")
	}
	res, err := db.Exec(`CREATE DATABASE ImageStore`)
	if err != nil {
		log.Print("Printing Error : ", err)
	}
	log.Print("Result", res)
	res, err = db.Exec(`USE ImageStore;`)
	res, err = db.Exec(`CREATE TABLE Albums( AlbumId int NOT NULL AUTO_INCREMENT, AlbumName varchar(255) NOT NULL UNIQUE,PRIMARY KEY (AlbumId));`)
	if err != nil {
		log.Print("Printing Error : ", err)
	}
	res, err = db.Exec(`CREATE TABLE Images( ImageId int NOT NULL AUTO_INCREMENT, ImageCaption varchar(255), Image LONGBLOB,AlbumId int, PRIMARY KEY (ImageId), CONSTRAINT FK_AlbumImage FOREIGN KEY (AlbumId) REFERENCES Albums(AlbumId));`)
	if err != nil {
		log.Print("Printing Error : ", err)
	}

	defer db.Close()

}
func session() {
	var err error
	db, err = sql.Open("mysql", dsn)
	_, err = db.Exec(`USE ImageStore`)
	if err != nil {
		log.Print("Printing Error : ", err)
	}

}

func CreateAlbumInDb(AlbumName string) error {
	session()
	res, err := db.Exec("Insert into Albums (AlbumName) VALUES( " + "'" + AlbumName + "'" + ");")
	if err != nil {
		log.Print("Error while inserting", err)
	}
	fmt.Println(res)
	defer db.Close()
	return err

}
func DeleteAlbumInDb(AlbumName string) (err error) {
	session()
	var id int
	err = db.QueryRow("SELECT AlbumId FROM Albums WHERE AlbumName=?", AlbumName).Scan(&id)
	if err != nil {
		if strings.Contains(err.Error(), "no rows in result set") {
			err = errors.New("Album not present in DB.")
			return

		}
		fmt.Println("Error in select ", err)
		return
	}
	statement := "Delete from Images where AlbumId=?"
	_, err = db.Exec(statement, id)
	if err != nil {
		log.Print("Error while Deleting images", err)
		return
	}
	statement = "Delete from Albums where AlbumId=?"
	_, err = db.Exec(statement, id)
	if err != nil {
		log.Print("Error while Deleting Albums", err)
		return
	}

	defer db.Close()
	return err

}
func DeleteImageInDb(caption string, AlbumName string) (err error) {
	session()
	var id int
	err = db.QueryRow("SELECT ImageId FROM Images WHERE ImageCaption=?", caption).Scan(&id)
	if err != nil {
		if strings.Contains(err.Error(), "no rows in result set") {
			err = errors.New("Image not present in DB.")
			return

		}
		fmt.Println("Error in select ", err)
		return
	}
	if AlbumName != "" {
		err = db.QueryRow("SELECT AlbumId FROM Albums WHERE AlbumName=?", AlbumName).Scan(&id)
		if err != nil {
			if strings.Contains(err.Error(), "no rows in result set") {
				err = errors.New("AlbumName " + AlbumName + "is not present in DB.")
				return

			}
			fmt.Println("Error in select ", err)
			return
		}
		var imageId int
		err = db.QueryRow("SELECT ImageId FROM Images WHERE AlbumId=? and ImageCaption=?", id, caption).Scan(&imageId)
		if err != nil {
			if strings.Contains(err.Error(), "no rows in result set") {
				err = errors.New("Image " + caption + " and AlbumName " + AlbumName + " combination is not present in DB.")
				return

			}
			fmt.Println("Error in select ", err)
			return
		}

		statement := "Delete from Images where ImageCaption=? and AlbumId=?"
		_, err = db.Exec(statement, caption, id)
		if err != nil {
			if strings.Contains(err.Error(), "no rows in result set") {
				fmt.Println("Not deleted")
				err = errors.New("Image not present in DB.")
				return

			}
			log.Print("Error while Deleting Images", err)
			return
		}
	} else {

		statement := "Delete from Images where ImageCaption=?"
		_, err = db.Exec(statement, caption)
		if err != nil {
			if strings.Contains(err.Error(), "no rows in result set") {
				err = errors.New("Image not present in DB.")
				return

			}
			log.Print("Error while Deleting Images", err)
			return
		}

	}
	fmt.Println("All success")
	defer db.Close()
	return err

}

func UploadImage(encodeImage string, AlbumName string, ImageCaption string) (err error) {
	session()
	//query := "SELECT AlbumId FROM Albums WHERE AlbumName=?"
	var id int
	err = db.QueryRow("SELECT AlbumId FROM Albums WHERE AlbumName=?", AlbumName).Scan(&id)
	//	fmt.Println("Id ", id)
	//defer rows.Close()
	/*var id int
	for rows.Next() {
		//	var id int
		err = rows.Scan(&id)
		if err != nil {
			fmt.Println("Error while scanning rows ", err)
		}
		fmt.Println("AlbumId ", id)
	}*/
	if err != nil {
		fmt.Println("Error in select ", err)
		return
	}
	//	fmt.Println(encodeImage)
	statement := "Insert into Images (ImageCaption,Image,AlbumId) VALUES(?,?,?)"
	_, err = db.Exec(statement, ImageCaption, encodeImage, id)
	if err != nil {
		log.Print("Error while inserting", err)
		return
	}
	//	fmt.Println(res)
	defer db.Close()
	return

}
func GetImageByAlbumAndCaption(AlbumName string, ImageCaption string) (err error, Image string) {
	session()
	var id int
	err = db.QueryRow("SELECT AlbumId FROM Albums WHERE AlbumName=?", AlbumName).Scan(&id)
	fmt.Println("Id ", id)
	if err != nil {
		fmt.Println("Error in select ", err)
		err = errors.New("Album not present in DB")
		return
	}
	fmt.Println("Image Caption", ImageCaption+"----")
	query := "Select Image FROM Images WHERE AlbumId=? and ImageCaption=?"
	err = db.QueryRow(query, id, ImageCaption).Scan(&Image)
	if err != nil {
		fmt.Println("Error in select 2", err)
		err = errors.New("Image not present in DB")
		return
	}
	return

}

func GetAllImageFromDb() (allImages []models.Images, err error) {
	session()
	rows, err := db.Query("Select ImageId,ImageCaption,Image,AlbumId from Images")
	if err != nil {
		fmt.Println("Error while fetching images")
		return

	}
	defer rows.Close()
	for rows.Next() {
		var data models.Images
		err = rows.Scan(&data.ImageId, &data.ImageCaption, &data.ImageData, &data.AlbumId)
		if err != nil {
			fmt.Println("Error while scanning rows ", err)
		}
		allImages = append(allImages, data)

	}
	return

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
