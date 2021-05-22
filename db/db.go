package db

import (
	"NewPhoto/utils"

	"database/sql"
	"errors"
	"fmt"
	"math"
	"os"
	"strings"

	"NewPhoto/log"

	"github.com/go-sql-driver/mysql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/sirupsen/logrus"
)

var (
	DefaultCapacity = 15 * math.Pow(10, 9)
)

type DB struct {
	db *sql.DB
}

type IDB interface {
	CreateDB()
	CreateTables()
	CloseDB()

	LoginUser(login, pass string) (string, error)
	RegisterUser(login, pass, firstname, secondname string) bool
	IsLogin(uid string) bool

	GetPhotos(userid string) [][5]interface{}
	GetVideos(userid string) [][1]interface{}
	UploadPhoto(userid string, photo, thumbnail []byte, extension string, size float64, tags string)
	UploadVideo(userid, extension string, video []byte, size float64)

	GetUserinfo(userid string) (string, string, float64)
	GetUserAvatar(userid string) []byte
	SetUserAvatar(userid string, avatar []byte)

	GetAlbums(userid string) [][3]interface{}
	GetPhotosFromAlbum(userid, name string) [][5]interface{}
	GetVideosFromAlbum(userid, name string) [][2]interface{}
	UploadPhotoToAlbum(userid, extension, album string, size float64, photo, thumbnail []byte)
	UploadVideoToAlbum(userid, album, extension string, video []byte, size float64)
	CreateAlbum(userid, name string) bool
	DeleteAlbum(userid, name string) bool
	DeletePhotoFromAlbum(userid, album string, photo []byte)
	DeleteVideoFromAlbum(userid, album string, video []byte)
	GetAlbumInfo(userid, album string) int64

	GetFullPhotoByThumbnail(userid string, thumbnail []byte) []byte
}

func (d *DB) CreateDB() {

	password, ok := os.LookupEnv("mysqlPassword")
	if !ok {
		log.Logger.Fatalln("mysqlPassword is not written in credentials.sh file")
	}
	username, ok := os.LookupEnv("mysqlUsername")
	if !ok {
		log.Logger.Fatalln("mysqlUsername is not written in credentials.sh file")
	}

	table, ok := os.LookupEnv("mysqlTable")
	if !ok {
		log.Logger.Fatalln("mysqlTable is not written in credentials.sh file")
	}

	addr, ok := os.LookupEnv("mysqlAddr")
	if !ok {
		log.Logger.Fatalln("mysqlAddr is not written in credentials.sh file")
	}

	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s)/%s", username, password, addr, table))
	if err != nil {
		log.Logger.Fatalln(err)
	}
	d.db = db
}

func (d *DB) CreateTables() {

	_, err := d.db.Exec(fmt.Sprintf("CREATE TABLE IF NOT EXISTS users (userid VARCHAR(200) UNIQUE, firstname VARCHAR(200), secondname VARCHAR(200), storage DOUBLE DEFAULT %f, avatar MEDIUMBLOB)", DefaultCapacity))
	if err != nil {
		log.Logger.Fatalln(err.Error())
	}
	_, err = d.db.Exec("CREATE TABLE IF NOT EXISTS photos (userid VARCHAR(200), photo MEDIUMBLOB, thumbnail MEDIUMBLOB, extension VARCHAR(10), size DOUBLE, album VARCHAR(200), tags VARCHAR(200))")
	if err != nil {
		log.Logger.Fatalln(err.Error())
	}
	_, err = d.db.Exec("CREATE TABLE IF NOT EXISTS videos (userid VARCHAR(200), video LONGBLOB, extension VARCHAR(10), size DOUBLE, album VARCHAR(200))")
	if err != nil {
		log.Logger.Fatalln(err.Error())
	}
	_, err = d.db.Exec("CREATE TABLE IF NOT EXISTS albums (userid VARCHAR(200), name VARCHAR(200))")
	if err != nil {
		log.Logger.Fatalln(err.Error())
	}
}

func (d *DB) CloseDB() {
	if err := d.db.Close(); err != nil {
		log.Logger.WithFields(logrus.Fields{"qt": "CloseDB"}).Fatalln(err)
	}
}

func (d *DB) LoginUser(login, pass string) (string, error) {
	passedEncodedCredentials := utils.EncodeLogin(login, pass)
	rows, err := d.db.Query("SELECT userid FROM users")
	if err != nil {
		log.Logger.WithFields(logrus.Fields{"qt": "LoginUser"}).Fatalln(err)
	}
	defer func() {
		err := rows.Close()
		if err != nil {
			log.Logger.WithFields(logrus.Fields{"qt": "LoginUser"}).Fatalln(err)
		}
	}()

	for rows.Next() {
		var userid string
		if err := rows.Scan(&userid); err != nil {
			log.Logger.WithFields(logrus.Fields{"qt": "LoginUser"}).Fatalln(err)
		}
		if userid == string(passedEncodedCredentials) {
			return string(passedEncodedCredentials), nil
		}
	}
	return "", errors.New("such user does not exist")
}

func (d *DB) RegisterUser(login, pass, firstname, secondname string) bool {

	_, err := d.db.Exec("INSERT INTO users (userid, firstname, secondname) VALUES (?, ?, ?)", utils.EncodeLogin(login, pass), firstname, secondname)
	if err != nil {
		mysqlError := err.(*mysql.MySQLError)
		if mysqlError.Number == 1062 {
			return false
		}
		log.Logger.WithFields(logrus.Fields{"qt": "RegisterUser"}).Fatalln(err)
	}
	return true
}

func (d *DB) IsLogin(uid string) bool {
	rows, err := d.db.Query("SELECT userid FROM users")
	if err != nil {
		log.Logger.WithFields(logrus.Fields{"qt": "IsLogin"}).Fatalln(err)
	}
	defer func() {
		err := rows.Close()
		if err != nil {
			log.Logger.WithFields(logrus.Fields{"qt": "IsLogin"}).Fatalln(err)
		}
	}()

	for rows.Next() {
		var userid string
		if err := rows.Scan(&userid); err != nil {
			log.Logger.WithFields(logrus.Fields{"qt": "IsLogin"}).Fatalln(err)
		}
		if userid == uid {
			return true
		}
	}
	return false
}

func (d *DB) GetPhotos(userid string) [][5]interface{} {
	// Returns all the photos uploaded by user ...

	rows, err := d.db.Query("SELECT * FROM photos WHERE userid = ?", userid)
	if err != nil {
		log.Logger.WithFields(logrus.Fields{"qt": "GetPhotos"}).Fatalln(err)
	}
	defer func() {
		err := rows.Close()
		if err != nil {
			log.Logger.WithFields(logrus.Fields{"qt": "GetPhotos"}).Fatalln(err)
		}
	}()
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
		if err := rows.Scan(&model.Userid, &model.Photo, &model.Thumbnail, &model.Extension, &model.Size, &model.Album, &model.Tags); err != nil{
			log.Logger.WithFields(logrus.Fields{"qt": "GetPhotos"}).Fatalln(err)
		}
		result = append(result, [5]interface{}{model.Photo, model.Thumbnail, model.Extension, model.Size, model.Tags.String})
	}
	return result
}

func (d *DB) GetVideos(userid string) [][1]interface{} {
	rows, err := d.db.Query("SELECT * FROM videos WHERE userid = ?", userid)
	if err != nil {
		log.Logger.WithFields(logrus.Fields{"qt": "GetVideos"}).Fatalln(err)
	}
	defer func() {
		err := rows.Close()
		if err != nil {
			log.Logger.WithFields(logrus.Fields{"qt": "GetVideos"}).Fatalln(err)
		}
	}()
	model := struct {
		Userid string
		Video  []byte
		Album  sql.NullString
	}{}
	var result [][1]interface{}
	for rows.Next() {
		if err := rows.Scan(&model.Userid, &model.Video, &model.Album); err != nil{
			log.Logger.WithFields(logrus.Fields{"qt": "GetVideos"}).Fatalln(err)
		}
		result = append(result, [1]interface{}{&model.Video})
	}
	return result
}

func (d *DB) UploadPhoto(userid string, photo, thumbnail []byte, extension string, size float64, tags string) {
	_, err := d.db.Exec("INSERT INTO photos (userid, photo, thumbnail, extension, size, tags) VALUES (?, ?, ?, ?, ?, ?)", userid, photo, thumbnail, extension, size, tags)
	if err != nil {
		log.Logger.WithFields(logrus.Fields{"qt": "UploadPhoto"}).Fatalln(err)
	}
}

func (d *DB) UploadVideo(userid, extension string, video []byte, size float64) {
	_, err := d.db.Exec("INSERT INTO videos (userid, video, extension, size) VALUES (?, ?, ?, ?)", userid, video, extension, size)
	if err != nil {
		log.Logger.WithFields(logrus.Fields{"qt": "UploadVideo"}).Fatalln(err)
	}
}

func (d *DB) GetUserinfo(userid string) (string, string, float64) {
	// Returns all the available storage for uploading by the passed user ...

	row := d.db.QueryRow("SELECT storage FROM users WHERE userid = ?", userid)
	var storage float64
	err := row.Scan(&storage)
	if err != nil {
		log.Logger.WithFields(logrus.Fields{"qt": "GetUserinfo"}).Fatalln(err)
	}

	rows, err := d.db.Query("SELECT size FROM photos WHERE userid = ?", userid)
	if err != nil {
		log.Logger.WithFields(logrus.Fields{"qt": "GetUserinfo"}).Fatalln(err)
	}
	defer func() {
		err := rows.Close()
		if err != nil {
			log.Logger.WithFields(logrus.Fields{"qt": "GetUserinfo"}).Fatalln(err)
		}
	}()

	var allphotossize float64
	for rows.Next() {
		var storage float64
		if err := rows.Scan(&storage); err != nil {
			log.Logger.WithFields(logrus.Fields{"qt": "GetUserinfo"}).Fatalln(err)
		}
		allphotossize += storage
	}

	row = d.db.QueryRow("SELECT firstname, secondname FROM users WHERE userid = ?", userid)

	var firstname string
	var secondname string

	if err := row.Scan(&firstname, &secondname); err != nil {
		log.Logger.WithFields(logrus.Fields{"qt": "GetUserinfo"}).Fatalln(err)
	}
	return firstname, secondname, storage - allphotossize
}

func (d *DB) GetUserAvatar(userid string) []byte {
	row := d.db.QueryRow("SELECT avatar FROM users WHERE userid = ?", userid)

	var avatar []byte

	if err := row.Scan(&avatar); err != nil {
		log.Logger.WithFields(logrus.Fields{"qt": "GetUserinfo"}).Fatalln(err)
	}
	return avatar
}

func (d *DB) SetUserAvatar(userid string, avatar []byte) {
	_, err := d.db.Exec("UPDATE users SET avatar = ? WHERE userid = ?", avatar, userid)
	if err != nil {
		log.Logger.WithFields(logrus.Fields{"qt": "SetUserAvatar"}).Fatalln(err)
	}
}

func (d *DB) GetAlbums(userid string) [][3]interface{} {
	rows, err := d.db.Query("SELECT name FROM albums WHERE userid = ?", userid)
	if err != nil {
		log.Logger.WithFields(logrus.Fields{"qt": "GetAlbums"}).Fatalln(err)
	}

	defer func() {
		err := rows.Close()
		if err != nil {
			log.Logger.WithFields(logrus.Fields{"qt": "GetAlbums"}).Fatalln(err)
		}
	}()

	var albums []string

	for rows.Next() {
		var album string
		if err := rows.Scan(&album); err != nil{
			log.Logger.WithFields(logrus.Fields{"qt": "GetAlbums"}).Fatalln(err)
		}
		albums = append(albums, album)
	}

	var result [][3]interface{}
	for _, album := range albums {
		row := d.db.QueryRow("SELECT photo, thumbnail FROM photos WHERE userid = ? AND (album like ? OR album like ? OR album like ?) ORDER BY userid DESC LIMIT 1", userid, "%"+album+",%", "%,"+album+"%", "%,"+album+",%")
		var photo []byte
		var thumbnail []byte
		if err := row.Scan(&photo, &thumbnail); err != nil && err != sql.ErrNoRows {
			log.Logger.WithFields(logrus.Fields{"qt": "GetAlbums"}).Fatalln(err)
		}
		result = append(result, [3]interface{}{album, photo, thumbnail})
	}
	return result
}

func (d *DB) GetPhotosFromAlbum(userid, name string) [][5]interface{} {
	rows, err := d.db.Query(`SELECT userid, photo, thumbnail, extension, size, album FROM photos WHERE userid = ? AND (album like ? OR album like ? OR album like ?)`, userid, "%"+name+",%", "%,"+name+"%", "%,"+name+",%")
	if err != nil {
		log.Logger.WithFields(logrus.Fields{"qt": "GetPhotosFromAlbum"}).Fatalln(err)
	}
	defer func() {
		err := rows.Close()
		if err != nil {
			log.Logger.WithFields(logrus.Fields{"qt": "GetPhotosFromAlbum"}).Fatalln(err)
		}
	}()
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
			log.Logger.WithFields(logrus.Fields{"qt": "GetPhotosFromAlbum"}).Fatalln(err)
		}
		result = append(result, [5]interface{}{model.Photo, model.Thumbnail, model.Extension, model.Size, model.Album.String})
	}
	return result
}

func (d *DB) GetVideosFromAlbum(userid, name string) [][2]interface{} {
	rows, err := d.db.Query("SELECT * FROM videos WHERE userid = ? AND (album like ? OR album like ? OR album like ?)", userid, "%"+name+",%", "%,"+name+"%", "%,"+name+",%")
	if err != nil {
		log.Logger.WithFields(logrus.Fields{"qt": "GetVideosFromAlbum"}).Fatalln(err)
	}
	defer func() {
		err := rows.Close()
		if err != nil {
			log.Logger.WithFields(logrus.Fields{"qt": "GetVideosFromAlbum"}).Fatalln(err)
		}
	}()
	model := struct {
		Userid    string
		Video     []byte
		Extension string
		Size      float64
		Album     sql.NullString
	}{}
	var result [][2]interface{}
	for rows.Next() {
		err := rows.Scan(&model.Userid, &model.Video, &model.Extension, &model.Size, &model.Album)
		if err != nil {
			log.Logger.WithFields(logrus.Fields{"qt": "GetVideosFromAlbum"}).Fatalln(err)
		}
		result = append(result, [2]interface{}{model.Video, model.Extension})
	}
	return result
}

func (d *DB) UploadPhotoToAlbum(userid, extension, album string, size float64, photo, thumbnail []byte) {
	_, err := d.db.Exec(
		"INSERT INTO photos (userid, photo, thumbnail, extension, size, album) SELECT ?, ?, ?, ?, ?, ? FROM DUAL WHERE NOT EXISTS (SELECT * FROM photos WHERE userid = ? AND photo = ?) ", userid, photo, thumbnail, extension, size, album+",", userid, photo)
	if err != nil {
		log.Logger.WithFields(logrus.Fields{"qt": "UploadPhotoToAlbum"}).Fatalln(err)
	}
	row := d.db.QueryRow("SELECT album FROM photos WHERE userid = ? AND photo = ?", userid, photo)
	var a sql.NullString
	if err := row.Scan(&a); err != nil {
		log.Logger.WithFields(logrus.Fields{"qt": "UploadPhotoToAlbum"}).Fatalln(err)
	}

	split := strings.Split(a.String, ",")
	if len(split) == 2 && split[0] == album {
		return
	}

	a.String = fmt.Sprintf("%s%s,", a.String, album)
	if _, err := d.db.Exec("UPDATE photos SET album = ? WHERE userid = ? AND photo = ?", a.String, userid, photo); err != nil {
		log.Logger.WithFields(logrus.Fields{"qt": "UploadPhotoToAlbum"}).Fatalln(err)
	}
}

func (d *DB) UploadVideoToAlbum(userid, album, extension string, video []byte, size float64) {
	_, err := d.db.Exec(
		"INSERT INTO videos (userid, video, extension, size, album) SELECT ?, ?, ?, ?, ? FROM DUAL WHERE NOT EXISTS (SELECT * FROM videos WHERE userid = ? AND video = ?) ", userid, video, extension, size, album+",", userid, video)
	if err != nil {
		log.Logger.WithFields(logrus.Fields{"qt": "UploadVideoToAlbum"}).Fatalln(err)
	}

	row := d.db.QueryRow("SELECT album FROM videos WHERE userid = ? AND video = ?", userid, video)
	var a sql.NullString
	if err := row.Scan(&a); err != nil {
		log.Logger.WithFields(logrus.Fields{"qt": "UploadVideoToAlbum"}).Fatalln(err)
	}
	split := strings.Split(a.String, ",")
	if len(split) == 2 && split[0] == album {
		return
	}

	a.String = fmt.Sprintf("%s%s,", a.String, album)
	if _, err := d.db.Exec("UPDATE videos SET album = ? WHERE userid = ? AND video = ?", a.String, userid, video); err != nil {
		log.Logger.WithFields(logrus.Fields{"qt": "UploadVideoToAlbum"}).Fatalln(err)
	}
}

func (d *DB) CreateAlbum(userid, name string) bool {
	_, err := d.db.Exec("INSERT IGNORE INTO albums (userid, name) VALUES(?, ?)", userid, name)
	if err != nil {
		log.Logger.WithFields(logrus.Fields{"qt": "CreateAlbum"}).Fatalln(err)
	}
	return true
}

func (d *DB) DeleteAlbum(userid, name string) bool {
	_, err := d.db.Exec("DELETE FROM albums WHERE userid = ? AND name = ?", userid, name)
	if err != nil {
		log.Logger.WithFields(logrus.Fields{"qt": "DeleteAlbum"}).Fatalln(err)
	}
	return true
}

func (d *DB) DeletePhotoFromAlbum(userid, album string, photo []byte) {
	row := d.db.QueryRow("SELECT album FROM photos WHERE userid = ? AND photo = ?", userid, photo)
	var a string
	if err := row.Scan(&a); err != nil {
		log.Logger.WithFields(logrus.Fields{"qt": "DeletePhotoFromAlbum"}).Fatalln(err)
	}
	split := strings.Split(a, ",")
	for i, v := range split {
		if v == album {
			if len(split) == 2 && len(split[1]) == 0 {
				a = ""
			} else {
				a = strings.Join(append(split[i:], split[i+1:]...), ",")
			}
			break
		}
	}

	if _, err := d.db.Exec("UPDATE photos SET album = ? WHERE userid = ? AND photo = ?", a, userid, photo); err != nil {
		log.Logger.WithFields(logrus.Fields{"qt": "DeletePhotoFromAlbum"}).Fatalln(err)
	}
}

func (d *DB) DeleteVideoFromAlbum(userid, album string, video []byte) {
	row := d.db.QueryRow("SELECT album FROM videos WHERE userid = ? AND video = ?", userid, video)
	var a string
	if err := row.Scan(&a); err != nil {
		log.Logger.WithFields(logrus.Fields{"qt": "DeleteVideoFromAlbum"}).Fatalln(err)
	}
	split := strings.Split(a, ",")
	for i, v := range split {
		if v == album {
			a = strings.Join(append(split[i:], split[i+1:]...), ",")
			break
		}
	}

	if _, err := d.db.Exec("UPDATE videos SET album = ? WHERE userid = ? AND video = ?", a, userid, video); err != nil {
		log.Logger.WithFields(logrus.Fields{"qt": "DeleteVideoFromAlbum"}).Fatalln(err)
	}
}

func (d *DB) GetAlbumInfo(userid, album string) int64 {
	row := d.db.QueryRow("SELECT COUNT(*) FROM photos as p WHERE p.userid = ? AND (p.album like ? OR p.album like ? OR p.album like ?)", userid, "%"+album+",%", "%,"+album+"%", "%,"+album+",%")
	var mediaNumPhotos int64
	if err := row.Scan(&mediaNumPhotos); err != nil {
		log.Logger.WithFields(logrus.Fields{"qt": "GetAlbumInfo"}).Fatalln(err)
	}

	row = d.db.QueryRow("SELECT COUNT(*) FROM videos as v WHERE v.userid = ? AND (v.album like ? OR v.album like ? OR v.album like ?)", userid, "%"+album+",%", "%,"+album+"%", "%,"+album+",%")
	var mediaNumVideos int64
	if err := row.Scan(&mediaNumVideos); err != nil {
		log.Logger.WithFields(logrus.Fields{"qt": "GetAlbumInfo"}).Fatalln(err)
	}
	return mediaNumPhotos + mediaNumVideos
}

func (d *DB) GetFullPhotoByThumbnail(userid string, thumbnail []byte) []byte {
	row := d.db.QueryRow("SELECT photo FROM photos WHERE userid = ? AND thumbnail = ?", userid, thumbnail)
	var photo []byte
	if err := row.Scan(&photo); err != nil {
		log.Logger.WithFields(logrus.Fields{"qt": "GetFullPhotoByThumbnail"}).Fatalln(err)
	}
	return photo
}

func New() *DB {
	db := new(DB)
	db.CreateDB()
	db.CreateTables()
	return db
}
