package meta

import (
	"context"
	"fmt"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqljson"
	"github.com/mangaweb4/mangaweb4-backend/ent"
	"github.com/mangaweb4/mangaweb4-backend/ent/meta"
	"github.com/mangaweb4/mangaweb4-backend/ent/tag"
	"github.com/mangaweb4/mangaweb4-backend/ent/user"
	"github.com/mangaweb4/mangaweb4-backend/grpc"
)

type QueryParams struct {
	SearchName  string
	SearchTag   string
	SortBy      grpc.SortField
	SortOrder   grpc.SortOrder
	Filter      grpc.Filter
	Page        int
	ItemPerPage int
}

func CreateQuery(ctx context.Context, client *ent.Client, u *ent.User, q QueryParams) (query *ent.MetaQuery, err error) {
	if q.SearchTag != "" {
		t, e := client.Tag.Query().Where(tag.Name(q.SearchTag)).Only(ctx)
		if e != nil {
			err = e
			return
		}

		query = t.QueryMeta()
	} else {
		query = client.Meta.Query()
	}

	query = query.Where(meta.Active(true), meta.Hidden(false))

	if q.SearchName != "" {
		query = query.Where(meta.NameContainsFold(q.SearchName))
	}

	switch q.Filter {
	case grpc.Filter_FILTER_FAVORITE_ITEMS:
		query = query.Where(
			meta.HasFavoriteOfUserWith(user.ID(u.ID)),
		)
	case grpc.Filter_FILTER_FAVORITE_TAGS:
		query = query.Where(
			meta.HasTagsWith(tag.HasFavoriteOfUserWith(user.ID(u.ID))),
		)
	}

	field := ""
	switch q.SortBy {
	case grpc.SortField_SORT_FIELD_NAME:
		field = meta.FieldName
	case grpc.SortField_SORT_FIELD_CREATION_TIME:
		field = meta.FieldCreateTime
	case grpc.SortField_SORT_FIELD_PAGECOUNT:
		field = meta.FieldFileIndices

	default:
		err = fmt.Errorf("invalid filter value: %v", q.SortBy)
	}

	switch q.SortOrder {
	case grpc.SortOrder_SORT_ORDER_ASCENDING:
		if q.SortBy == grpc.SortField_SORT_FIELD_PAGECOUNT {
			query = query.Order(sqljson.OrderLen(meta.FieldFileIndices)).Unique(false)
		} else {
			query = query.Order(ent.Asc(string(field)))
		}
	case grpc.SortOrder_SORT_ORDER_DESCENDING:
		if q.SortBy == grpc.SortField_SORT_FIELD_PAGECOUNT {
			query = query.Order(sqljson.OrderLenDesc(meta.FieldFileIndices)).Unique(false)
		} else {
			query = query.Order(ent.Desc(string(field)))
		}
	}

	if q.ItemPerPage > 0 {
		query = query.Limit(q.ItemPerPage).Offset(q.ItemPerPage * q.Page)
	}

	return
}

func ReadPage(ctx context.Context, client *ent.Client, u *ent.User, q QueryParams) (items []*ent.Meta, err error) {
	query, err := CreateQuery(ctx, client, u, q)
	if err != nil {
		return
	}

	return query.All(ctx)
}

func Count(ctx context.Context, client *ent.Client, u *ent.User, q QueryParams) (count int, err error) {
	query, err := CreateQuery(ctx, client, u, q)
	if err != nil {
		return
	}

	return query.Count(ctx)
}

func IsItemExist(ctx context.Context, client *ent.Client, name string) bool {
	count, err := client.Meta.Query().Where(meta.Name(name)).Count(ctx)
	if err != nil {
		return false
	}

	return count > 0
}

func Write(ctx context.Context, client *ent.Client, m *ent.Meta) error {
	return client.Meta.Create().
		SetName(m.Name).
		SetCreateTime(m.CreateTime).
		SetFavorite(m.Favorite).
		SetFileIndices(m.FileIndices).
		SetRead(m.Read).
		SetActive(m.Active).
		SetContainerType(m.ContainerType).
		SetThumbnailIndex(m.ThumbnailIndex).
		SetThumbnailX(m.ThumbnailX).
		SetThumbnailY(m.ThumbnailY).
		SetThumbnailWidth(m.ThumbnailWidth).
		SetThumbnailHeight(m.ThumbnailHeight).
		OnConflict(sql.ConflictColumns(meta.FieldName)).
		UpdateNewValues().Exec(ctx)
}

func Read(ctx context.Context, client *ent.Client, name string) (m *ent.Meta, err error) {
	return client.Meta.Query().Where(meta.Name(name)).Only(ctx)
}

func ReadAll(ctx context.Context, client *ent.Client) (items []*ent.Meta, err error) {
	return client.Meta.Query().Where(meta.Active(true)).All(ctx)
}
