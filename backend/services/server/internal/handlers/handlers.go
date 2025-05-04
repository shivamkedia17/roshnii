package handlers

import (
	"github.com/shivamkedia17/roshnii/shared/pkg/config"
	"github.com/shivamkedia17/roshnii/shared/pkg/db"
	"github.com/shivamkedia17/roshnii/shared/pkg/jwt"
	"github.com/shivamkedia17/roshnii/shared/pkg/storage"
)

type Handlers struct {
	OAuth GoogleOAuthService
	Img   ImageHandler
	Album AlbumHandler
	User  UserHandler
	// TODO Search
}

func InitHandlers(config *config.Config, db db.Store, storage storage.BlobStorage, jwt jwt.JWTService) Handlers {
	googleOAuthService := NewGoogleOAuthService(config, db, jwt)
	imageHandler := NewImageHandler(config, db, storage)
	albumHandler := NewAlbumHandler(config, db)
	userHandler := NewUserHandler(config, db)
	// TODO search

	return Handlers{
		OAuth: *googleOAuthService,
		Img:   *imageHandler,
		Album: *albumHandler,
		User:  *userHandler,
	}
}
