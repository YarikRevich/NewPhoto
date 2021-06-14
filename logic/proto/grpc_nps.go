package proto

import (
	"context"
	"github.com/YarikRevich/NewPhoto/caching"
	"github.com/YarikRevich/NewPhoto/db"
	"github.com/YarikRevich/NewPhoto/log"
)

type NewPhoto struct {
	Tag        TagClient
	DBInstanse db.IDB
}

func (s *NewPhoto) GetPhotos(r *GetPhotosRequest, stream NewPhotos_GetPhotosServer) error {
	log.Logger.UsingAccessLogFile().CInfoln("GetPhotos")
	select {
	case <-stream.Context().Done():
		return nil
	default:
		if cr, cached := caching.RedisInstanse.IsCached(s.DBInstanse.GetUserID(r.GetAccessToken(), r.GetLoginToken()), caching.GET_PHOTOS); cached {
			definer := caching.Definer{Model: []caching.GetPhotosModel{}, Data: cr}
			result := definer.Define()
			converted, ok := result.([]caching.GetPhotosModel)
			if !ok {
				log.Logger.UsingErrorLogFile().CFatalln("GetPhotos", caching.ErrConverting)
			}
			for _, value := range converted {
				if err := stream.Send(&GetPhotosResponse{Thumbnail: value.Thumbnail, Tags: value.Tags, Ok: true}); err != nil {
					continue
				}
			}
			return nil
		}

		result := s.DBInstanse.GetPhotos(s.DBInstanse.GetUserID(r.GetAccessToken(), r.GetLoginToken()), r.GetOffset(), r.GetPage())

		var model []caching.GetPhotosModel
		for _, value := range result {
			model = append(model, caching.GetPhotosModel{Thumbnail: value.Thumbnail, Tags: []string{}})
			if err := stream.Send(&GetPhotosResponse{Thumbnail: value.Thumbnail, Tags: []string{}, Ok: true}); err != nil {
				continue
			}
		}
		conf := caching.DataConfigurator{Model: model}
		caching.RedisInstanse.Set(s.DBInstanse.GetUserID(r.GetAccessToken(), r.GetLoginToken()), caching.GET_PHOTOS, string(conf.Configure()))
		return nil
	}
}

func (s *NewPhoto) GetPhotosNum(ctx context.Context, r *GetPhotosNumRequest) (*GetPhotosNumResponse, error) {
	log.Logger.UsingAccessLogFile().CInfoln("GetPhotosNum")

	n := s.DBInstanse.GetPhotosNum(s.DBInstanse.GetUserID(r.GetAccessToken(), r.GetLoginToken()))
	return &GetPhotosNumResponse{Ok: true, Num: n}, nil
}

func (s *NewPhoto) GetVideos(r *GetVideosRequest, stream NewPhotos_GetVideosServer) error {
	log.Logger.UsingAccessLogFile().CInfoln("GetVideos")
	select {
	case <-stream.Context().Done():
		return nil
	default:
		if cr, cached := caching.RedisInstanse.IsCached(s.DBInstanse.GetUserID(r.GetAccessToken(), r.GetLoginToken()), caching.GET_VIDEOS); cached {
			definer := caching.Definer{Model: []caching.GetVideosModel{}, Data: cr}
			result := definer.Define()
			converted, ok := result.([]caching.GetVideosModel)
			if !ok {
				log.Logger.UsingErrorLogFile().CFatalln("GetVideos", caching.ErrConverting)
			}
			for _, value := range converted {
				if err := stream.Send(&GetVideosResponse{Thumbnail: value.Thumbnail, Tags: value.Tags, Ok: true}); err != nil {
					continue
				}
			}
			return nil
		}

		result := s.DBInstanse.GetVideos(s.DBInstanse.GetUserID(r.GetAccessToken(), r.GetLoginToken()), r.GetOffset(), r.GetPage())

		var model []caching.GetVideosModel
		for _, value := range result {
			model = append(model, caching.GetVideosModel{Thumbnail: value.Thumbnail})
			if err := stream.Send(&GetVideosResponse{Thumbnail: value.Thumbnail, Ok: true}); err != nil {
				continue
			}
		}
		conf := caching.DataConfigurator{Model: model}
		caching.RedisInstanse.Set(s.DBInstanse.GetUserID(r.GetAccessToken(), r.GetLoginToken()), caching.GET_VIDEOS, string(conf.Configure()))

		return nil
	}
}

func (s *NewPhoto) GetVideosNum(ctx context.Context, r *GetVideosNumRequest) (*GetVideosNumResponse, error) {
	log.Logger.UsingAccessLogFile().CInfoln("GetVideosNum")

	n := s.DBInstanse.GetVideosNum(s.DBInstanse.GetUserID(r.GetAccessToken(), r.GetLoginToken()))
	return &GetVideosNumResponse{Ok: true, Num: n}, nil
}

func (s *NewPhoto) UploadPhoto(stream NewPhotos_UploadPhotoServer) error {
	log.Logger.UsingAccessLogFile().CInfoln("UploadPhoto")
	select {
	case <-stream.Context().Done():
		return nil
	default:
		for {
			msg, err := stream.Recv()
			if err != nil {
				break
			}
			// tags, err := s.Tag.RecognizeObject(context.Background(), &RecognizeObjectRequest{Photo: msg.GetPhoto()})
			// if err != nil {
			// 	log.Logger.Fatalln(err)
			// }
			s.DBInstanse.UploadPhoto(s.DBInstanse.GetUserID(msg.GetAccessToken(), msg.GetLoginToken()), msg.GetPhoto(), msg.GetThumbnail(), msg.GetExtension(), msg.GetSize(), []string{})
		}
		if err := stream.SendAndClose(&UploadPhotoResponse{Ok: true}); err != nil {
			log.Logger.UsingErrorLogFile().CFatalln("UploadPhoto", err)
		}
	}
	caching.RedisInstanse.Clean(caching.GET_PHOTOS)
	return nil
}

func (s *NewPhoto) UploadVideo(stream NewPhotos_UploadVideoServer) error {
	log.Logger.UsingAccessLogFile().CInfoln("UploadVideo")
	select {
	case <-stream.Context().Done():
		return nil
	default:
		for {
			msg, err := stream.Recv()
			if err != nil {
				break
			}
			s.DBInstanse.UploadVideo(s.DBInstanse.GetUserID(msg.GetAccessToken(), msg.GetLoginToken()), msg.GetExtension(), msg.GetVideo(), msg.GetThumbnail(), msg.GetSize(), []string{})
		}
		if err := stream.SendAndClose(&UploadVideoResponse{Ok: true}); err != nil {
			log.Logger.UsingErrorLogFile().CFatalln("UploadVideo", err)
		}
	}
	caching.RedisInstanse.Clean(caching.GET_VIDEOS)
	return nil
}

func (s *NewPhoto) DeleteAccount(ctx context.Context, r *DeleteAccountRequest) (*DeleteAccountResponse, error) {
	log.Logger.UsingAccessLogFile().CInfoln("DeleteAccount")

	s.DBInstanse.DeleteAccount(s.DBInstanse.GetUserID(r.GetAccessToken(), r.GetLoginToken()))
	return &DeleteAccountResponse{Ok: true}, nil
}

func (s *NewPhoto) GetUserinfo(cxt context.Context, r *GetUserinfoRequest) (*GetUserinfoResponse, error) {
	log.Logger.UsingAccessLogFile().CInfoln("GetUserinfo")
	if cr, cached := caching.RedisInstanse.IsCached(s.DBInstanse.GetUserID(r.GetAccessToken(), r.GetLoginToken()), caching.GET_USER_INFO); cached {
		definer := caching.Definer{Model: caching.GetUserinfoModel{}, Data: cr}
		result := definer.Define()
		converted, ok := result.(caching.GetUserinfoModel)
		if !ok {
			log.Logger.UsingErrorLogFile().CFatalln("GetUserinfo", caching.ErrConverting)
		}

		return &GetUserinfoResponse{Firstname: converted.Firstname, Secondname: converted.Secondname, Storage: converted.Storage, Ok: true}, nil
	}

	firstname, secondname, storage := s.DBInstanse.GetUserinfo(s.DBInstanse.GetUserID(r.GetAccessToken(), r.GetLoginToken()))

	model := caching.GetUserinfoModel{Firstname: firstname, Secondname: secondname, Storage: storage}
	conf := caching.DataConfigurator{Model: model}
	caching.RedisInstanse.Set(s.DBInstanse.GetUserID(r.GetAccessToken(), r.GetLoginToken()), caching.GET_USER_INFO, string(conf.Configure()))

	return &GetUserinfoResponse{Firstname: firstname, Secondname: secondname, Storage: storage, Ok: true}, nil
}

func (s *NewPhoto) GetUserAvatar(ctx context.Context, r *GetUserAvatarRequest) (*GetUserAvatarResponse, error) {
	log.Logger.UsingAccessLogFile().CInfoln("GetUserAvatar")

	avatar := s.DBInstanse.GetUserAvatar(s.DBInstanse.GetUserID(r.GetAccessToken(), r.GetLoginToken()))
	return &GetUserAvatarResponse{Avatar: avatar, Ok: true}, nil
}

func (s *NewPhoto) SetUserAvatar(ctx context.Context, r *SetUserAvatarRequest) (*SetUserAvatarResponse, error) {
	log.Logger.UsingAccessLogFile().CInfoln("SetUserAvatar")

	s.DBInstanse.SetUserAvatar(s.DBInstanse.GetUserID(r.GetAccessToken(), r.GetLoginToken()), r.GetAvatar())
	return &SetUserAvatarResponse{Ok: true}, nil
}

func (s *NewPhoto) GetPhotosFromAlbum(r *GetPhotosFromAlbumRequest, stream NewPhotos_GetPhotosFromAlbumServer) error {
	log.Logger.UsingAccessLogFile().CInfoln("GetPhotosFromAlbum")

	select {
	case <-stream.Context().Done():
		return nil
	default:
		if cr, cached := caching.RedisInstanse.IsCached(s.DBInstanse.GetUserID(r.GetAccessToken(), r.GetLoginToken()), caching.GET_PHOTOS_FROM_ALBUM); cached {
			definer := caching.Definer{Model: []caching.GetPhotosFromAlbum{}, Data: cr}
			result := definer.Define()
			converted, ok := result.([]caching.GetPhotosFromAlbum)
			if !ok {
				log.Logger.UsingErrorLogFile().CFatalln("GetPhotosFromAlbum", caching.ErrConverting)
			}
			for _, value := range converted {
				if err := stream.Send(&GetPhotosFromAlbumResponse{Thumbnail: value.Thumbnail, Tags: value.Tags, Ok: true}); err != nil {
					continue
				}
			}
			return nil
		}

		result := s.DBInstanse.GetPhotosFromAlbum(s.DBInstanse.GetUserID(r.GetAccessToken(), r.GetLoginToken()), r.GetName(), r.GetOffset(), r.GetPage())
		var model []caching.GetPhotosFromAlbum
		for _, value := range result {
			model = append(model, caching.GetPhotosFromAlbum{Thumbnail: value.Thumbnail, Tags: value.Tags})
			if err := stream.Send(&GetPhotosFromAlbumResponse{Thumbnail: value.Thumbnail, Tags: value.Tags, Ok: true}); err != nil {
				log.Logger.UsingErrorLogFile().CFatalln("GetPhotosFromAlbum", err)
			}
		}
		conf := caching.DataConfigurator{Model: model}
		caching.RedisInstanse.Set(s.DBInstanse.GetUserID(r.GetAccessToken(), r.GetLoginToken()), caching.GET_PHOTOS_FROM_ALBUM, string(conf.Configure()))
		return nil
	}
}

func (s *NewPhoto) GetPhotosInAlbumNum(ctx context.Context, r *GetPhotosInAlbumNumRequest) (*GetPhotosInAlbumNumResponse, error) {
	log.Logger.UsingAccessLogFile().CInfoln("GetPhotosInAlbumNum")

	n := s.DBInstanse.GetPhotosInAlbumNum(s.DBInstanse.GetUserID(r.GetAccessToken(), r.GetLoginToken()), r.GetName())
	return &GetPhotosInAlbumNumResponse{Ok: true, Num: n}, nil
}

func (s *NewPhoto) GetVideosFromAlbum(r *GetVideosFromAlbumRequest, stream NewPhotos_GetVideosFromAlbumServer) error {
	log.Logger.UsingAccessLogFile().CInfoln("GetVideosFromAlbum")

	select {
	case <-stream.Context().Done():
		return nil
	default:
		if cr, cached := caching.RedisInstanse.IsCached(s.DBInstanse.GetUserID(r.GetAccessToken(), r.GetLoginToken()), caching.GET_VIDEOS_FROM_ALBUM); cached {
			definer := caching.Definer{Model: []caching.GetVideosFromAlbum{}, Data: cr}
			result := definer.Define()
			converted, ok := result.([]caching.GetVideosFromAlbum)
			if !ok {
				log.Logger.UsingErrorLogFile().CFatalln("GetVideosFromAlbum", caching.ErrConverting)
			}
			for _, value := range converted {
				if err := stream.Send(&GetVideosFromAlbumResponse{Thumbnail: value.Thumbnail, Tags: value.Tags, Ok: true}); err != nil {
					continue
				}
			}
			return nil
		}

		result := s.DBInstanse.GetVideosFromAlbum(s.DBInstanse.GetUserID(r.GetAccessToken(), r.GetLoginToken()), r.GetName(), r.GetOffset(), r.GetPage())

		var model []caching.GetVideosFromAlbum
		for _, value := range result {
			model = append(model, caching.GetVideosFromAlbum{Thumbnail: value.Thumbnail, Tags: value.Tags})
			if err := stream.Send(&GetVideosFromAlbumResponse{Thumbnail: value.Thumbnail, Tags: value.Tags, Ok: true}); err != nil {
				continue
			}
		}
		conf := caching.DataConfigurator{Model: model}
		caching.RedisInstanse.Set(s.DBInstanse.GetUserID(r.GetAccessToken(), r.GetLoginToken()), caching.GET_VIDEOS_FROM_ALBUM, string(conf.Configure()))

		return nil
	}
}

func (s *NewPhoto) GetVideosInAlbumNum(ctx context.Context, r *GetVideosInAlbumNumRequest) (*GetVideosInAlbumNumResponse, error) {
	log.Logger.UsingAccessLogFile().CInfoln("GetVideosInAlbumNum")

	n := s.DBInstanse.GetVideosInAlbumNum(s.DBInstanse.GetUserID(r.GetAccessToken(), r.GetLoginToken()), r.GetName())
	return &GetVideosInAlbumNumResponse{Ok: true, Num: n}, nil
}

func (s *NewPhoto) UploadPhotoToAlbum(stream NewPhotos_UploadPhotoToAlbumServer) error {
	log.Logger.UsingAccessLogFile().CInfoln("UploadPhotoToAlbum")

	select {
	case <-stream.Context().Done():
		return nil
	default:
		for {
			recv, err := stream.Recv()
			if err != nil {
				break
			}
			s.DBInstanse.UploadPhotoToAlbum(s.DBInstanse.GetUserID(recv.GetAccessToken(), recv.GetLoginToken()), recv.GetExtension(), recv.GetAlbum(), recv.GetSize(), recv.GetPhoto(), recv.GetThumbnail(), []string{})
		}
		if err := stream.SendAndClose(&UploadPhotoToAlbumResponse{Ok: true}); err != nil {
			log.Logger.UsingErrorLogFile().CFatalln("UploadPhotoToAlbum", err)
		}
	}
	caching.RedisInstanse.Clean(caching.GET_PHOTOS_FROM_ALBUM)
	return nil
}

func (s *NewPhoto) UploadVideoToAlbum(stream NewPhotos_UploadVideoToAlbumServer) error {
	log.Logger.UsingAccessLogFile().CInfoln("UploadVideoToAlbum")

	select {
	case <-stream.Context().Done():
		return nil
	default:
		for {
			recv, err := stream.Recv()
			if err != nil {
				break
			}
			s.DBInstanse.UploadVideoToAlbum(s.DBInstanse.GetUserID(recv.GetAccessToken(), recv.GetLoginToken()), recv.GetExtension(), recv.GetAlbum(), recv.GetVideo(), recv.GetThumbnail(), recv.GetSize(), []string{})
		}
		if err := stream.SendAndClose(&UploadVideoToAlbumResponse{Ok: true}); err != nil {
			log.Logger.UsingErrorLogFile().CFatalln("UploadVideoToAlbum", err)
		}
	}
	caching.RedisInstanse.Clean(caching.GET_VIDEOS_FROM_ALBUM)
	return nil
}

func (s *NewPhoto) GetAlbums(r *GetAlbumsRequest, stream NewPhotos_GetAlbumsServer) error {
	log.Logger.UsingAccessLogFile().CInfoln("GetAlbums")

	select {
	case <-stream.Context().Done():
		return nil
	default:
		if cr, cached := caching.RedisInstanse.IsCached(s.DBInstanse.GetUserID(r.GetAccessToken(), r.GetLoginToken()), caching.GET_ALBUMS); cached {
			definer := caching.Definer{Model: []caching.GetAlbumsModel{}, Data: cr}
			result := definer.Define()
			converted, ok := result.([]caching.GetAlbumsModel)
			if !ok {
				log.Logger.UsingErrorLogFile().CFatalln("GetAlbums", caching.ErrConverting)
			}
			for _, value := range converted {
				if err := stream.Send(&GetAlbumsResponse{Name: value.Name, LatestPhoto: value.LatestPhoto, Ok: true}); err != nil {
					continue
				}
			}
			return nil
		}

		result := s.DBInstanse.GetAlbums(s.DBInstanse.GetUserID(r.GetAccessToken(), r.GetLoginToken()))

		var model []caching.GetAlbumsModel

		for _, value := range result {
			model = append(model, caching.GetAlbumsModel{Name: value.Name, LatestPhoto: value.Photo})
			if err := stream.Send(&GetAlbumsResponse{Name: value.Name, LatestPhoto: value.Photo, Ok: true}); err != nil {
				continue
			}
		}
		conf := caching.DataConfigurator{Model: model}
		caching.RedisInstanse.Set(s.DBInstanse.GetUserID(r.GetAccessToken(), r.GetLoginToken()), caching.GET_ALBUMS, string(conf.Configure()))
		return nil
	}
}

func (s *NewPhoto) DeletePhotoFromAlbum(stream NewPhotos_DeletePhotoFromAlbumServer) error {
	log.Logger.UsingAccessLogFile().CInfoln("DeletePhotoFromAlbum")

	select {
	case <-stream.Context().Done():
		return nil
	default:
		for {
			msg, err := stream.Recv()
			if err != nil {
				break
			}
			s.DBInstanse.DeletePhotoFromAlbum(s.DBInstanse.GetUserID(msg.GetAccessToken(), msg.GetLoginToken()), msg.GetAlbum(), msg.GetPhoto())
		}
		if err := stream.SendAndClose(&DeletePhotoFromAlbumResponse{Ok: true}); err != nil {
			log.Logger.UsingErrorLogFile().CFatalln("DeletePhotoFromAlbum", err)
		}
	}
	caching.RedisInstanse.Clean(caching.GET_PHOTOS_FROM_ALBUM)
	return nil
}

func (s *NewPhoto) DeleteVideoFromAlbum(stream NewPhotos_DeleteVideoFromAlbumServer) error {
	log.Logger.UsingAccessLogFile().CInfoln("DeleteVideoFromAlbum")

	select {
	case <-stream.Context().Done():
		return nil
	default:
		for {
			msg, err := stream.Recv()
			if err != nil {
				break
			}
			s.DBInstanse.DeleteVideoFromAlbum(s.DBInstanse.GetUserID(msg.GetAccessToken(), msg.GetLoginToken()), msg.GetAlbum(), msg.GetVideo())
		}
		if err := stream.SendAndClose(&DeleteVideoFromAlbumResponse{Ok: true}); err != nil {
			log.Logger.UsingErrorLogFile().CFatalln("DeleteVideoFromAlbum", err)
		}
	}
	caching.RedisInstanse.Clean(caching.GET_VIDEOS_FROM_ALBUM)
	return nil
}

func (s *NewPhoto) CreateAlbum(ctx context.Context, r *CreateAlbumRequest) (*CreateAlbumResponse, error) {
	log.Logger.UsingAccessLogFile().CInfoln("CreateAlbum")

	s.DBInstanse.CreateAlbum(s.DBInstanse.GetUserID(r.GetAccessToken(), r.GetLoginToken()), r.GetName())
	caching.RedisInstanse.Clean(caching.GET_ALBUMS)
	return &CreateAlbumResponse{Ok: true}, nil
}

func (s *NewPhoto) DeleteAlbum(ctx context.Context, r *DeleteAlbumRequest) (*DeleteAlbumResponse, error) {
	log.Logger.UsingAccessLogFile().CInfoln("DeleteAlbum")

	s.DBInstanse.DeleteAlbum(s.DBInstanse.GetUserID(r.GetAccessToken(), r.GetLoginToken()), r.GetName())
	caching.RedisInstanse.Clean(caching.GET_ALBUMS)
	return &DeleteAlbumResponse{Ok: true}, nil
}

func (s *NewPhoto) GetAlbumsNum(ctx context.Context, r *GetAlbumsNumRequest) (*GetAlbumsNumResponse, error) {
	log.Logger.UsingAccessLogFile().CInfoln("GetAlbumInfo")

	n := s.DBInstanse.GetAlbumsNum(s.DBInstanse.GetUserID(r.GetAccessToken(), r.GetLoginToken()))
	return &GetAlbumsNumResponse{Ok: true, Num: n}, nil
}

func (s *NewPhoto) Ping(ctx context.Context, r *PingRequest) (*PingResponse, error) {
	return &PingResponse{Pong: true}, nil
}

func (s *NewPhoto) GetFullMediaByThumbnail(ctx context.Context, r *GetFullMediaByThumbnailRequest) (*GetFullMediaByThumbnailResponse, error) {
	log.Logger.UsingAccessLogFile().CInfoln("GetFullMediaByThumbnail")

	if cr, cached := caching.RedisInstanse.IsCached(s.DBInstanse.GetUserID(r.GetAccessToken(), r.GetLoginToken()), caching.GET_FULL_MEDIA_BY_THUMBNAIL); cached {
		definer := caching.Definer{Model: caching.GetFullMediaByThumbnail{}, Data: cr}
		result := definer.Define()
		converted, ok := result.(*caching.GetFullMediaByThumbnail)
		if !ok {
			log.Logger.UsingErrorLogFile().CFatalln("GetFullPhotoByThumbnail", caching.ErrConverting)
		}

		return &GetFullMediaByThumbnailResponse{Media: converted.Media}, nil
	}

	media := s.DBInstanse.GetFullMediaByThumbnail(s.DBInstanse.GetUserID(r.GetAccessToken(), r.GetLoginToken()), r.GetThumbnail(), &db.MediaSize{Width: 100, Height: 100}, db.MediaType(r.GetMediaType().Number()))

	model := caching.GetFullMediaByThumbnail{Media: media}
	conf := caching.DataConfigurator{Model: model}
	caching.RedisInstanse.Set(s.DBInstanse.GetUserID(r.GetAccessToken(), r.GetLoginToken()), caching.GET_FULL_MEDIA_BY_THUMBNAIL, string(conf.Configure()))

	return &GetFullMediaByThumbnailResponse{Media: media, Ok: true}, nil
}

func (s *NewPhoto) mustEmbedUnimplementedNewPhotosServer() {}

func NewNewPhoto() *NewPhoto {
	r := new(NewPhoto)
	r.DBInstanse = db.New()
	return r
}
