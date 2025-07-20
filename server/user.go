package server

import (
	"context"

	"github.com/mangaweb4/mangaweb4-backend/database"
	"github.com/mangaweb4/mangaweb4-backend/ent/meta"
	"github.com/mangaweb4/mangaweb4-backend/ent/tag"
	"github.com/mangaweb4/mangaweb4-backend/grpc"
	"github.com/mangaweb4/mangaweb4-backend/user"
	"github.com/rs/zerolog/log"
)

type UserServer struct {
	grpc.UnimplementedUserServer
}

func (s *UserServer) Info(ctx context.Context, req *grpc.UserInfoRequest) (resp *grpc.UserInfoResponse, err error) {
	defer func() { log.Err(err).Msg("UserServer.Info") }()

	client := database.CreateEntClient()

	defer func() { log.Err(client.Close()).Msg("database client close on MangaServer.List") }()

	u, err := user.GetUser(ctx, client, req.User)
	if err != nil {
		return
	}

	countReadManga, err := u.QueryProgress().Count(ctx)
	if err != nil {
		return
	}

	countFavoriteManga, err := u.QueryFavoriteItems().Where(meta.Active(true), meta.Hidden(false)).Count(ctx)
	if err != nil {
		return
	}

	countFavoriteTag, err := u.QueryFavoriteTags().Where(tag.Hidden(false)).Count(ctx)
	if err != nil {
		return
	}

	resp = &grpc.UserInfoResponse{
		UserId:            int32(u.ID),
		FavoriteItemCount: int32(countFavoriteManga),
		FavoriteTagCount:  int32(countFavoriteTag),
		ReadItemCount:     int32(countReadManga),
	}

	return
}
