package db

import (
	"NewPhoto/utils"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"math"
	"os"

	. "NewPhoto/config"

	_ "github.com/go-sql-driver/mysql"
)

var (
	DBInstanse      = New()
	DefaultCapacity = 15 * math.Pow(10, 9)
)

type DB struct {
	db *sql.DB
}

func (d *DB) CreateDB() {

	password, ok := os.LookupEnv("mysqlPassword")
	if !ok{
		log.Fatalln("mysqlPassword is not written in credentials.sh file")
	}
	username, ok := os.LookupEnv("mysqlUsername")
	if !ok{
		log.Fatalln("mysqlUsername is not written in credentials.sh file")
	}

	table, ok := os.LookupEnv("mysqlTable")
	if !ok{
		log.Fatalln("mysqlTable is not written in credentials.sh file")
	}

	addr, ok := os.LookupEnv("mysqlAddr")
	if !ok{
		log.Fatalln("mysqlAddr is not written in credentials.sh file")
	}

	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s)/%s", username, password, addr, table))
	if err != nil {
		Logger.WriteFatal(err.Error())
	}
	d.db = db
}

func (d *DB) CreateTables() {

	_, err := d.db.Exec(fmt.Sprintf("CREATE TABLE IF NOT EXISTS users (userid VARCHAR(200), firstname VARCHAR(200), secondname VARCHAR(200), storage DOUBLE DEFAULT %f)", DefaultCapacity))
	if err != nil {
		Logger.WriteFatal(err.Error())
	}
	_, err = d.db.Exec("CREATE TABLE IF NOT EXISTS photos (userid VARCHAR(200), photo MEDIUMBLOB, thumbnail MEDIUMBLOB, extension VARCHAR(10), size DOUBLE, album VARCHAR(200), tags VARCHAR(200));")
	if err != nil {
		Logger.WriteFatal(err.Error())
	}
	_, err = d.db.Exec("CREATE TABLE IF NOT EXISTS albums (userid VARCHAR(200), name VARCHAR(200));")
	if err != nil {
		Logger.WriteFatal(err.Error())
	}
}

func (d *DB) LoginUser(login, pass string) (string, error) {
	passedEncodedCredentials := utils.EncodeLogin(login, pass)
	rows, err := d.db.Query("SELECT userid FROM users")
	if err != nil {
		Logger.WriteFatal(err.Error())
	}
	defer rows.Close()

	for rows.Next() {
		var userid string
		rows.Scan(&userid)
		if userid == string(passedEncodedCredentials) {
			return string(passedEncodedCredentials), nil
		}
	}
	return "", errors.New("such user does not exist")
}

func (d *DB) RegisterUser(login, pass, firstname, secondname string) {
	_, err := d.db.Exec("INSERT IGNORE INTO users (userid, firstname, secondname) VALUES (?, ?, ?)", utils.EncodeLogin(login, pass), firstname, secondname)
	if err != nil {
		Logger.WriteFatal(err.Error())
	}
}

func (d *DB) IsLogin(uid string) bool {
	rows, err := d.db.Query("SELECT userid FROM users")
	if err != nil {
		Logger.WriteFatal(err.Error())
	}
	defer rows.Close()

	for rows.Next() {
		var userid string
		rows.Scan(&userid)
		if userid == uid {
			return true
		}
	}
	return false
}

func (d *DB) UploadEqualPhoto(userid string, photo, thumbnail []byte, extension string, size float64, tags string) {
	// Inserts passed photo and its options ...

	_, err := d.db.Exec("INSERT INTO photos (userid, photo, thumbnail, extension, size, tags) VALUES (?, ?, ?, ?, ?, ?)",
		userid,
		photo,
		thumbnail,
		extension,
		size,
		tags,
	)
	if err != nil {
		Logger.WriteFatal(err.Error())
	}
}

func (d *DB) GetFullPhotoByThumbnail(userid string, thumbnail []byte) []byte {
	row := d.db.QueryRow("SELECT photo FROM photos WHERE userid = ? AND thumbnail = ?", userid, thumbnail)
	var photo []byte
	row.Scan(&photo)
	return photo
}

func (d *DB) GetUserinfo(userid string) (string, string, float64) {
	// Returns all the available storage for uploading by the passed user ...

	row := d.db.QueryRow("SELECT storage FROM users WHERE userid = ?", userid)
	var storage float64
	err := row.Scan(&storage)
	if err != nil {
		Logger.WriteFatal(err.Error())
	}

	rows, err := d.db.Query("SELECT size FROM photos WHERE userid = ?", userid)

	if err != nil {
		Logger.WriteFatal(err.Error())
	}
	defer rows.Close()

	var allphotossize float64
	for rows.Next() {
		var storage float64
		rows.Scan(&storage)
		allphotossize += storage
	}

	row = d.db.QueryRow("SELECT firstname, secondname FROM users WHERE userid = ?", userid)

	var firstname string
	var secondname string

	row.Scan(&firstname, &secondname)
	//fmt.Println(storage, allphotossize, "It is a size of all photos", storage - allphotossize)
	return firstname, secondname, storage - allphotossize
}

func (d *DB) AllPhotos(userid string) [][5]interface{} {
	// Returns all the photos uploaded by user ...

	rows, err := d.db.Query("SELECT * FROM photos WHERE userid = ?", userid)
	if err != nil {
		Logger.WriteFatal(err.Error())
	}
	defer rows.Close()
	model := struct {
		Userid    string
		Photo     []byte
		Thumbnail []byte
		Extension string
		Size      float64
		Album     sql.NullString
		Tags      sql.NullString
	}{}
	var result [][5]interface{}
	for rows.Next() {
		err := rows.Scan(&model.Userid, &model.Photo, &model.Thumbnail, &model.Extension, &model.Size, &model.Album, &model.Tags)
		if err != nil{
			Logger.WriteFatal(err.Error())
		}
		result = append(result, [5]interface{}{model.Photo, model.Thumbnail, model.Extension, model.Size, model.Tags.String})
	}
	return result
}

func (d *DB) CreateAlbum(userid, name string) bool {
	_, err := d.db.Exec("INSERT IGNORE INTO albums (userid, name) VALUES(?, ?)", userid, name)
	if err != nil {
		Logger.WriteFatal(err.Error())
	}
	return true
}

func (d *DB) AlbumPhotos(userid, name string) [][5]interface{} {
	rows, err := d.db.Query("SELECT userid, photo, thumbnail, extension, size, album FROM photos WHERE userid = ? AND album = ?", userid, name)
	if err != nil {
		Logger.WriteFatal(err.Error())
	}
	defer rows.Close()
	var model struct {
		Userid    string
		Photo     []byte
		Thumbnail []byte
		Extension string
		Size      float64
		Album     sql.NullString
	}
	var result [][5]interface{}
	for rows.Next() {
		err := rows.Scan(&model.Userid, &model.Photo, &model.Thumbnail, &model.Extension, &model.Size, &model.Album)
		if err != nil {
			continue
		}
		result = append(result, [5]interface{}{model.Photo, model.Thumbnail, model.Extension, model.Size, model.Album.String})
	}
	return result
}

func (d *DB) GetAllAlbums(userid string) [][3]interface{} {
	//Returns all the albums of user with passed userid and the last added photo to
	//each album

	
	rows, err := d.db.Query("SELECT name FROM albums WHERE userid = ?", userid)
	if err != nil {
		Logger.WriteFatal(err.Error())
	}
	
	defer rows.Close()

	var albums []string

	for rows.Next() {
		var album string
		err := rows.Scan(&album)
		if err != nil {
			Logger.WriteFatal(err.Error())
			continue
		}
		albums = append(albums, album)
	}
	

	var result [][3]interface{}
	for _, album := range albums {
		row := d.db.QueryRow("SELECT photo, thumbnail FROM photos WHERE userid = ? AND album = ? ORDER BY userid DESC LIMIT 1", userid, album)
		var photo []byte
		var thumbnail []byte
		row.Scan(&photo, &thumbnail)
		result = append(result, [3]interface{}{album, photo, thumbnail})
	}
	return result
}

func (d *DB) DeleteAlbum(userid, name string) bool {
	_, err := d.db.Exec("DELETE FROM albums WHERE userid = ? AND name = ?", userid, name)
	if err != nil {
		Logger.WriteFatal(err.Error())
	}
	return true
}

func (d *DB) UploadEqualPhotoToAlbum(userid, extension, album string, size float64, photo, thumbnail []byte) {
	_, err := d.db.Exec("INSERT INTO photos (userid, photo, thumbnail, extension, size, album) VALUES (?, ?, ?, ?, ?, ?)", userid, photo, thumbnail, extension, size, album)
	if err != nil {
		Logger.WriteFatal(err.Error())
	}
}

func (d *DB) CloseDB() {
	// Closes the opened db ...

	if err := d.db.Close(); err != nil {
		Logger.WriteFatal(err.Error())
	}
}

func New() *DB {
	db := new(DB)
	db.CreateDB()
	db.CreateTables()
	return db
}
