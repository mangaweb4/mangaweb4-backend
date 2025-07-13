package tag

import (
	"context"

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

func CreateQuery(client *ent.Client, u *ent.User, q QueryParams) *ent.TagQuery {
	query := client.Tag.Query()
	if q.ItemPerPage > 0 {
		query = query.Limit(q.ItemPerPage).
			Offset(q.Page * q.ItemPerPage)
	}

	if q.Filter == grpc.Filter_FILTER_FAVORITE_ITEMS {
		query = query.Where(tag.HasFavoriteOfUserWith(user.ID(u.ID)))
	}
	if q.Search != "" {
		query = query.Where(tag.NameContainsFold(q.Search))
	}

	switch q.Sort {
	case grpc.SortField_SORT_FIELD_NAME:
		if q.Order == grpc.SortOrder_SORT_ORDER_ASCENDING {
			query = query.Order(tag.ByName(sql.OrderAsc()))
		} else {
			query = query.Order(tag.ByName(sql.OrderDesc()))
		}
	case grpc.SortField_SORT_FIELD_ITEMCOUNT:
		if q.Order == grpc.SortOrder_SORT_ORDER_ASCENDING {
			query = query.Order(tag.ByMetaCount(sql.OrderAsc()))
		} else {
			query = query.Order(tag.ByMetaCount(sql.OrderDesc()))
		}
	}

	return query
}

func ReadPage(ctx context.Context, client *ent.Client, u *ent.User, q QueryParams) (tags []*ent.Tag, err error) {
	query := CreateQuery(client, u, q)
	return query.All(ctx)
}

func Count(ctx context.Context, client *ent.Client, u *ent.User, q QueryParams) (count int, err error) {
	query := CreateQuery(client, u, q)

	return query.Count(ctx)
}

func Write(ctx context.Context, client *ent.Client, t *ent.Tag) error {
	return client.Tag.Create().
		SetName(t.Name).
		SetHidden(t.Hidden).
		SetFavorite(t.Favorite).
		OnConflict(sql.ConflictColumns(tag.FieldName)).
		UpdateNewValues().
		Exec(ctx)
}
