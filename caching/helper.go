package caching

import (
	"encoding/json"
	"errors"
	"log"
)

type Definer interface {
	Define(string, string) (interface{}, error)
}

type D struct{}

// Model to unparse AllPhotos cached values
type AllPhotosModel struct {
	Photo     []byte
	Thumbnail []byte
	Extension string
	Size      float64
	Tags      string
}

type AllPhotosAlbumModel struct {
	Photo     []byte
	Thumbnail []byte
	Extension string
	Size      float64
	Album     string
}

type GetAllAlbumsModel struct {
	Name                 string
	LatestPhoto          []byte
	LatestPhotoThumbnail []byte
}

type GetUserinfoModel struct {
	Firstname  string
	Secondname string
	Storage    float64
}

type GetFullPhotoByThumbnail struct {
	Photo []byte
}

func (d *D) Define(c string, args string) (interface{}, error) {
	// Defines type of command .. then unparses it's values ...

	switch c {
	case "AllPhotos":
		var stat []AllPhotosModel
		json.Unmarshal([]byte(args), &stat)
		return stat, nil
	case "AllPhotosAlbum":
		var stat []AllPhotosAlbumModel
		json.Unmarshal([]byte(args), &stat)
		return stat, nil
	case "GetAllAlbums":
		var stat []GetAllAlbumsModel
		json.Unmarshal([]byte(args), &stat)
		return stat, nil
	case "GetUserinfo":
		var stat GetUserinfoModel
		json.Unmarshal([]byte(args), &stat)
		return stat, nil
	case "GetFullPhotoByThumbnail":
		var stat GetFullPhotoByThumbnail
		json.Unmarshal([]byte(args), &stat)
		return stat, nil
	default:
		return nil, errors.New("type is not defined")
	}
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
		log.Fatalln(err)
	}
	return string(result)
}

func NewConfigurator() DataConfigurator {
	return new(DC)
}
