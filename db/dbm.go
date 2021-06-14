//Contains models to parse the results from db
package db

import "github.com/lib/pq"

//Parses the result from GetAlbums request
type GetAlbumsModel struct {
	Name  string `db:"album"`
	Photo []byte `db:"photo"`
}

//Parses the result from GetPhotosFromAlbum request
type GetPhotosFromAlbumModel struct {
	Thumbnail []byte         `db:"thumbnail"`
	Tags      pq.StringArray `db:"tags"`
}

//Parses the result from GetVideosFromAlbum request
type GetVideosFromAlbumModel struct {
	Thumbnail []byte         `db:"thumbnail"`
	Tags      pq.StringArray `db:"tags"`
}

//Parses the result from GetPhotos request
type GetPhotosModel struct {
	Photo     []byte         `db:"photo"`
	Thumbnail []byte         `db:"thumbnail"`
	Extension string         `db:"extension"`
	Size      float64        `db:"size"`
	Tags      pq.StringArray `db:"tags"`
}

//Parses the result from GetVideos request
type GetVideosModel struct {
	Thumbnail []byte         `db:"thumbnail"`
	Tags      pq.StringArray `db:"tags"`
}
