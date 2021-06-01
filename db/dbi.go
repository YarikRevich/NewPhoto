package db

import "google.golang.org/protobuf/reflect/protoreflect"

type IMediaSize interface {
	//Implements the methods to get the with and height
	//of the passed media

	GetWidth() float64
	GetHeight() float64
}

type MediaSize struct {
	Width  float64
	Height float64
}

func (i *MediaSize) GetWidth() float64 {
	return i.Width
}

func (i *MediaSize) GetHeight() float64 {
	return i.Height
}

type IMediaType interface {
	//Gets the opportunity to get media type(enum)

	GetScanData() (string, string)
}

type MediaType protoreflect.EnumNumber

const (
	Photo MediaType = iota
	Video
)

func (m MediaType) GetScanData() (string, string) {
	switch m {
	case Photo:
		return "photo", "photos"
	case Video:
		return "video", "videos"
	}
	return "", ""
}

type Service interface {
	CreateDB()
	CreateTables()
	CloseDB()
}

type Auth interface {
	//Checks if user's cred are ok and then
	//Returns his tokens
	Login(login, pass, sourceType string) (string, string, error)

	//Logouts user and then removes tokens
	Logout(userid, sourceType string) error

	//Gets accessToken and loginToken. Then checks them
	//If it is ok it returns ok status and new tokens
	IsTokenCorrect(accessToken, loginToken, sourceType string) bool

	//Gets userID using accessToken and loginToken
	GetUserID(accessToken, loginToken string) string

	RegisterUser(login, pass, firstname, secondname string) error
}

type Home interface {
	GetPhotos(userid string, offset, page int64) []GetPhotosModel
	GetVideos(userid string) []GetVideosModel
	UploadPhoto(userid string, photo, thumbnail []byte, extension string, size float64, tags []string)
	UploadVideo(userid, extension string, video []byte, size float64)
}

type Account interface {
	DeleteAccount(userid string)
	GetUserinfo(userid string) (string, string, float64)
	GetUserAvatar(userid string) []byte
	SetUserAvatar(userid string, avatar []byte)
}

type Album interface {
	//Creates album due to the given name
	CreateAlbum(userid, name string) bool

	//Deletes album due to the given name
	DeleteAlbum(userid, name string) bool

	//Gets album due to the given name
	GetAlbums(userid string) []GetAlbumsModel

	//Gets photos from the album due to the given name
	GetPhotosFromAlbum(userid, name string, offset, page int64) []GetPhotosFromAlbumModel

	//Gets videos from the album due to the given name
	GetVideosFromAlbum(userid, name string, offset, page int64) []GetVideosFromAlbumModel

	//Uploads video to the album due to the given name
	UploadPhotoToAlbum(userid, extension, album string, size float64, photo, thumbnail []byte)

	//Uploads video to the album due to the given name
	UploadVideoToAlbum(userid, extension, album string, video []byte, size float64)

	//Deletes photo from the album due to the given name
	DeletePhotoFromAlbum(userid, album string, photo []byte)

	//Deletes video from the album due to the given name
	DeleteVideoFromAlbum(userid, album string, video []byte)
}

type Info interface {
	//Gets num of user's photos
	GetPhotosNum(userid string) int64

	//Gets num of user's videos
	GetVideosNum(userid string) int64

	//Gets num of user's photos in the album
	GetPhotosInAlbumNum(userid, name string) int64

	//Gets num of user's videos in the album
	GetVideosInAlbumNum(userid, name string) int64

	//Gets num of user's albums
	GetAlbumsNum(userid string) int64
}

type Util interface {
	//Will load media which is bigger in size
	//It will get it by the thumbnail and will resize it on the fly
	GetFullMediaByThumbnail(userid string, thumbnail []byte, mediaSize IMediaSize, mediaType IMediaType) []byte
}

type IDB interface {

	//Provides access to system funcs
	Service

	//Provides access to auth funcs (login, reg, isLogin)
	Auth

	//Provides access to photo and video funcs
	Home

	//Provides access to account funcs
	Account

	//Provides access to album funcs(general of equal)
	Album

	//Provides access to get info about user's data
	Info

	//Provides access to util funcs(getting of full photo by thumbnail or sth else)
	Util
}
