package proto

import (
	"NewPhoto/caching"
	"NewPhoto/db"
	"NewPhoto/log"
	"context"
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
			definer := caching.NewDefiner()
			result, err := definer.Define(caching.GET_PHOTOS, cr)
			if err != nil {
				log.Logger.UsingErrorLogFile().CFatalln("GetPhotos", err)
			}
			converted, ok := result.([]caching.GetPhotosModel)
			if !ok {
				log.Logger.UsingErrorLogFile().CFatalln("GetPhotos", err)
			}
			for _, value := range converted {
				if err := stream.Send(&GetPhotosResponse{Photo: value.Photo, Thumbnail: value.Thumbnail, Extension: value.Extension, Size: value.Size, Tags: value.Tags, Ok: true}); err != nil {
					continue
				}
			}
			return nil
		}

		result := s.DBInstanse.GetPhotos(s.DBInstanse.GetUserID(r.GetAccessToken(), r.GetLoginToken()))

		var model []caching.GetPhotosModel
		for _, value := range result {
			model = append(model, caching.GetPhotosModel{Photo: value.Photo, Thumbnail: value.Thumbnail, Extension: value.Extension, Size: value.Size, Tags: value.Tags.String})
			if err := stream.Send(&GetPhotosResponse{Photo: value.Photo, Thumbnail: value.Thumbnail, Extension: value.Extension, Size: value.Size, Tags: value.Tags.String, Ok: true}); err != nil {
				continue
			}
		}
		conf := caching.NewConfigurator()
		caching.RedisInstanse.Set(s.DBInstanse.GetUserID(r.GetAccessToken(), r.GetLoginToken()), caching.GET_PHOTOS, conf.Configure(model))
		return nil
	}
}

func (s *NewPhoto) GetVideos(r *GetVideosRequest, stream NewPhotos_GetVideosServer) error {
	log.Logger.UsingAccessLogFile().CInfoln("GetVideos")
	select {
	case <-stream.Context().Done():
		return nil
	default:
		if cr, cached := caching.RedisInstanse.IsCached(s.DBInstanse.GetUserID(r.GetAccessToken(), r.GetLoginToken()), caching.GET_VIDEOS); cached {
			definer := caching.NewDefiner()
			result, err := definer.Define(caching.GET_VIDEOS, cr)
			if err != nil {
				log.Logger.UsingErrorLogFile().CFatalln("GetVideos", err)
			}
			converted, ok := result.([]caching.GetVideosModel)
			if !ok {
				log.Logger.UsingErrorLogFile().CFatalln("GetVideos", err)
			}
			for _, value := range converted {
				if err := stream.Send(&GetVideosResponse{Video: value.Video, Ok: true}); err != nil {
					continue
				}
			}
			return nil
		}

		result := s.DBInstanse.GetVideos(s.DBInstanse.GetUserID(r.GetAccessToken(), r.GetLoginToken()))

		var model []caching.GetVideosModel
		for _, value := range result {
			model = append(model, caching.GetVideosModel{Video: value.Video})
			if err := stream.Send(&GetVideosResponse{Video: value.Video, Ok: true}); err != nil {
				continue
			}
		}
		conf := caching.NewConfigurator()
		caching.RedisInstanse.Set(s.DBInstanse.GetUserID(r.GetAccessToken(), r.GetLoginToken()), caching.GET_VIDEOS, conf.Configure(model))

		return nil
	}
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
			s.DBInstanse.UploadVideo(s.DBInstanse.GetUserID(msg.GetAccessToken(), msg.GetLoginToken()), msg.GetExtension(), msg.GetVideo(), msg.GetSize())
		}
		if err := stream.SendAndClose(&UploadVideoResponse{Ok: true}); err != nil {
			log.Logger.UsingErrorLogFile().CFatalln("UploadVideo", err)
		}
	}
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
		definer := caching.NewDefiner()
		result, err := definer.Define(caching.GET_USER_INFO, cr)
		if err != nil {
			log.Logger.UsingErrorLogFile().CFatalln("GetUserinfo", err)
		}
		converted, ok := result.(caching.GetUserinfoModel)
		if !ok {
			log.Logger.UsingErrorLogFile().CFatalln("GetUserinfo", err)
		}

		return &GetUserinfoResponse{Firstname: converted.Firstname, Secondname: converted.Secondname, Storage: converted.Storage, Ok: true}, nil
	}

	firstname, secondname, storage := s.DBInstanse.GetUserinfo(s.DBInstanse.GetUserID(r.GetAccessToken(), r.GetLoginToken()))

	// Adds the newest info to redis ...
	model := caching.GetUserinfoModel{Firstname: firstname, Secondname: secondname, Storage: storage}
	conf := caching.NewConfigurator()
	caching.RedisInstanse.Set(s.DBInstanse.GetUserID(r.GetAccessToken(), r.GetLoginToken()), caching.GET_USER_INFO, conf.Configure(model))

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
		result := s.DBInstanse.GetPhotosFromAlbum(s.DBInstanse.GetUserID(r.GetAccessToken(), r.GetLoginToken()), r.GetName())
		var model []caching.GetPhotosFromAlbum
		for _, value := range result {
			model = append(model, caching.GetPhotosFromAlbum{Photo: value.Photo, Thumbnail: value.Thumbnail, Extension: value.Extension, Size: value.Size, Album: value.Album.String})
			if err := stream.Send(&GetPhotosFromAlbumResponse{Photo: value.Photo, Thumbnail: value.Thumbnail, Extension: value.Extension, Size: value.Size, Album: value.Album.String, Ok: true}); err != nil {
				log.Logger.UsingErrorLogFile().CFatalln("GetPhotosFromAlbum", err)
			}
		}
		conf := caching.NewConfigurator()
		caching.RedisInstanse.Set(s.DBInstanse.GetUserID(r.GetAccessToken(), r.GetLoginToken()), caching.GET_PHOTOS_FROM_ALBUM, conf.Configure(model))
		return nil
	}
}

func (s *NewPhoto) GetVideosFromAlbum(r *GetVideosFromAlbumRequest, stream NewPhotos_GetVideosFromAlbumServer) error {
	log.Logger.UsingAccessLogFile().CInfoln("GetVideosFromAlbum")

	select {
	case <-stream.Context().Done():
		return nil
	default:
		result := s.DBInstanse.GetVideosFromAlbum(s.DBInstanse.GetUserID(r.GetAccessToken(), r.GetLoginToken()), r.GetName())

		var model []caching.GetVideosFromAlbum
		for _, value := range result {
			model = append(model, caching.GetVideosFromAlbum{Video: value.Video, Extension: value.Extension})
			if err := stream.Send(&GetVideosFromAlbumResponse{Video: value.Video, Extension: value.Extension, Ok: true}); err != nil {
				continue
			}
		}
		conf := caching.NewConfigurator()
		caching.RedisInstanse.Set(s.DBInstanse.GetUserID(r.GetAccessToken(), r.GetLoginToken()), caching.GET_VIDEOS_FROM_ALBUM, conf.Configure(model))

		return nil
	}
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
			s.DBInstanse.UploadPhotoToAlbum(s.DBInstanse.GetUserID(recv.GetAccessToken(), recv.GetLoginToken()), recv.GetExtension(), recv.GetAlbum(), recv.GetSize(), recv.GetPhoto(), recv.GetThumbnail())
		}
		if err := stream.SendAndClose(&UploadPhotoToAlbumResponse{Ok: true}); err != nil {
			log.Logger.UsingErrorLogFile().CFatalln("UploadPhotoToAlbum", err)
		}
	}
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
			s.DBInstanse.UploadVideoToAlbum(s.DBInstanse.GetUserID(recv.GetAccessToken(), recv.GetLoginToken()), recv.GetExtension(), recv.GetAlbum(), recv.GetVideo(), recv.GetSize())
		}
		if err := stream.SendAndClose(&UploadVideoToAlbumResponse{Ok: true}); err != nil {
			log.Logger.UsingErrorLogFile().CFatalln("UploadVideoToAlbum", err)
		}
	}
	return nil
}

func (s *NewPhoto) GetAlbums(r *GetAlbumsRequest, stream NewPhotos_GetAlbumsServer) error {
	log.Logger.UsingAccessLogFile().CInfoln("GetAlbums")

	select {
	case <-stream.Context().Done():
		return nil
	default:
		result := s.DBInstanse.GetAlbums(s.DBInstanse.GetUserID(r.GetAccessToken(), r.GetLoginToken()))

		var model []caching.GetAlbumsModel

		for _, value := range result {
			model = append(model, caching.GetAlbumsModel{Name: value.Album, LatestPhoto: value.Photo})
			if err := stream.Send(&GetAlbumsResponse{Name: value.Album, LatestPhoto: value.Photo, Ok: true}); err != nil {
				continue
			}
		}
		conf := caching.NewConfigurator()
		caching.RedisInstanse.Set(s.DBInstanse.GetUserID(r.GetAccessToken(), r.GetLoginToken()), caching.GET_ALBUMS, conf.Configure(model))
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
		return nil
	}
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
		return nil
	}
}

func (s *NewPhoto) CreateAlbum(ctx context.Context, r *CreateAlbumRequest) (*CreateAlbumResponse, error) {
	log.Logger.UsingAccessLogFile().CInfoln("CreateAlbum")

	s.DBInstanse.CreateAlbum(s.DBInstanse.GetUserID(r.GetAccessToken(), r.GetLoginToken()), r.GetName())
	return &CreateAlbumResponse{Ok: true}, nil
}

func (s *NewPhoto) DeleteAlbum(ctx context.Context, r *DeleteAlbumRequest) (*DeleteAlbumResponse, error) {
	log.Logger.UsingAccessLogFile().CInfoln("DeleteAlbum")

	s.DBInstanse.DeleteAlbum(s.DBInstanse.GetUserID(r.GetAccessToken(), r.GetLoginToken()), r.GetName())
	return &DeleteAlbumResponse{Ok: true}, nil
}

func (s *NewPhoto) GetAlbumInfo(ctx context.Context, r *GetAlbumInfoRequest) (*GetAlbumInfoResponse, error) {
	log.Logger.UsingAccessLogFile().CInfoln("GetAlbumInfo")

	n := s.DBInstanse.GetAlbumInfo(s.DBInstanse.GetUserID(r.GetAccessToken(), r.GetLoginToken()), r.GetAlbum())
	return &GetAlbumInfoResponse{Ok: true, MediaNum: n}, nil
}

func (s *NewPhoto) Ping(ctx context.Context, r *PingRequest) (*PingResponse, error) {
	return &PingResponse{Pong: true}, nil
}

func (s *NewPhoto) GetFullPhotoByThumbnail(ctx context.Context, r *GetFullPhotoByThumbnailRequest) (*GetFullPhotoByThumbnailResponse, error) {
	log.Logger.UsingAccessLogFile().CInfoln("GetFullPhotoByThumbnail")

	if cr, cached := caching.RedisInstanse.IsCached(s.DBInstanse.GetUserID(r.GetAccessToken(), r.GetLoginToken()), caching.GET_FULL_PHOTO_BY_THUMBNAIL); cached {
		definer := caching.NewDefiner()
		result, err := definer.Define(caching.GET_FULL_PHOTO_BY_THUMBNAIL, cr)
		if err != nil {
			log.Logger.UsingErrorLogFile().CFatalln("GetFullPhotoByThumbnail", err)
		}
		converted, ok := result.(caching.GetFullPhotoByThumbnail)
		if !ok {
			log.Logger.UsingErrorLogFile().CFatalln("GetFullPhotoByThumbnail", err)
		}

		return &GetFullPhotoByThumbnailResponse{Photo: converted.Photo}, nil
	}
	photo := s.DBInstanse.GetFullPhotoByThumbnail(s.DBInstanse.GetUserID(r.GetAccessToken(), r.GetLoginToken()), r.GetThumbnail())

	model := caching.GetFullPhotoByThumbnail{Photo: photo}
	conf := caching.NewConfigurator()
	caching.RedisInstanse.Set(s.DBInstanse.GetUserID(r.GetAccessToken(), r.GetLoginToken()), caching.GET_FULL_PHOTO_BY_THUMBNAIL, conf.Configure(model))

	return &GetFullPhotoByThumbnailResponse{Photo: photo, Ok: true}, nil
}

func (s *NewPhoto) mustEmbedUnimplementedNewPhotosServer() {}

func NewNewPhoto() *NewPhoto {
	r := new(NewPhoto)
	r.DBInstanse = db.New()
	return r
}
