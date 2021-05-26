//Contains models to parse the results from db
package db

import "database/sql"

//Parses the result from GetAlbums request
type GetAlbumsModel struct {
	Album string `db:"album"`
	Photo []byte `db:"photo"`
}

//Parses the result from GetPhotosFromAlbum request
type GetPhotosFromAlbumModel struct {
	Userid    string         `db:"userid"`
	Photo     []byte         `db:"photo"`
	Thumbnail []byte         `db:"thumbnail"`
	Extension string         `db:"extension"`
	Size      float64        `db:"size"`
	Album     sql.NullString `db:"album"`
}

//Parses the result from GetVideosFromAlbum request
type GetVideosFromAlbumModel struct {
	Video     []byte `db:"video"`
	Extension string `db:"extension"`
}

//Parses the result from GetPhotos request
type GetPhotosModel struct {
	Photo     []byte         `db:"photo"`
	Thumbnail []byte         `db:"thumbnail"`
	Extension string         `db:"extension"`
	Size      float64        `db:"size"`
	Tags      sql.NullString `db:"tags"`
}

//Parses the result from GetVideos request
type GetVideosModel struct {
	Video []byte `db:"video"`
}
