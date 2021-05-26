package db

type Service interface {
	CreateDB()
	CreateTables()
	CloseDB()
}

type Auth interface {
	Login(login, pass string) (string, error)
	Logout(userid string) error
	RegisterUser(login, pass, firstname, secondname string) error
	IsLogin(uid string) bool
}

type Home interface {
	GetPhotos(userid string) []GetPhotosModel
	GetVideos(userid string) []GetVideosModel
	UploadPhoto(userid string, photo, thumbnail []byte, extension string, size float64, tags []string)
	UploadVideo(userid, extension string, video []byte, size float64)
}

type Account interface {
	GetUserinfo(userid string) (string, string, float64)
	GetUserAvatar(userid string) []byte
	SetUserAvatar(userid string, avatar []byte)
}

type Album interface {
	GetAlbums(userid string) []GetAlbumsModel
	GetPhotosFromAlbum(userid, name string) []GetPhotosFromAlbumModel
	GetVideosFromAlbum(userid, name string) []GetVideosFromAlbumModel
	UploadPhotoToAlbum(userid, extension, album string, size float64, photo, thumbnail []byte)
	UploadVideoToAlbum(userid, extension, album string, video []byte, size float64)
	CreateAlbum(userid, name string) bool
	DeleteAlbum(userid, name string) bool
	DeletePhotoFromAlbum(userid, album string, photo []byte)
	DeleteVideoFromAlbum(userid, album string, video []byte)
	GetAlbumInfo(userid, album string) int64
}

type Util interface {
	GetFullPhotoByThumbnail(userid string, thumbnail []byte) []byte
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

	//Provides access to util funcs(getting of full photo by thumbnail or sth else)
	Util
}