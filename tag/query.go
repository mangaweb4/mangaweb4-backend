package tag

import (
	"context"

	"entgo.io/ent/dialect/sql"
	"github.com/mangaweb4/mangaweb4-backend/ent"
	"github.com/mangaweb4/mangaweb4-backend/ent/tag"
	"github.com/mangaweb4/mangaweb4-backend/ent/user"
)

type SortField string
type SortOrder string
type Filter string

const (
	SortFieldName      = SortField("name")
	SortFieldPageCount = SortField("itemCount")

	SortOrderAscending  = SortOrder("ascending")
	SortOrderDescending = SortOrder("descending")
)

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
	FavoriteOnly bool
	Search       string
	Page         int
	ItemPerPage  int
	Sort         SortField
	Order        SortOrder
}

func CreateQuery(client *ent.Client, u *ent.User, q QueryParams) *ent.TagQuery {
	query := client.Tag.Query()
	if q.ItemPerPage > 0 {
		query = query.Limit(q.ItemPerPage).
			Offset(q.Page * q.ItemPerPage)
	}

	if q.FavoriteOnly {
		query = query.Where(tag.HasFavoriteOfUserWith(user.ID(u.ID)))
	}
	if q.Search != "" {
		query = query.Where(tag.NameContainsFold(q.Search))
	}

	switch q.Sort {
	case SortFieldName:
		if q.Order == SortOrderAscending {
			query = query.Order(tag.ByName(sql.OrderAsc()))
		} else {
			query = query.Order(tag.ByName(sql.OrderDesc()))
		}
	case SortFieldPageCount:
		if q.Order == SortOrderAscending {
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
