package caching

const (
	GET_PHOTOS                  = "GetPhotos"
	GET_VIDEOS                  = "GetVideos"
	GET_PHOTOS_FROM_ALBUM       = "GetPhotosFromAlbum"
	GET_VIDEOS_FROM_ALBUM       = "GetVideosFromAlbum"
	GET_ALBUMS                  = "GetAlbums"
	GET_USER_INFO               = "GetUserinfo"
	GET_FULL_MEDIA_BY_THUMBNAIL = "GetFullMediaByThumbnail"

	ErrConverting = "an error happened during the converting"
)

type GetPhotosModel struct {
	Thumbnail []byte
	Tags      []string
}

type GetVideosModel struct {
	Thumbnail []byte
	Tags      []string
}

type GetPhotosFromAlbum struct {
	Thumbnail []byte
	Tags      []string
}

type GetVideosFromAlbum struct {
	Thumbnail []byte
	Tags      []string
}

type GetAlbumsModel struct {
	Name                string
	LatestPhoto          []byte
}

type GetUserinfoModel struct {
	Firstname  string
	Secondname string
	Storage    float64
}

type GetFullMediaByThumbnail struct {
	Media []byte
}
