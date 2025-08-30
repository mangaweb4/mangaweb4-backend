package tag

import (
	"context"
	"fmt"

	"entgo.io/ent/dialect/sql/sqljson"
	"github.com/mangaweb4/mangaweb4-backend/ent"
	"github.com/mangaweb4/mangaweb4-backend/ent/meta"
	"github.com/mangaweb4/mangaweb4-backend/ent/tag"
	"github.com/mangaweb4/mangaweb4-backend/ent/user"
	"github.com/mangaweb4/mangaweb4-backend/grpc"
)

type QueryMetaParams struct {
	SearchName  string
	SortBy      grpc.SortField
	SortOrder   grpc.SortOrder
	Filter      grpc.Filter
	Page        int
	ItemPerPage int
}

func CreateMetaQuery(
	ctx context.Context,
	client *ent.Client,
	t *ent.Tag,
	u *ent.User,
	q QueryMetaParams,
) (query *ent.MetaQuery, err error) {
	query = t.QueryMeta().Where(meta.Active(true), meta.Hidden(false))

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

func ReadMetaPage(
	ctx context.Context,
	client *ent.Client,
	t *ent.Tag,
	u *ent.User,
	q QueryMetaParams,
) (items []*ent.Meta, err error) {
	query, err := CreateMetaQuery(ctx, client, t, u, q)
	if err != nil {
		return
	}

	return query.All(ctx)
}

func MetaCount(
	ctx context.Context,
	client *ent.Client,
	t *ent.Tag,
	u *ent.User,
	q QueryMetaParams,
) (count int, err error) {
	query, err := CreateMetaQuery(ctx, client, t, u, q)
	if err != nil {
		return
	}

	return query.Count(ctx)
}
