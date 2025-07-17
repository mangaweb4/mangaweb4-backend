package tag

import (
	"context"
	"fmt"

	"entgo.io/ent/dialect/sql"
	"github.com/mangaweb4/mangaweb4-backend/ent"
	"github.com/mangaweb4/mangaweb4-backend/ent/tag"
	"github.com/mangaweb4/mangaweb4-backend/ent/user"
	"github.com/mangaweb4/mangaweb4-backend/grpc"
)

type Filter string

func IsTagExist(ctx context.Context, client *ent.Client, name string) bool {
	count, err := client.Tag.Query().Where(tag.Name(name)).Count(ctx)
	if err != nil {
		return false
	}

	return count > 0
}

func Read(ctx context.Context, client *ent.Client, name string) (t *ent.Tag, err error) {
	return client.Tag.Query().Where(tag.Name(name)).First(ctx)
}

func ReadAll(ctx context.Context, client *ent.Client) (tags []*ent.Tag, err error) {
	return client.Tag.Query().Order(tag.ByName()).All(ctx)
}

type QueryParams struct {
	Filter      grpc.Filter
	Search      string
	Page        int
	ItemPerPage int
	Sort        grpc.SortField
	Order       grpc.SortOrder
}

func CreateQuery(client *ent.Client, u *ent.User, params QueryParams) (query *ent.TagQuery, err error) {
	query = client.Tag.Query()
	if params.ItemPerPage > 0 {
		query = query.Limit(params.ItemPerPage).
			Offset(params.Page * params.ItemPerPage)
	}

	switch params.Filter {
	case grpc.Filter_FILTER_UNKNOWN:
		break

	case grpc.Filter_FILTER_FAVORITE_TAGS:
		query = query.Where(tag.HasFavoriteOfUserWith(user.ID(u.ID)))

	default:
		query = nil
		err = fmt.Errorf("invalid filter value: %v", params.Filter)
		return
	}

	if params.Search != "" {
		query = query.Where(tag.NameContainsFold(params.Search))
	}

	switch params.Sort {
	case grpc.SortField_SORT_FIELD_NAME:
		if params.Order == grpc.SortOrder_SORT_ORDER_ASCENDING {
			query = query.Order(tag.ByName(sql.OrderAsc()))
		} else {
			query = query.Order(tag.ByName(sql.OrderDesc()))
		}
	case grpc.SortField_SORT_FIELD_ITEMCOUNT:
		if params.Order == grpc.SortOrder_SORT_ORDER_ASCENDING {
			query = query.Order(tag.ByMetaCount(sql.OrderAsc()))
		} else {
			query = query.Order(tag.ByMetaCount(sql.OrderDesc()))
		}
	}

	return
}

func ReadPage(ctx context.Context, client *ent.Client, u *ent.User, params QueryParams) (tags []*ent.Tag, err error) {
	query, err := CreateQuery(client, u, params)
	if err != nil {
		return
	}
	return query.All(ctx)
}

func Count(ctx context.Context, client *ent.Client, u *ent.User, params QueryParams) (count int, err error) {
	query, err := CreateQuery(client, u, params)
	if err != nil {
		return
	}

	return query.Count(ctx)
}

func Write(ctx context.Context, client *ent.Client, t *ent.Tag) error {
	return client.Tag.Create().
		SetName(t.Name).
		SetHidden(t.Hidden).
		OnConflict(sql.ConflictColumns(tag.FieldName)).
		UpdateNewValues().
		Exec(ctx)
}
