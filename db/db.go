package db

import (
	"errors"
	"fmt"
	"math"
	"os"

	"NewPhoto/log"
	"NewPhoto/utils"

	"database/sql"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

var (
	DefaultCapacity = 15 * math.Pow(10, 9)
)

type DB struct {
	db *sqlx.DB
}

func (d *DB) CreateDB() {

	password, ok := os.LookupEnv("dbPassword")
	if !ok {
		log.Logger.Fatalln("dbPassword is not written in credentials.sh file")
	}
	username, ok := os.LookupEnv("dbUsername")
	if !ok {
		log.Logger.Fatalln("dbUsername is not written in credentials.sh file")
	}

	dbDatabase, ok := os.LookupEnv("dbDatabase")
	if !ok {
		log.Logger.Fatalln("dbTable is not written in credentials.sh file")
	}

	addr, ok := os.LookupEnv("dbAddr")
	if !ok {
		log.Logger.Fatalln("dbAddr is not written in credentials.sh file")
	}

	db, err := sqlx.Connect("postgres", fmt.Sprintf("user=%s dbname=%s host=%s password=%s sslmode=disable", username, dbDatabase, addr, password))
	if err != nil {
		log.Logger.Fatalln(err)
	}
	d.db = db
}

func (d *DB) CreateTables() {
	_, err := d.db.Exec(fmt.Sprintf("CREATE TABLE IF NOT EXISTS users (userid text, firstname text, secondname text, storage float8 DEFAULT %f, avatar bytea, PRIMARY KEY(userid))", DefaultCapacity))
	if err != nil {
		log.Logger.Fatalln(err.Error())
	}
	_, err = d.db.Exec("CREATE TABLE IF NOT EXISTS token (userid text, token uuid, UNIQUE(token), FOREIGN KEY(userid) REFERENCES users(userid) ON DELETE CASCADE)")
	if err != nil {
		log.Logger.Fatalln(err.Error())
	}
	_, err = d.db.Exec("CREATE TABLE IF NOT EXISTS photos (userid text, photo bytea, thumbnail bytea, extension text, size float8, album text[], tags text[], FOREIGN KEY(userid) REFERENCES users(userid))")
	if err != nil {
		log.Logger.Fatalln(err.Error())
	}
	_, err = d.db.Exec("CREATE TABLE IF NOT EXISTS videos (userid text, video bytea, extension text, size float8, album text[], FOREIGN KEY(userid) REFERENCES users(userid))")
	if err != nil {
		log.Logger.Fatalln(err.Error())
	}
	_, err = d.db.Exec("CREATE TABLE IF NOT EXISTS albums (userid text, name text, FOREIGN KEY(userid) REFERENCES users(userid))")
	if err != nil {
		log.Logger.Fatalln(err.Error())
	}
}

func (d *DB) CloseDB() {
	if err := d.db.Close(); err != nil {
		log.Logger.WithFields(logrus.Fields{"qt": "CloseDB"}).Fatalln(err)
	}
}

func (d *DB) Login(login, pass string) (string, error) {
	c := string(utils.EncodeLogin(login, pass))
	var userid string
	tx, err := d.db.Begin()
	if err != nil {
		log.Logger.WithFields(logrus.Fields{"qt": "Login"}).Fatalln(err)
	}
	if err := d.db.Get(&userid, "SELECT userid FROM users WHERE userid = $1", c); err != nil {
		log.Logger.WithFields(logrus.Fields{"qt": "Login"}).Fatalln(err)
	}
	if len(userid) != 0 {
		t := uuid.NewString()
		if _, err := d.db.Exec("INSERT INTO token (userid, token) VALUES ($1, $2)", c, t); err != nil {
			log.Logger.WithFields(logrus.Fields{"qt": "Login"}).Fatalln(err)
			return c, nil
		}
	}
	if err := tx.Commit(); err != nil {
		log.Logger.WithFields(logrus.Fields{"qt": "Login"}).Fatalln(err)
	}
	return "", errors.New(LOGIN_ERROR)
}

func (d *DB) Logout(userid string) error {
	tx, err := d.db.Begin()
	if err != nil {
		log.Logger.WithFields(logrus.Fields{"qt": "Logout"}).Fatalln(err)
	}
	if _, err := d.db.Exec("DELETE FROM token WHERE userid = $1", userid); err != nil {
		log.Logger.WithFields(logrus.Fields{"qt": "Logout"}).Fatalln(err)
	}
	if err := tx.Commit(); err != nil {
		log.Logger.WithFields(logrus.Fields{"qt": "Logout"}).Fatalln(err)
	}
	return nil
}

func (d *DB) RegisterUser(login, pass, firstname, secondname string) error {
	_, err := d.db.Exec("INSERT INTO users (userid, firstname, secondname) VALUES ($1, $2, $3)", utils.EncodeLogin(login, pass), firstname, secondname)
	if err != nil {
		if pqerr, ok := err.(*pq.Error); ok {
			if pqerr.Code == "23505" {
				return errors.New(REGISTRAION_ERROR)
			}
		}
		log.Logger.WithFields(logrus.Fields{"qt": "RegisterUser"}).Fatalln(err)
	}
	return nil
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

func (d *DB) GetPhotos(userid string) []GetPhotosModel {
	r := []GetPhotosModel{}
	if err := d.db.Select(&r, "SELECT photo, thumbnail, extension, size, tags FROM photos WHERE userid = $1", userid); err != nil {
		log.Logger.WithFields(logrus.Fields{"qt": "GetPhotos"}).Fatalln(err)
	}
	return r
}

func (d *DB) GetVideos(userid string) []GetVideosModel {
	r := []GetVideosModel{}
	if err := d.db.Select(&r, "SELECT video FROM videos WHERE userid = $1", userid); err != nil {
		log.Logger.WithFields(logrus.Fields{"qt": "GetVideos"}).Fatalln(err)
	}
	return r
}

func (d *DB) UploadPhoto(userid string, photo, thumbnail []byte, extension string, size float64, tags []string) {
	d.db.MustExec("INSERT INTO photos (userid, photo, thumbnail, extension, size, tags) VALUES ($1, $2, $3, $4, $5, $6)", userid, photo, thumbnail, extension, size, pq.Array(tags))
}

func (d *DB) UploadVideo(userid, extension string, video []byte, size float64) {
	d.db.MustExec("INSERT INTO videos (userid, video, extension, size) VALUES ($1, $2, $3, $4)", userid, video, extension, size)
}

func (d *DB) GetUserinfo(userid string) (string, string, float64) {
	// Returns all the available storage for uploading by the passed user ...

	var storage float64
	if err := d.db.Select(&storage, "SELECT storage FROM users WHERE userid = $1", userid); err != nil {
		log.Logger.WithFields(logrus.Fields{"qt": "GetUserinfo"}).Fatalln(err)
	}

	var size float64
	if err := d.db.Select(&size, "SELECT SUM(size) FROM photos WHERE userid = $1", userid); err != nil {
		log.Logger.WithFields(logrus.Fields{"qt": "GetUserinfo"}).Fatalln(err)
	}

	var Userinfo struct {
		firstname  string
		secondname string
	}
	if err := d.db.Select(&Userinfo, "SELECT firstname, secondname FROM users WHERE userid = $1", userid); err != nil {
		log.Logger.WithFields(logrus.Fields{"qt": "GetUserinfo"}).Fatalln(err)
	}

	return Userinfo.firstname, Userinfo.secondname, storage - size
}

func (d *DB) GetUserAvatar(userid string) []byte {
	var avatar []interface{}
	if err := d.db.Select(&avatar, "SELECT avatar FROM users WHERE userid = $1", userid); err != nil {
		log.Logger.WithFields(logrus.Fields{"qt": "GetUserAvatar"}).Fatalln(err)
	}
	switch x := avatar[0].(type) {
	case []byte:
		return []byte(x)
	}
	return []byte("")
}

func (d *DB) SetUserAvatar(userid string, avatar []byte) {
	d.db.MustExec("UPDATE users SET avatar = $1 WHERE userid = $2", avatar, userid)
}

func AlbumInResult(w string, c []GetAlbumsModel) bool {
	for _, v := range c {
		if v.Album == w {
			return true
		}
	}
	return false
}

func (d *DB) GetAlbums(userid string) []GetAlbumsModel {
	var albums []string
	if err := d.db.Select(&albums, "SELECT name FROM albums WHERE userid = $1", userid); err != nil {
		log.Logger.WithFields(logrus.Fields{"qt": "GetAlbums"}).Fatalln(err)
	}

	r := []GetAlbumsModel{}
	for _, v := range albums {
		a := GetAlbumsModel{}
		if err := d.db.Get(&a, "SELECT album, photo FROM photos WHERE userid = $1 AND $2=ANY(album) ORDER BY photo DESC LIMIT 1", userid, v); err != nil && err != sql.ErrNoRows {
			log.Logger.WithFields(logrus.Fields{"qt": "GetAlbums"}).Fatalln(err)
		}
		a.Album = v
		r = append(r, a)
	}
	return r
}

func (d *DB) GetPhotosFromAlbum(userid, name string) []GetPhotosFromAlbumModel {
	r := []GetPhotosFromAlbumModel{}
	if err := d.db.Select(&r, `SELECT userid, photo, thumbnail, extension, size, album FROM photos WHERE userid = $1 AND $2=ANY(album)`, userid, name); err != nil {
		log.Logger.WithFields(logrus.Fields{"qt": "GetPhotosFromAlbum"}).Fatalln(err)
	}
	return r
}

func (d *DB) GetVideosFromAlbum(userid, name string) []GetVideosFromAlbumModel {
	r := []GetVideosFromAlbumModel{}
	if err := d.db.Select(&r, "SELECT video, extension FROM videos WHERE userid = $1 AND $2=ANY(album)", userid, name); err != nil {
		log.Logger.WithFields(logrus.Fields{"qt": "GetVideosFromAlbum"}).Fatalln(err)
	}
	return r
}

func (d *DB) UploadPhotoToAlbum(userid, extension, album string, size float64, photo, thumbnail []byte) {
	_, err := d.db.Exec(
		"INSERT INTO photos (userid, photo, thumbnail, extension, size) (SELECT $1, $2, $3, $4, $5 WHERE NOT EXISTS (SELECT photo FROM photos WHERE userid = $1 AND photo = $2))", userid, photo, thumbnail, extension, size)
	if err != nil {
		log.Logger.WithFields(logrus.Fields{"qt": "UploadPhotoToAlbum"}).Fatalln(err)
	}
	if _, err := d.db.Exec("UPDATE photos SET album = array_append(album, $1) WHERE userid = $2 AND photo = $3", album, userid, photo); err != nil {
		log.Logger.WithFields(logrus.Fields{"qt": "UploadPhotoToAlbum"}).Fatalln(err)
	}
}

func (d *DB) UploadVideoToAlbum(userid, extension, album string, video []byte, size float64) {
	err := d.db.MustExec(
		"INSERT INTO videos (userid, video, extension, size) (SELECT $1, $2, $3, $4 WHERE NOT EXISTS (SELECT video FROM videos WHERE userid = $1 AND video = $2))", userid, video, extension, size)
	if err != nil {
		log.Logger.WithFields(logrus.Fields{"qt": "UploadVideoToAlbum"}).Fatalln(err)
	}
	if err := d.db.MustExec("UPDATE videos SET album = array_append(album, $1) WHERE userid = $2 AND video = $3", album, userid, video); err != nil {
		log.Logger.WithFields(logrus.Fields{"qt": "UploadVideoToAlbum"}).Fatalln(err)
	}
}

func (d *DB) CreateAlbum(userid, name string) bool {
	d.db.MustExec("INSERT INTO albums (userid, name) (SELECT $1, $2 WHERE NOT EXISTS (SELECT * FROM albums WHERE userid = $1 AND name = $2))", userid, name)
	return true
}

func (d *DB) DeleteAlbum(userid, album string) bool {
	d.db.MustExec("DELETE FROM albums WHERE userid = $1 AND name = $2", userid, album)
	d.db.MustExec("UPDATE photos videos set album = array_remove(album, $1) WHERE userid = $2", album, userid)
	return true
}

func (d *DB) DeletePhotoFromAlbum(userid, album string, photo []byte) {
	d.db.MustExec("UPDATE photos SET album = array_remove(album, $1) WHERE userid = $2, AND photo = $3", album, userid, photo)
}

func (d *DB) DeleteVideoFromAlbum(userid, album string, video []byte) {
	d.db.MustExec("UPDATE videos SET album = array_remove(album, $1) WHERE userid = $2, AND video = $3", album, userid, video)
}

func (d *DB) GetAlbumInfo(userid, album string) int64 {
	var mediaNumPhotos []interface{}
	if err := d.db.Select(&mediaNumPhotos, "SELECT COUNT(*) FROM photos WHERE userid = $1 AND $2=ANY(album)", userid, album); err != nil {
		log.Logger.WithFields(logrus.Fields{"qt": "GetAlbumInfo"}).Fatalln(err)
	}

	var mediaNumVideos []interface{}
	if err := d.db.Select(&mediaNumVideos, "SELECT COUNT(*) FROM videos WHERE userid = $1 AND $2=ANY(album)", userid, album); err != nil {
		log.Logger.WithFields(logrus.Fields{"qt": "GetAlbumInfo"}).Fatalln(err)
	}
	return mediaNumPhotos[0].(int64) + mediaNumVideos[0].(int64)
}

func (d *DB) GetFullPhotoByThumbnail(userid string, thumbnail []byte) []byte {
	var photo []byte
	if err := d.db.Get(&photo, "SELECT photo FROM photos WHERE userid = $1 AND thumbnail = $2", userid, thumbnail); err != nil {
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
