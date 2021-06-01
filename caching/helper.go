package caching

import (
	"NewPhoto/log"
	"encoding/json"
	"errors"
)

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

type Definer interface {
	Define(string, string) (interface{}, error)
}

type D struct{}

type GetPhotosModel struct {
	Photo     []byte
	Thumbnail []byte
	Extension string
	Size      float64
	Tags      string
}

type GetVideosModel struct {
	Video []byte
}

type GetPhotosFromAlbum struct {
	Photo     []byte
	Thumbnail []byte
	Extension string
	Size      float64
	Tags      string
	Album     string
}

type GetVideosFromAlbum struct {
	Video     []byte
	Extension string
	Album     string
}

type GetAlbumsModel struct {
	Name                 string
	LatestPhoto          []byte
	LatestPhotoThumbnail []byte
}

type GetUserinfoModel struct {
	Firstname  string
	Secondname string
	Storage    float64
}

type GetFullMediaByThumbnail struct {
	Media []byte
}

func (d *D) Define(c string, args string) (interface{}, error) {
	// Defines type of command .. then unparses it's values ...

	var stat interface{} = 0
	switch c {
	case GET_PHOTOS:
		stat = new([]GetPhotosModel)
	case GET_VIDEOS:
		stat = new([]GetVideosModel)
	case GET_VIDEOS_FROM_ALBUM:
		stat = new([]GetVideosFromAlbum)
	case GET_PHOTOS_FROM_ALBUM:
		stat = new([]GetPhotosFromAlbum)
	case GET_ALBUMS:
		stat = new([]GetAlbumsModel)
	case GET_USER_INFO:
		stat = new(GetUserinfoModel)
	case GET_FULL_MEDIA_BY_THUMBNAIL:
		stat = new(GetFullMediaByThumbnail)
	}
	if err := json.Unmarshal([]byte(args), &stat); err != nil {
		log.Logger.UsingErrorLogFile().CFatalln("CacheDefine", err)
	}
	var err error
	if _, ok := stat.(int); ok {
		err = errors.New("type is not defined")
	}
	return stat, err
}

func NewDefiner() Definer {
	return new(D)
}

type DataConfigurator interface {
	Configure(interface{}) string
}

type DC struct{}

func (dc *DC) Configure(data interface{}) string {
	result, err := json.Marshal(data)
	if err != nil {
		log.Logger.UsingErrorLogFile().CFatalln("CacheConfigure", err)
	}
	return string(result)
}

func NewConfigurator() DataConfigurator {
	return new(DC)
}
