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
	select {
	case <-stream.Context().Done():
		return nil
	default:
		if cr, cached := caching.RedisInstanse.IsCached(r.GetUserid(), caching.GET_PHOTOS); cached {
			definer := caching.NewDefiner()
			result, err := definer.Define(caching.GET_PHOTOS, cr)
			if err != nil {
				log.Logger.Fatalln(err)
			}
			converted, ok := result.([]caching.GetPhotosModel)
			if !ok {
				log.Logger.Fatalln(err)
			}
			for _, value := range converted {
				if err := stream.Send(&GetPhotosResponse{Photo: value.Photo, Thumbnail: value.Thumbnail, Extension: value.Extension, Size: value.Size, Tags: value.Tags, Ok: true}); err != nil {
					continue
				}
			}
			return nil
		}

		result := s.DBInstanse.GetPhotos(r.GetUserid())

		var model []caching.GetPhotosModel
		for _, value := range result {
			model = append(model, caching.GetPhotosModel{Photo: value[0].([]byte), Thumbnail: value[1].([]byte), Extension: value[2].(string), Size: value[3].(float64), Tags: value[4].(string)})
		}
		conf := caching.NewConfigurator()
		caching.RedisInstanse.Set(r.GetUserid(), caching.GET_PHOTOS, conf.Configure(model))
		for _, value := range result {
			if err := stream.Send(&GetPhotosResponse{Photo: value[0].([]byte), Thumbnail: value[1].([]byte), Extension: value[2].(string), Size: value[3].(float64), Tags: value[4].(string), Ok: true}); err != nil {
				continue
			}
		}
		return nil
	}
}

func (s *NewPhoto) GetVideos(r *GetVideosRequest, stream NewPhotos_GetVideosServer) error {
	select {
	case <-stream.Context().Done():
		return nil
	default:
		if cr, cached := caching.RedisInstanse.IsCached(r.GetUserid(), caching.GET_VIDEOS); cached {
			definer := caching.NewDefiner()
			result, err := definer.Define(caching.GET_VIDEOS, cr)
			if err != nil {
				log.Logger.Fatalln(err)
			}
			converted, ok := result.([]caching.GetVideosModel)
			if !ok {
				log.Logger.Fatalln(err)
			}
			for _, value := range converted {
				if err := stream.Send(&GetVideosResponse{Video: value.Video, Ok: true}); err != nil {
					continue
				}
			}
			return nil
		}

		result := s.DBInstanse.GetVideos(r.GetUserid())

		var model []caching.GetVideosModel
		for _, value := range result {
			model = append(model, caching.GetVideosModel{Video: value[0].([]byte)})
		}
		conf := caching.NewConfigurator()
		caching.RedisInstanse.Set(r.GetUserid(), caching.GET_VIDEOS, conf.Configure(model))
		for _, value := range result {
			if err := stream.Send(&GetVideosResponse{Video: value[0].([]byte), Ok: true}); err != nil {
				continue
			}
		}
		return nil
	}
}

func (s *NewPhoto) UploadPhoto(stream NewPhotos_UploadPhotoServer) error {
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
			s.DBInstanse.UploadPhoto(msg.GetUserid(), msg.GetPhoto(), msg.GetThumbnail(), msg.GetExtension(), msg.GetSize(), "")
		}
		if err := stream.SendAndClose(&UploadPhotoResponse{Ok: true}); err != nil {
			log.Logger.Fatalln(err)
		}
	}
	return nil
}

func (s *NewPhoto) UploadVideo(stream NewPhotos_UploadVideoServer) error {
	select {
	case <-stream.Context().Done():
		return nil
	default:
		for {
			msg, err := stream.Recv()
			if err != nil {
				break
			}
			s.DBInstanse.UploadVideo(msg.GetUserid(), msg.GetExtension(), msg.GetVideo(), msg.GetSize())
		}
		if err := stream.SendAndClose(&UploadVideoResponse{Ok: true}); err != nil {
			log.Logger.Fatalln(err)
		}
	}
	return nil
}

func (s *NewPhoto) GetUserinfo(cxt context.Context, r *GetUserinfoRequest) (*GetUserinfoResponse, error) {
	if cr, cached := caching.RedisInstanse.IsCached(r.GetUserid(), caching.GET_USER_INFO); cached {
		definer := caching.NewDefiner()
		result, err := definer.Define(caching.GET_USER_INFO, cr)
		if err != nil {
			log.Logger.Fatalln(err)
		}
		converted, ok := result.(caching.GetUserinfoModel)
		if !ok {
			log.Logger.Fatalln(err)
		}

		return &GetUserinfoResponse{Firstname: converted.Firstname, Secondname: converted.Secondname, Storage: converted.Storage, Ok: true}, nil
	}

	firstname, secondname, storage := s.DBInstanse.GetUserinfo(r.GetUserid())

	// Adds the newest info to redis ...
	model := caching.GetUserinfoModel{Firstname: firstname, Secondname: secondname, Storage: storage}
	conf := caching.NewConfigurator()
	caching.RedisInstanse.Set(r.GetUserid(), caching.GET_USER_INFO, conf.Configure(model))

	return &GetUserinfoResponse{Firstname: firstname, Secondname: secondname, Storage: storage, Ok: true}, nil
}

func (s *NewPhoto) GetUserAvatar(ctx context.Context, r *GetUserAvatarRequest) (*GetUserAvatarResponse, error) {
	avatar := s.DBInstanse.GetUserAvatar(r.GetUserid())
	return &GetUserAvatarResponse{Avatar: avatar, Ok: true}, nil
}

func (s *NewPhoto) SetUserAvatar(ctx context.Context, r *SetUserAvatarRequest) (*SetUserAvatarResponse, error) {
	s.DBInstanse.SetUserAvatar(r.GetUserid(), r.GetAvatar())
	return &SetUserAvatarResponse{Ok: true}, nil
}

func (s *NewPhoto) GetPhotosFromAlbum(r *GetPhotosFromAlbumRequest, stream NewPhotos_GetPhotosFromAlbumServer) error {
	select {
	case <-stream.Context().Done():
		return nil
	default:
		result := s.DBInstanse.GetPhotosFromAlbum(r.GetUserid(), r.GetName())
		var model []caching.GetPhotosFromAlbum
		for _, value := range result {
			model = append(model, caching.GetPhotosFromAlbum{Photo: value[0].([]byte), Thumbnail: value[1].([]byte), Extension: value[2].(string), Size: value[3].(float64), Album: value[4].(string)})
		}
		conf := caching.NewConfigurator()
		caching.RedisInstanse.Set(r.GetUserid(), caching.GET_PHOTOS_FROM_ALBUM, conf.Configure(model))
		for _, value := range result {
			if err := stream.Send(&GetPhotosFromAlbumResponse{Photo: value[0].([]byte), Thumbnail: value[1].([]byte), Extension: value[2].(string), Size: value[3].(float64), Album: value[4].(string), Ok: true}); err != nil {
				log.Logger.Fatalln(err)
				continue
			}
		}
		return nil
	}
}

func (s *NewPhoto) GetVideosFromAlbum(r *GetVideosFromAlbumRequest, stream NewPhotos_GetVideosFromAlbumServer) error {
	select {
	case <-stream.Context().Done():
		return nil
	default:
		result := s.DBInstanse.GetVideosFromAlbum(r.GetUserid(), r.GetName())

		var model []caching.GetVideosFromAlbum
		for _, value := range result {
			model = append(model, caching.GetVideosFromAlbum{Video: value[0].([]byte), Extension: value[1].(string)})
		}
		conf := caching.NewConfigurator()
		caching.RedisInstanse.Set(r.GetUserid(), caching.GET_VIDEOS_FROM_ALBUM, conf.Configure(model))
		for _, value := range result {
			if err := stream.Send(&GetVideosFromAlbumResponse{Video: value[0].([]byte), Extension: value[1].(string), Ok: true}); err != nil {
				continue
			}
		}
		return nil
	}
}

func (s *NewPhoto) UploadPhotoToAlbum(stream NewPhotos_UploadPhotoToAlbumServer) error {
	select {
	case <-stream.Context().Done():
		return nil
	default:
		for {
			recv, err := stream.Recv()
			if err != nil {
				break
			}
			s.DBInstanse.UploadPhotoToAlbum(recv.GetUserid(), recv.GetExtension(), recv.GetAlbum(), recv.GetSize(), recv.GetPhoto(), recv.GetThumbnail())
		}
		if err := stream.SendAndClose(&UploadPhotoToAlbumResponse{Ok: true}); err != nil {
			log.Logger.Fatalln(err)
		}
	}
	return nil
}

func (s *NewPhoto) UploadVideoToAlbum(stream NewPhotos_UploadVideoToAlbumServer) error {
	select {
	case <-stream.Context().Done():
		return nil
	default:
		for {
			recv, err := stream.Recv()
			if err != nil {
				break
			}
			s.DBInstanse.UploadVideoToAlbum(recv.GetUserid(), recv.GetAlbum(), recv.GetExtension(), recv.GetVideo(), recv.GetSize())
		}
		if err := stream.SendAndClose(&UploadVideoToAlbumResponse{Ok: true}); err != nil {
			log.Logger.Fatalln(err)
		}
	}
	return nil
}

func (s *NewPhoto) GetAlbums(r *GetAlbumsRequest, stream NewPhotos_GetAlbumsServer) error {
	select {
	case <-stream.Context().Done():
		return nil
	default:
		result := s.DBInstanse.GetAlbums(r.GetUserid())

		var model []caching.GetAlbumsModel
		for _, value := range result {
			model = append(model, caching.GetAlbumsModel{Name: value[0].(string), LatestPhoto: value[1].([]byte), LatestPhotoThumbnail: value[2].([]byte)})
		}

		conf := caching.NewConfigurator()
		caching.RedisInstanse.Set(r.GetUserid(), caching.GET_ALBUMS, conf.Configure(model))
		for _, value := range result {
			if err := stream.Send(&GetAlbumsResponse{Name: value[0].(string), LatestPhoto: value[1].([]byte), LatestPhotoThumbnail: value[2].([]byte), Ok: true}); err != nil {
				continue
			}
		}
		return nil
	}
}

func (s *NewPhoto) DeletePhotoFromAlbum(stream NewPhotos_DeletePhotoFromAlbumServer) error {
	select {
	case <-stream.Context().Done():
		return nil
	default:
		for {
			msg, err := stream.Recv()
			if err != nil {
				break
			}
			s.DBInstanse.DeletePhotoFromAlbum(msg.GetUserid(), msg.GetAlbum(), msg.GetPhoto())
		}
		if err := stream.SendAndClose(&DeletePhotoFromAlbumResponse{Ok: true}); err != nil {
			log.Logger.Fatalln(err)
		}
		return nil
	}
}

func (s *NewPhoto) DeleteVideoFromAlbum(stream NewPhotos_DeleteVideoFromAlbumServer) error {
	select {
	case <-stream.Context().Done():
		return nil
	default:
		for {
			msg, err := stream.Recv()
			if err != nil {
				break
			}
			s.DBInstanse.DeleteVideoFromAlbum(msg.GetUserid(), msg.GetAlbum(), msg.GetVideo())
		}
		if err := stream.SendAndClose(&DeleteVideoFromAlbumResponse{Ok: true}); err != nil {
			log.Logger.Fatalln(err)
		}
		return nil
	}
}

func (s *NewPhoto) CreateAlbum(ctx context.Context, r *CreateAlbumRequest) (*CreateAlbumResponse, error) {
	s.DBInstanse.CreateAlbum(r.GetUserid(), r.GetName())
	return &CreateAlbumResponse{Ok: true}, nil
}

func (s *NewPhoto) DeleteAlbum(ctx context.Context, r *DeleteAlbumRequest) (*DeleteAlbumResponse, error) {
	s.DBInstanse.DeleteAlbum(r.GetUserid(), r.GetName())
	return &DeleteAlbumResponse{Ok: true}, nil
}

func (s *NewPhoto) GetAlbumInfo(ctx context.Context, r *GetAlbumInfoRequest) (*GetAlbumInfoResponse, error) {
	n := s.DBInstanse.GetAlbumInfo(r.GetUserid(), r.GetAlbum())
	return &GetAlbumInfoResponse{Ok: true, MediaNum: n}, nil
}

func (s *NewPhoto) GetFullPhotoByThumbnail(ctx context.Context, r *GetFullPhotoByThumbnailRequest) (*GetFullPhotoByThumbnailResponse, error) {

	if cr, cached := caching.RedisInstanse.IsCached(r.GetUserid(), caching.GET_FULL_PHOTO_BY_THUMBNAIL); cached {
		definer := caching.NewDefiner()
		result, err := definer.Define(caching.GET_FULL_PHOTO_BY_THUMBNAIL, cr)
		if err != nil {
			log.Logger.Fatalln(err)
		}
		converted, ok := result.(caching.GetFullPhotoByThumbnail)
		if !ok {
			log.Logger.Fatalln(err)
		}

		return &GetFullPhotoByThumbnailResponse{Photo: converted.Photo}, nil
	}
	photo := s.DBInstanse.GetFullPhotoByThumbnail(r.GetUserid(), r.GetThumbnail())

	model := caching.GetFullPhotoByThumbnail{Photo: photo}
	conf := caching.NewConfigurator()
	caching.RedisInstanse.Set(r.GetUserid(), caching.GET_FULL_PHOTO_BY_THUMBNAIL, conf.Configure(model))

	return &GetFullPhotoByThumbnailResponse{Photo: photo, Ok: true}, nil
}

func (s *NewPhoto) mustEmbedUnimplementedNewPhotosServer() {}

func NewNewPhoto() *NewPhoto {
	r := new(NewPhoto)
	r.DBInstanse = db.New()
	return r
}
