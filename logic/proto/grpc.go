package proto

import (
	"NewPhoto/caching"
	"NewPhoto/db"
	"context"
	"log"
	"os"

	"google.golang.org/grpc"
)

type Authentication struct{}

func (a *Authentication) LoginUser(ctx context.Context, r *UserLoginRequest) (*UserLoginResponse, error) {
	userid, err := db.DBInstanse.LoginUser(r.GetLogin(), r.GetPassword())
	if err != nil {
		return &UserLoginResponse{Userid: userid, Error: err.Error()}, nil
	}
	return &UserLoginResponse{Userid: userid, Error: "OK"}, nil
}

func (a *Authentication) RegisterUser(ctx context.Context, r *UserRegisterRequest) (*UserRegisterResponse, error) {
	db.DBInstanse.RegisterUser(r.GetLogin(), r.GetPassword(), r.GetFirstname(), r.GetSecondname())
	return &UserRegisterResponse{Error: "OK"}, nil
}

func (a *Authentication) mustEmbedUnimplementedAuthenticationServer() {}

type NewPhoto struct {
	Tag TagClient
}

func (s *NewPhoto) AllPhotos(r *AllPhotosRequest, stream NewPhotos_AllPhotosServer) error {
	select {
	case <-stream.Context().Done():
		return nil
	default:
		islogin := db.DBInstanse.IsLogin(r.GetUserid())
		if islogin {
			// Firstly checks whether similar request in the cache ...

			if cr, cached := caching.RedisInstanse.IsCached(r.GetUserid(), "AllPhotos"); cached {
				definer := caching.NewDefiner()
				result, err := definer.Define("AllPhotos", cr)
				if err != nil {
					log.Fatalln(err)
				}
				converted, ok := result.([]caching.AllPhotosModel)
				if !ok {
					log.Fatalln(err)
				}
				for _, value := range converted {
					err := stream.Send(&AllPhotosResponse{Photo: value.Photo, Thumbnail: value.Thumbnail, Extension: value.Extension, Size: value.Size, Tags: value.Tags, Error: "OK"})
					if err != nil {
						continue
					}
				}
				return nil
			}

			result := db.DBInstanse.AllPhotos(r.GetUserid())

			// Adds the newest info to redis ...
			var model []caching.AllPhotosModel
			for _, value := range result {
				model = append(model, caching.AllPhotosModel{Photo: value[0].([]byte), Thumbnail: value[1].([]byte), Extension: value[2].(string), Size: value[3].(float64), Tags: value[4].(string)})
			}
			conf := caching.NewConfigurator()
			caching.RedisInstanse.Set(r.GetUserid(), "AllPhotos", conf.Configure(model))
			for _, value := range result {
				err := stream.Send(&AllPhotosResponse{Photo: value[0].([]byte), Thumbnail: value[1].([]byte), Extension: value[2].(string), Size: value[3].(float64), Tags: value[4].(string), Error: "OK"})
				if err != nil {
					continue
				}
			}
			return nil
		}
		stream.Send(&AllPhotosResponse{Error: "user is not login"})
		return nil
	}
}

func (s *NewPhoto) UploadEqualPhoto(ctx context.Context, r *UploadEqualPhotoRequest) (*UploadEqualPhotoResponse, error) {

	tags, err := s.Tag.RecognizeObject(context.Background(), &RecognizeObjectRequest{Photo: r.GetPhoto()})
	if err != nil{
		log.Fatalln(err)
	}

	db.DBInstanse.UploadEqualPhoto(r.GetUserid(), r.GetPhoto(), r.GetThumbnail(), r.GetExtension(), r.GetSize(), tags.GetTags())
	return &UploadEqualPhotoResponse{Error: "OK"}, nil
}

func (s *NewPhoto) GetUserinfo(cxt context.Context, r *GetUserinfoRequest) (*GetUserinfoResponse, error) {

	islogin := db.DBInstanse.IsLogin(r.GetUserid())
	if islogin {

		// Firstly checks whether similar request in the cache ...

		if cr, cached := caching.RedisInstanse.IsCached(r.GetUserid(), "GetUserinfo"); cached {
			definer := caching.NewDefiner()
			result, err := definer.Define("GetUserinfo", cr)
			if err != nil {
				log.Fatalln(err)
			}
			converted, ok := result.(caching.GetUserinfoModel)
			if !ok {
				log.Fatalln(err)
			}

			return &GetUserinfoResponse{Firstname: converted.Firstname, Secondname: converted.Secondname, Storage: converted.Storage, Error: "OK"}, nil
		}

		firstname, secondname, storage := db.DBInstanse.GetUserinfo(r.GetUserid())

		// Adds the newest info to redis ...
		model := caching.GetUserinfoModel{Firstname: firstname, Secondname: secondname, Storage: storage}
		conf := caching.NewConfigurator()
		caching.RedisInstanse.Set(r.GetUserid(), "GetUserinfo", conf.Configure(model))

		return &GetUserinfoResponse{Firstname: firstname, Secondname: secondname, Storage: storage, Error: "OK"}, nil
	}
	return &GetUserinfoResponse{Error: "user is not login"}, nil
}

func (s *NewPhoto) AllPhotosAlbum(r *AllPhotosAlbumRequest, stream NewPhotos_AllPhotosAlbumServer) error {

	select {
	case <-stream.Context().Done():
		return nil
	default:
		islogin := db.DBInstanse.IsLogin(r.GetUserid())
		if islogin {
			// Firstly checks whether similar request in the cache ...

			if cr, cached := caching.RedisInstanse.IsCached(r.GetUserid(), "AllPhotosAlbum"); cached {

				definer := caching.NewDefiner()
				result, err := definer.Define("AllPhotosAlbum", cr)
				if err != nil {
					log.Fatalln(err)
				}
				converted, ok := result.([]caching.AllPhotosAlbumModel)
				if !ok {
					log.Fatalln("Is can't be converted")
				}
				for _, value := range converted {
					err := stream.Send(&AllPhotosAlbumResponse{Photo: value.Photo, Thumbnail: value.Thumbnail, Extension: value.Extension, Size: value.Size, Error: "OK", Album: r.GetName()})
					if err != nil {
						continue
					}
				}
				return nil
			}

			result := db.DBInstanse.AlbumPhotos(r.GetUserid(), r.GetName())

			// Adds the newest info to redis ...
			var model []caching.AllPhotosAlbumModel
			for _, value := range result {
				model = append(model, caching.AllPhotosAlbumModel{Photo: value[0].([]byte), Thumbnail: value[1].([]byte), Extension: value[2].(string), Size: value[3].(float64), Album: value[4].(string)})
			}
			conf := caching.NewConfigurator()
			caching.RedisInstanse.Set(r.GetUserid(), "AllPhotosAlbum", conf.Configure(model))
			for _, value := range result {
				err := stream.Send(&AllPhotosAlbumResponse{Photo: value[0].([]byte), Thumbnail: value[1].([]byte), Extension: value[2].(string), Size: value[3].(float64), Album: value[4].(string), Error: "OK"})
				if err != nil {
					continue
				}
			}
			return nil
		}
		stream.Send(&AllPhotosAlbumResponse{Error: "user is not login"})
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
				return nil
			}
			islogin := db.DBInstanse.IsLogin(recv.GetUserid())
			if islogin {
				db.DBInstanse.UploadEqualPhotoToAlbum(recv.GetUserid(), recv.GetExtension(), recv.GetAlbum(), recv.GetSize(), recv.GetPhoto(), recv.GetThumbnail())
				continue
			}
			return nil
		}
	}
}

func (s *NewPhoto) GetAllAlbums(r *GetAllAlbumsRequest, stream NewPhotos_GetAllAlbumsServer) error {

	select {
	case <-stream.Context().Done():
		return nil
	default:
		islogin := db.DBInstanse.IsLogin(r.GetUserid())
		if islogin {
			if cr, cached := caching.RedisInstanse.IsCached(r.GetUserid(), "GetAllAlbums"); cached {

				definer := caching.NewDefiner()
				result, err := definer.Define("GetAllAlbums", cr)
				if err != nil {

					log.Fatalln(err)
				}
				converted, ok := result.([]caching.GetAllAlbumsModel)
				if !ok {
					log.Fatalln(err)
				}
				for _, value := range converted {
					err := stream.Send(&GetAllAlbumsResponse{Name: value.Name, LatestPhoto: value.LatestPhoto, LatestPhotoThumbnail: value.LatestPhotoThumbnail})
					if err != nil {
						continue
					}
				}
				return nil
			}

			result := db.DBInstanse.GetAllAlbums(r.GetUserid())

			// Adds the newest info to redis ...
			var model []caching.GetAllAlbumsModel
			for _, value := range result {
				model = append(model, caching.GetAllAlbumsModel{Name: value[0].(string), LatestPhoto: value[1].([]byte), LatestPhotoThumbnail: value[2].([]byte)})
			}

			conf := caching.NewConfigurator()
			caching.RedisInstanse.Set(r.GetUserid(), "GetAllAlbums", conf.Configure(model))
			for _, value := range result {

				err := stream.Send(&GetAllAlbumsResponse{Name: value[0].(string), LatestPhoto: value[1].([]byte), LatestPhotoThumbnail: value[2].([]byte), Error: "OK"})
				if err != nil {

					continue
				}
			}

			return nil
		}
		stream.Send(&GetAllAlbumsResponse{Error: "user is not login"})
		return nil
	}
}

func (s *NewPhoto) GetFullPhotoByThumbnail(ctx context.Context, r *GetFullPhotoByThumbnailRequest) (*GetFullPhotoByThumbnailResponse, error) {
	islogin := db.DBInstanse.IsLogin(r.GetUserid())
	if islogin {
		if cr, cached := caching.RedisInstanse.IsCached(r.GetUserid(), "GetFullPhotoByThumbnail"); cached {
			definer := caching.NewDefiner()
			result, err := definer.Define("GetFullPhotoByThumbnail", cr)
			if err != nil {
				log.Fatalln(err)
			}
			converted, ok := result.(caching.GetFullPhotoByThumbnail)
			if !ok {
				log.Fatalln(err)
			}

			return &GetFullPhotoByThumbnailResponse{Photo: converted.Photo}, nil
		}
		photo := db.DBInstanse.GetFullPhotoByThumbnail(r.GetUserid(), r.GetThumbnail())

		model := caching.GetFullPhotoByThumbnail{Photo: photo}
		conf := caching.NewConfigurator()
		caching.RedisInstanse.Set(r.GetUserid(), "GetFullPhotoByThumbnail", conf.Configure(model))

		return &GetFullPhotoByThumbnailResponse{Photo: photo, Error: "OK"}, nil
	}
	return &GetFullPhotoByThumbnailResponse{Error: "User is not loggedin"}, nil
}

func (s *NewPhoto) CreateAlbum(ctx context.Context, r *CreateAlbumRequest) (*CreateAlbumResponse, error) {

	islogin := db.DBInstanse.IsLogin(r.GetUserid())
	if islogin {
		db.DBInstanse.CreateAlbum(r.GetUserid(), r.GetName())
		return &CreateAlbumResponse{Error: "OK"}, nil
	}
	return &CreateAlbumResponse{Error: "user is not login"}, nil
}

func (s *NewPhoto) DeleteAlbum(ctx context.Context, r *DeleteAlbumRequest) (*DeleteAlbumResponse, error) {

	islogin := db.DBInstanse.IsLogin(r.GetUserid())
	if islogin {

		db.DBInstanse.DeleteAlbum(r.GetUserid(), r.GetName())
		return &DeleteAlbumResponse{Error: "OK"}, nil
	}
	return &DeleteAlbumResponse{Error: "user is not login"}, nil
}

func (s *NewPhoto) mustEmbedUnimplementedNewPhotosServer() {}

func (s *NewPhoto) InitTagClient(){
	s.Tag = NewTag()
}

func NewNewPhoto() *NewPhoto {
	// Creates NewPhoto server instanse ...

	return new(NewPhoto)
}

func NewAuthentication() *Authentication {
	// Creates Authentication server instanse ...

	return new(Authentication)
}

func NewTag() TagClient {
	// Creates Tag client instanse ...

	opts := []grpc.DialOption{
		grpc.WithInsecure(),
		grpc.WithDefaultCallOptions(
			grpc.MaxCallRecvMsgSize(50 * 10e6),
			grpc.MaxCallSendMsgSize(50 * 10e6),
		),
	}

	aiAddr, ok := os.LookupEnv("aiAddr")
	if !ok{
		log.Fatalln("aiAddr is not written in credentials.sh file")
	}

	client, err := grpc.Dial(aiAddr, opts...)
	if err != nil {
		log.Fatalln(err.Error())
	}
	return NewTagClient(client)
}
