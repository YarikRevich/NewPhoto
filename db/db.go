package db

import (
	"errors"
	"fmt"
	"math"
	"os"

	"github.com/YarikRevich/NewPhoto/exceptions"
	"github.com/YarikRevich/NewPhoto/log"
	"github.com/YarikRevich/NewPhoto/utils"

	"database/sql"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"

	"github.com/google/uuid"
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
		log.Logger.UsingErrorLogFile().CFatalln("InitDB", "dbPassword is not written in credentials.sh file")
	}
	username, ok := os.LookupEnv("dbUsername")
	if !ok {
		log.Logger.UsingErrorLogFile().CFatalln("InitDB", "dbUsername is not written in credentials.sh file")
	}

	dbDatabase, ok := os.LookupEnv("dbDatabase")
	if !ok {
		log.Logger.UsingErrorLogFile().CFatalln("InitDB", "dbTable is not written in credentials.sh file")
	}

	addr, ok := os.LookupEnv("dbAddr")
	if !ok {
		log.Logger.UsingErrorLogFile().CFatalln("InitDB", "dbAddr is not written in credentials.sh file")
	}

	db, err := sqlx.Connect("postgres", fmt.Sprintf("user=%s dbname=%s host=%s password=%s sslmode=disable", username, dbDatabase, addr, password))
	if err != nil {
		log.Logger.UsingErrorLogFile().CFatalln("InitDB", err)
	}
	d.db = db
}

func (d *DB) CreateTables() {
	_, err := d.db.Exec(fmt.Sprintf("CREATE TABLE IF NOT EXISTS users (userid text, firstname text, secondname text, storage float8 DEFAULT %f, avatar bytea, PRIMARY KEY(userid))", DefaultCapacity))
	if err != nil {
		log.Logger.UsingErrorLogFile().CFatalln("CreateDBs", err)
	}
	t, err := d.db.Begin()
	if err != nil {
		log.Logger.UsingErrorLogFile().CFatalln("CreateDBs", err)
	}
	_, err = d.db.Exec("CREATE TYPE source_type as ENUM ('web', 'mobile')")
	if err != nil {
		if p, ok := err.(*pq.Error); ok && p.Code != "42710" {
			log.Logger.UsingErrorLogFile().CFatalln("CreateDBs", err)
		}
	}
	if err := t.Commit(); err != nil {
		log.Logger.UsingErrorLogFile().CFatalln("CreateDBs", err)
	}
	_, err = d.db.Exec("CREATE TABLE IF NOT EXISTS token (userid text, access_token uuid, login_token uuid, source_type source_type, UNIQUE(userid ,access_token, login_token), FOREIGN KEY(userid) REFERENCES users(userid) ON DELETE CASCADE)")
	if err != nil {
		log.Logger.UsingErrorLogFile().CFatalln("CreateDBs", err)
	}
	_, err = d.db.Exec("CREATE TABLE IF NOT EXISTS photos (userid text, photo bytea, thumbnail bytea, extension text, size float8, album text[], tags text[], FOREIGN KEY(userid) REFERENCES users(userid) ON DELETE CASCADE)")
	if err != nil {
		log.Logger.UsingErrorLogFile().CFatalln("CreateDBs", err)
	}
	_, err = d.db.Exec("CREATE TABLE IF NOT EXISTS videos (userid text, video bytea, thumbnail bytea, extension text, size float8, album text[], tags text[], FOREIGN KEY(userid) REFERENCES users(userid) ON DELETE CASCADE)")
	if err != nil {
		log.Logger.UsingErrorLogFile().CFatalln("CreateDBs", err)
	}
	_, err = d.db.Exec("CREATE TABLE IF NOT EXISTS albums (userid text, name text, FOREIGN KEY(userid) REFERENCES users(userid) ON DELETE CASCADE)")
	if err != nil {
		log.Logger.UsingErrorLogFile().CFatalln("CreateDBs", err)
	}
}

func (d *DB) CloseDB() {
	if err := d.db.Close(); err != nil {
		log.Logger.UsingErrorLogFile().CFatalln("CloseDB", err)
	}
}

func (d *DB) Login(login, pass string, sourceType ISourceType) (string, string, error) {
	c := string(utils.EncodeLogin(login, pass))
	sT := sourceType.GetScanData()
	var userid string
	tx, err := d.db.Begin()
	if err != nil {
		log.Logger.UsingErrorLogFile().CFatalln("Login", err)
	}
	if err := d.db.Get(&userid, "SELECT userid FROM users WHERE userid = $1", c); err != nil && err != sql.ErrNoRows {
		log.Logger.UsingErrorLogFile().CFatalln("Login", err)
	}
	if len(userid) != 0 {
		var newAccessToken string
		var newLoginToken string
		for {
			newAccessToken = uuid.NewString()
			var e bool
			if err := d.db.Get(&e, "SELECT COUNT((SELECT access_token FROM token WHERE access_token = $1 AND source_type = $2)) = 0", newAccessToken, sT); err != nil && err != sql.ErrNoRows {
				log.Logger.UsingErrorLogFile().CFatalln("Login", err)
			}
			if e {
				break
			}
		}
		for {
			newLoginToken = uuid.NewString()
			var e bool
			if err := d.db.Get(&e, "SELECT COUNT((SELECT login_token FROM token WHERE login_token = $1 AND source_type=$2)) = 0", newLoginToken, sT); err != nil && err != sql.ErrNoRows {
				log.Logger.UsingErrorLogFile().CFatalln("Login", err)
			}
			if e {
				break
			}
		}
		if _, err := d.db.Exec("INSERT INTO token (userid, access_token, login_token, source_type) VALUES($1, $2, $3, $4)", c, newAccessToken, newLoginToken, sT); err != nil {
			log.Logger.UsingErrorLogFile().CFatalln("Login", err)
		}
		return newAccessToken, newLoginToken, nil
	}
	if err := tx.Commit(); err != nil {
		log.Logger.UsingErrorLogFile().CFatalln("Login", err)
	}
	return "", "", errors.New(exceptions.LOGIN_ERROR)
}

func (d *DB) Logout(userid string, sourceType ISourceType) error {
	sT := sourceType.GetScanData()
	tx, err := d.db.Begin()
	if err != nil {
		log.Logger.UsingErrorLogFile().CFatalln("Logout", err)
	}
	if _, err := d.db.Exec("DELETE FROM token WHERE userid = $1 AND source_type=$2", userid, sT); err != nil {
		log.Logger.UsingErrorLogFile().CFatalln("Logout", err)
	}
	if err := tx.Commit(); err != nil {
		log.Logger.UsingErrorLogFile().CFatalln("Logout", err)
	}
	return nil
}

func (d *DB) IsTokenCorrect(accessToken, loginToken string, sourceType ISourceType) bool {
	sT := sourceType.GetScanData()
	var ok bool
	if err := d.db.Get(&ok, "SELECT COUNT((SELECT userid from token WHERE access_token = $1 AND login_token = $2 AND source_type = $3)) != 0", accessToken, loginToken, sT); err != nil {
		log.Logger.UsingErrorLogFile().CFatalln("IsTokenCorrect", err)
	}
	return ok
}

func (d *DB) GetUserID(accessToken, loginToken string) string {
	var userID string
	if err := d.db.Get(&userID, "SELECT userid FROM token WHERE access_token = $1 AND login_token = $2", accessToken, loginToken); err != nil {
		log.Logger.UsingErrorLogFile().CFatalln("GetUserID", err)
	}
	return userID
}

func (d *DB) RegisterUser(login, pass, firstname, secondname string) error {
	_, err := d.db.Exec("INSERT INTO users (userid, firstname, secondname) VALUES ($1, $2, $3)", string(utils.EncodeLogin(login, pass)), firstname, secondname)
	if err != nil {
		if pqerr, ok := err.(*pq.Error); ok {
			if pqerr.Code == "23505" {
				return errors.New(exceptions.REGISTRAION_ERROR)
			}
		}
		log.Logger.UsingErrorLogFile().CFatalln("RegisterUser", err)
	}
	return nil
}

func (d *DB) GetPhotos(userid string, offset, page int64) []GetPhotosModel {
	r := []GetPhotosModel{}
	if err := d.db.Select(&r, "SELECT photo, thumbnail, extension, size, tags FROM photos WHERE userid = $1", userid); err != nil {
		log.Logger.UsingErrorLogFile().CFatalln("GetPhotos", err)
	}
	if int64(len(r)) >= (page-1)*offset+offset {
		return r[(page-1)*offset : (page-1)*offset+offset]
	}
	return r
}

func (d *DB) GetPhotosNum(userid string) int64 {
	var num int64
	if err := d.db.Get(&num, "SELECT COUNT(*) FROM photos WHERE userid = $1", userid); err != nil {
		log.Logger.UsingErrorLogFile().CFatalln("GetPhotosNum", err)
	}
	return num
}

func (d *DB) GetVideos(userid string, offset, page int64) []GetVideosModel {
	r := []GetVideosModel{}
	if err := d.db.Select(&r, "SELECT thumbnail FROM videos WHERE userid = $1", userid); err != nil {
		log.Logger.UsingErrorLogFile().CFatalln("GetVideos", err)
	}
	if int64(len(r)) >= (page-1)*offset+offset {
		return r[(page-1)*offset : (page-1)*offset+offset]
	}
	return r
}

func (d *DB) GetVideosNum(userid string) int64 {
	var num int64
	if err := d.db.Get(&num, "SELECT COUNT(*) FROM videos WHERE userid = $1", userid); err != nil {
		log.Logger.UsingErrorLogFile().CFatalln("GetVideosNum", err)
	}
	return num
}

func (d *DB) UploadPhoto(userid string, photo, thumbnail []byte, extension string, size float64, tags []string) {
	d.db.MustExec("INSERT INTO photos (userid, photo, thumbnail, extension, size, tags) VALUES ($1, $2, $3, $4, $5, $6)", userid, photo, thumbnail, extension, size, pq.Array(tags))
}

func (d *DB) UploadVideo(userid, extension string, video, thumbnail []byte, size float64, tags []string) {
	d.db.MustExec("INSERT INTO videos (userid, video, thumbnail, extension, size, tags) VALUES ($1, $2, $3, $4, $5, $6)", userid, video, thumbnail, extension, size, pq.Array(tags))
}

func (d *DB) DeleteAccount(userid string) {
	if _, err := d.db.Exec("DELETE FROM users WHERE userid = $1", userid); err != nil {
		log.Logger.CFatalln("DeleteAccount", err)
	}
}

func (d *DB) GetUserinfo(userid string) (string, string, float64) {
	// Returns all the available storage for uploading by the passed user ...

	var storage float64
	if err := d.db.Get(&storage, "SELECT storage FROM users WHERE userid = $1", userid); err != nil {
		log.Logger.UsingErrorLogFile().CFatalln("GetUserinfo", err)
	}

	var size sql.NullFloat64
	if err := d.db.Get(&size, "SELECT SUM(size) FROM photos WHERE userid = $1", userid); err != nil {
		log.Logger.UsingErrorLogFile().CFatalln("GetUserinfo", err)
	}

	var Userinfo struct {
		Firstname  string
		Secondname string
	}
	if err := d.db.Get(&Userinfo, "SELECT firstname, secondname FROM users WHERE userid = $1", userid); err != nil {
		log.Logger.UsingErrorLogFile().CFatalln("GetUserinfo", err)
	}

	return Userinfo.Firstname, Userinfo.Secondname, storage - size.Float64
}

func (d *DB) GetUserAvatar(userid string) []byte {
	var avatar []interface{}
	if err := d.db.Select(&avatar, "SELECT avatar FROM users WHERE userid = $1", userid); err != nil {
		log.Logger.UsingErrorLogFile().CFatalln("GetUserAvatar", err)
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
		if v.Name == w {
			return true
		}
	}
	return false
}

func (d *DB) GetAlbums(userid string) []GetAlbumsModel {
	var albums []string
	if err := d.db.Select(&albums, "SELECT name FROM albums WHERE userid = $1", userid); err != nil {
		log.Logger.UsingErrorLogFile().CFatalln("GetAlbums", err)
	}

	r := []GetAlbumsModel{}
	for _, v := range albums {
		a := GetAlbumsModel{}
		if err := d.db.Get(&a, "SELECT album, photo FROM photos WHERE userid = $1 AND $2=ANY(album) ORDER BY photo DESC LIMIT 1", userid, v); err != nil && err != sql.ErrNoRows {
			log.Logger.UsingErrorLogFile().CFatalln("GetAlbums", err)
		}
		a.Name = v
		r = append(r, a)
	}
	return r
}

func (d *DB) GetPhotosFromAlbum(userid, name string, offset, page int64) []GetPhotosFromAlbumModel {
	r := []GetPhotosFromAlbumModel{}
	if err := d.db.Select(&r, `SELECT thumbnail, tags FROM photos WHERE userid = $1 AND $2=ANY(album)`, userid, name); err != nil {
		log.Logger.UsingErrorLogFile().CFatalln("GetPhotosFromAlbum", err)
	}
	if int64(len(r)) >= (page-1)*offset+offset {
		return r[(page-1)*offset : (page-1)*offset+offset]
	}
	return r
}

func (d *DB) GetPhotosInAlbumNum(userid, name string) int64 {
	var num int64
	if err := d.db.Get(&num, "SELECT COUNT(*) FROM photos WHERE userid = $1 AND $2=ANY(album)", userid, name); err != nil {
		log.Logger.UsingErrorLogFile().CFatalln("GetPhotosInAlbumNum", err)
	}
	return num
}

func (d *DB) GetVideosFromAlbum(userid, name string, offset, page int64) []GetVideosFromAlbumModel {
	r := []GetVideosFromAlbumModel{}
	if err := d.db.Select(&r, "SELECT thumbnail, tags FROM videos WHERE userid = $1 AND $2=ANY(album)", userid, name); err != nil {
		log.Logger.UsingErrorLogFile().CFatalln("GetVideosFromAlbum", err)
	}
	if int64(len(r)) >= (page-1)*offset+offset {
		return r[(page-1)*offset : (page-1)*offset+offset]
	}
	return r
}

func (d *DB) GetVideosInAlbumNum(userid, name string) int64 {
	var num int64
	if err := d.db.Get(&num, "SELECT COUNT(*) FROM videos WHERE userid = $1 AND $2=ANY(album)", userid, name); err != nil {
		log.Logger.UsingErrorLogFile().CFatalln("GetVideosInAlbumNum", err)
	}
	return num
}

func (d *DB) UploadPhotoToAlbum(userid, extension, album string, size float64, photo, thumbnail []byte, tags []string) {
	_, err := d.db.Exec(
		"INSERT INTO photos (userid, photo, thumbnail, extension, size) (SELECT $1, $2, $3, $4, $5 WHERE NOT EXISTS (SELECT photo FROM photos WHERE userid = $1 AND photo = $2))", userid, photo, thumbnail, extension, size)
	if err != nil {
		log.Logger.UsingErrorLogFile().CFatalln("UploadPhotoToAlbum", err)
	}
	if _, err := d.db.Exec("UPDATE photos SET album = array_append(album, $1), tags = $2 WHERE userid = $3 AND photo = $4", album, pq.Array(tags), userid, photo); err != nil {
		log.Logger.UsingErrorLogFile().CFatalln("UploadPhotoToAlbum", err)
	}
}

func (d *DB) UploadVideoToAlbum(userid, extension, album string, video, thumbnail []byte, size float64, tags []string) {
	if _, err := d.db.Exec(
		"INSERT INTO videos (userid, video, thumbnail, extension, size) (SELECT $1, $2, $3, $4, $5 WHERE NOT EXISTS (SELECT video FROM videos WHERE userid = $1 AND video = $2))", userid, video, thumbnail, extension, size); err != nil {
		log.Logger.UsingErrorLogFile().CFatalln("UploadVideoToAlbum", err)
	}
	if _, err := d.db.Exec("UPDATE videos SET album = array_append(album, $1), tags = $2 WHERE userid = $3 AND video = $4", album, pq.Array(tags), userid, video); err != nil {
		log.Logger.UsingErrorLogFile().CFatalln("UploadVideoToAlbum", err)
	}
}

func (d *DB) CreateAlbum(userid, name string) bool {
	d.db.MustExec("INSERT INTO albums (userid, name) (SELECT $1, $2 WHERE NOT EXISTS (SELECT * FROM albums WHERE userid = $1 AND name = $2))", userid, name)
	return true
}

func (d *DB) DeleteAlbum(userid, album string) bool {
	t, err := d.db.Begin()
	if err != nil {
		log.Logger.UsingErrorLogFile().CFatalln("DeleteAlbum", err)
	}
	defer func() {
		if err := t.Commit(); err != nil {
			log.Logger.UsingErrorLogFile().CFatalln("DeleteAlbum", err)
		}
	}()
	d.db.MustExec("DELETE FROM albums WHERE userid = $1 AND name = $2", userid, album)
	d.db.MustExec("UPDATE photos videos set album = array_remove(album, $1) WHERE userid = $2", album, userid)

	return true
}

func (d *DB) DeletePhotoFromAlbum(userid, album string, photo []byte) {
	d.db.MustExec("UPDATE photos SET album = array_remove(album, $1) WHERE userid = $2 AND photo = $3", album, userid, photo)
}

func (d *DB) DeleteVideoFromAlbum(userid, album string, video []byte) {
	d.db.MustExec("UPDATE videos SET album = array_remove(album, $1) WHERE userid = $2 AND video = $3", album, userid, video)
}

func (d *DB) GetAlbumsNum(userid string) int64 {
	var num int64
	if err := d.db.Get(&num, "SELECT COUNT(*) FROM albums WHERE userid = $1", userid); err != nil {
		log.Logger.UsingErrorLogFile().CFatalln("GetAlbumsNum", err)
	}
	return num
}

func (d *DB) GetFullMediaByThumbnail(userid string, thumbnail []byte, mediaSize IMediaSize, mediaType IMediaType) []byte {
	r, t := mediaType.GetScanData()

	var media []byte
	if err := d.db.Get(&media, fmt.Sprintf("SELECT %s FROM %s WHERE userid = $1 AND thumbnail = $2", r, t), userid, thumbnail); err != nil {
		log.Logger.UsingErrorLogFile().CFatalln("GetFullMediaByThumbnail", err)
	}
	return media
}

func New() *DB {
	db := new(DB)
	db.CreateDB()
	db.CreateTables()
	return db
}
