package server

import (
	"context"

	"entgo.io/ent/dialect/sql"
	"github.com/mangaweb4/mangaweb4-backend/database"
	"github.com/mangaweb4/mangaweb4-backend/ent/history"
	ent_meta "github.com/mangaweb4/mangaweb4-backend/ent/meta"
	"github.com/mangaweb4/mangaweb4-backend/grpc"
	"github.com/mangaweb4/mangaweb4-backend/user"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/rs/zerolog/log"
)

type HistoryServer struct{}

func (s *HistoryServer) List(ctx context.Context, req *grpc.HistoryListRequest) (resp *grpc.HistoryListResponse, err error) {
	client := database.CreateEntClient()
	defer client.Close()

	u, err := user.GetUser(ctx, client, req.User)
	if err != nil {

		return
	}

	histories, err := client.User.QueryHistories(u).
		Order(history.ByCreateTime(sql.OrderDesc())).
		Limit(int(req.ItemPerPage)).
		Offset(int(req.ItemPerPage * req.Page)).All(ctx)

	if err != nil {

		return
	}

	items := make([]*grpc.HistoryListResponseItem, len(histories))

	for i, h := range histories {
		m, e := h.QueryItem().Only(ctx)
		if e != nil {
			err = e
			return
		}

		items[i] = &grpc.HistoryListResponseItem{
			Id:         int32(m.ID),
			Name:       m.Name,
			IsFavorite: u.QueryFavoriteItems().Where(ent_meta.ID(m.ID)).ExistX(ctx),
			IsRead:     true,
			PageCount:  int32(len(m.FileIndices)),
			AccessTime: timestamppb.New(h.CreateTime),
		}

		tags, e := m.QueryTags().All(ctx)
		if e != nil {
			return
		}

		for _, t := range tags {
			if t.Favorite {
				items[i].HasFavoriteTag = true
				break
			}
		}
	}

	count, err := client.History.Query().Count(ctx)

	if err != nil {

		return
	}

	pageCount := int32(count) / req.ItemPerPage
	if int32(count)%req.ItemPerPage > 0 {
		pageCount++
	}

	if req.Page > pageCount || req.Page < 0 {
		req.Page = 0
	}

	log.Info().
		Interface("request", req).
		Msg("Browse")

	resp = &grpc.HistoryListResponse{
		Items:     items,
		TotalPage: pageCount,
	}

	return
}
