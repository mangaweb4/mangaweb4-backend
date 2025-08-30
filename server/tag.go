package server

import (
	"context"

	"github.com/mangaweb4/mangaweb4-backend/database"
	"github.com/mangaweb4/mangaweb4-backend/ent"
	"github.com/mangaweb4/mangaweb4-backend/ent/progress"
	ent_tag "github.com/mangaweb4/mangaweb4-backend/ent/tag"
	ent_user "github.com/mangaweb4/mangaweb4-backend/ent/user"
	"github.com/mangaweb4/mangaweb4-backend/grpc"
	"github.com/mangaweb4/mangaweb4-backend/meta"
	"github.com/mangaweb4/mangaweb4-backend/tag"
	"github.com/mangaweb4/mangaweb4-backend/user"

	"github.com/rs/zerolog/log"
)

type TagServer struct {
	grpc.UnimplementedTagServer
}

func (s *TagServer) List(
	ctx context.Context,
	req *grpc.TagListRequest,
) (resp *grpc.TagListResponse, err error) {
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
			Id:         int32(t.ID),
			Name:       t.Name,
			IsFavorite: u.QueryFavoriteTags().Where(ent_tag.ID(t.ID)).ExistX(ctx),
			PageCount:  int32(len(items)),
		}
	}

	return
}

func (s *TagServer) Detail(
	ctx context.Context,
	req *grpc.TagDetailRequest,
) (resp *grpc.TagDetailResponse, err error) {
	defer func() { log.Err(err).Interface("request", req).Msg("TagServer.Detail") }()

	client := database.CreateEntClient()
	defer func() { log.Err(client.Close()).Msg("database client close on TagServer.Detail") }()

	t, err := client.Tag.Get(ctx, int(req.Id))
	if err != nil {
		return
	}

	u, err := user.GetUser(ctx, client, req.User)
	if err != nil {
		return
	}

	items, err := tag.ReadMetaPage(ctx, client, t, u, tag.QueryMetaParams{
		SearchName:  req.Search,
		SortBy:      req.Sort,
		SortOrder:   req.Order,
		Filter:      req.Filter,
		Page:        int(req.Page),
		ItemPerPage: int(req.ItemPerPage),
	})
	if err != nil {
		return
	}

	count, err := tag.MetaCount(ctx, client, t, u, tag.QueryMetaParams{
		SearchName:  req.Search,
		SortBy:      req.Sort,
		SortOrder:   req.Order,
		Filter:      req.Filter,
		Page:        0,
		ItemPerPage: 0,
	})
	if err != nil {
		return
	}

	resp = &grpc.TagDetailResponse{
		Name:           t.Name,
		TotalItemCount: int32(count),
	}

	resp.TagFavorite, err = u.QueryFavoriteTags().Where(ent_tag.ID(t.ID)).Exist(ctx)
	if err != nil {
		return
	}

	for _, i := range items {
		p, e := i.QueryProgress().Where(progress.UserID(u.ID)).First(ctx)

		if e != nil {
			if !ent.IsNotFound(e) {
				err = e
				return
			}
			p = nil
		}

		item := &grpc.TagDetailResponseItem{
			Id:         int32(i.ID),
			Name:       i.Name,
			IsFavorite: i.QueryFavoriteOfUser().Where(ent_user.ID(u.ID)).ExistX(ctx),
			IsRead:     false,
			PageCount:  int32(len(i.FileIndices)),
			HasFavoriteTag: i.QueryTags().
				QueryFavoriteOfUser().
				Where(ent_user.ID(u.ID)).
				ExistX(ctx),
			CurrentPage: 0,
			MaxProgress: 0,
		}

		if p != nil {
			item.IsRead = true
			item.CurrentPage = int32(p.Page)
			item.MaxProgress = int32(p.Max)
		}
		resp.Items = append(resp.Items, item)
	}

	return
}

func (s *TagServer) Thumbnail(
	ctx context.Context,
	req *grpc.TagThumbnailRequest,
) (resp *grpc.TagThumbnailResponse, err error) {
	defer func() { log.Err(err).Interface("request", req).Msg("TagServer.Thumbnail") }()

	client := database.CreateEntClient()
	defer func() { log.Err(client.Close()).Msg("database client close on TagServer.Thumbnail") }()

	t, err := client.Tag.Get(ctx, int(req.Id))
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

func (s *TagServer) SetFavorite(
	ctx context.Context,
	req *grpc.TagSetFavoriteRequest,
) (resp *grpc.TagSetFavoriteResponse, err error) {
	defer func() { log.Err(err).Interface("request", req).Msg("TagServer.SetFavorite") }()

	client := database.CreateEntClient()
	defer func() { log.Err(client.Close()).Msg("database client close on TagServer.SetFavorite") }()

	t, err := client.Tag.Get(ctx, int(req.Id))
	if err != nil {
		return
	}

	u, err := user.GetUser(ctx, client, req.User)
	if err != nil {
		return
	}

	if req.Favorite {
		_, err = u.Update().AddFavoriteTags(t).Save(ctx)
	} else {
		_, err = u.Update().RemoveFavoriteTags(t).Save(ctx)
	}
	if err != nil {
		return
	}

	resp = &grpc.TagSetFavoriteResponse{
		Tag:      t.Name,
		Favorite: req.Favorite,
	}

	return
}
