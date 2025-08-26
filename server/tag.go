package server

import (
	"context"

	"github.com/mangaweb4/mangaweb4-backend/database"
	"github.com/mangaweb4/mangaweb4-backend/ent"
	ent_tag "github.com/mangaweb4/mangaweb4-backend/ent/tag"
	"github.com/mangaweb4/mangaweb4-backend/ent/taguser"
	"github.com/mangaweb4/mangaweb4-backend/grpc"
	"github.com/mangaweb4/mangaweb4-backend/meta"
	"github.com/mangaweb4/mangaweb4-backend/tag"
	"github.com/mangaweb4/mangaweb4-backend/user"

	"github.com/rs/zerolog/log"
)

type TagServer struct {
	grpc.UnimplementedTagServer
}

func (s *TagServer) List(ctx context.Context, req *grpc.TagListRequest) (
	resp *grpc.TagListResponse, err error,
) {
	defer func() { log.Err(err).Interface("request", req).Msg("TagServer.List") }()

	client := database.CreateEntClient()
	defer func() { log.Err(client.Close()).Msg("database client close on TagServer.List") }()

	u, err := user.GetUser(ctx, client, req.User)
	if err != nil {
		return
	}

	allTags, err := tag.ReadPage(ctx, client, u,
		tag.QueryParams{
			Filter:      req.Filter,
			Search:      req.Search,
			Page:        int(req.Page),
			ItemPerPage: int(req.ItemPerPage),
			Sort:        req.Sort,
			Order:       req.Order,
		})

	if err != nil {
		return
	}

	total, err := tag.Count(ctx, client, u,
		tag.QueryParams{
			Filter:      req.Filter,
			Search:      req.Search,
			Page:        0,
			ItemPerPage: 0,
			Sort:        req.Sort,
			Order:       req.Order,
		})

	if err != nil {
		return
	}

	resp = &grpc.TagListResponse{
		TotalPage: (int32(total) / req.ItemPerPage) + 1,
	}

	resp.Items = make([]*grpc.TagListResponseItem, len(allTags))
	for i, t := range allTags {
		items, e := t.QueryMeta().All(ctx)
		if e != nil {
			err = e
			return
		}
		resp.Items[i] = &grpc.TagListResponseItem{
			Name:       t.Name,
			IsFavorite: u.QueryFavoriteTags().Where(ent_tag.ID(t.ID)).ExistX(ctx),
			PageCount:  int32(len(items)),
		}
	}

	return
}

func (s *TagServer) Thumbnail(ctx context.Context, req *grpc.TagThumbnailRequest) (
	resp *grpc.TagThumbnailResponse, err error,
) {
	defer func() { log.Err(err).Interface("request", req).Msg("TagServer.Thumbnail") }()

	client := database.CreateEntClient()
	defer func() { log.Err(client.Close()).Msg("database client close on TagServer.Thumbnail") }()

	t, err := tag.Read(ctx, client, req.Name)
	if err != nil {
		return
	}

	m, err := t.QueryMeta().First(ctx)
	if err != nil {
		return
	}

	thumbnail, err := meta.GetThumbnailBytes(m)
	if err != nil {
		return
	}

	resp = &grpc.TagThumbnailResponse{
		Data:        thumbnail,
		ContentType: "image/jpeg",
	}

	return
}

func (s *TagServer) SetFavorite(ctx context.Context, req *grpc.TagSetFavoriteRequest) (
	resp *grpc.TagSetFavoriteResponse, err error,
) {
	defer func() { log.Err(err).Interface("request", req).Msg("TagServer.SetFavorite") }()

	client := database.CreateEntClient()
	defer func() { log.Err(client.Close()).Msg("database client close on TagServer.SetFavorite") }()

	m, err := tag.Read(ctx, client, req.Tag)
	if err != nil {
		return
	}

	u, err := user.GetUser(ctx, client, req.User)
	if err != nil {
		return
	}

	userTag, err := client.TagUser.Query().Where(taguser.TagID(m.ID), taguser.UserID(u.ID)).First(ctx)
	if ent.IsNotFound(err) {
		userTag, err = client.TagUser.Create().SetTag(m).SetUser(u).Save(ctx)
	}

	userTag.IsFavorite = req.Favorite

	_, err = userTag.Update().Save(ctx)

	if err != nil {
		return
	}

	resp = &grpc.TagSetFavoriteResponse{
		Tag:      req.Tag,
		Favorite: req.Favorite,
	}

	return
}
