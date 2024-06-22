package meta

import (
	"context"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqljson"
	"github.com/wutipong/mangaweb3-backend/ent"
	"github.com/wutipong/mangaweb3-backend/ent/meta"
	"github.com/wutipong/mangaweb3-backend/ent/tag"
)

type SortField string
type SortOrder string

const (
	SortFieldName       = SortField("name")
	SortFieldCreateTime = SortField("createTime")
	SortFieldPageCount  = SortField("pageCount")

	SortOrderAscending  = SortOrder("ascending")
	SortOrderDescending = SortOrder("descending")
)

type QueryParams struct {
	SearchName   string
	SearchTag    string
	SortBy       SortField
	SortOrder    SortOrder
	FavoriteOnly bool
	Page         int
	ItemPerPage  int
}

func CreateQuery(ctx context.Context, client *ent.Client, q QueryParams) (query *ent.MetaQuery, err error) {
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

	query = query.Where(meta.Active(true))

	if q.SearchName != "" {
		query = query.Where(meta.NameContainsFold(q.SearchName))
	}

	if q.FavoriteOnly {
		query = query.Where(meta.Favorite(true))
	}

	field := ""
	switch q.SortBy {
	case SortFieldName:
		field = meta.FieldName
	case SortFieldCreateTime:
		field = meta.FieldCreateTime
	case SortFieldPageCount:
		field = meta.FieldFileIndices
	}

	switch q.SortOrder {
	case SortOrderAscending:
		if q.SortBy == SortFieldPageCount {
			query = query.Order(sqljson.OrderLen(meta.FieldFileIndices)).Unique(false)
		} else {
			query = query.Order(ent.Asc(string(field)))
		}
	case SortOrderDescending:
		if q.SortBy == SortFieldPageCount {
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

func ReadPage(ctx context.Context, client *ent.Client, q QueryParams) (items []*ent.Meta, err error) {
	query, err := CreateQuery(ctx, client, q)
	if err != nil {
		return
	}

	return query.All(ctx)
}

func Count(ctx context.Context, client *ent.Client, q QueryParams) (count int, err error) {
	query, err := CreateQuery(ctx, client, q)
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
		SetThumbnail(m.Thumbnail).
		SetRead(m.Read).
		SetActive(m.Active).
		OnConflict(sql.ConflictColumns(meta.FieldName)).
		UpdateNewValues().Exec(ctx)
}

func Read(ctx context.Context, client *ent.Client, name string) (m *ent.Meta, err error) {
	return client.Meta.Query().Where(meta.Name(name)).Only(ctx)
}

func ReadAll(ctx context.Context, client *ent.Client) (items []*ent.Meta, err error) {
	return client.Meta.Query().Where(meta.Active(true)).All(ctx)
}
