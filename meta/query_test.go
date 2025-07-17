package meta

import (
	"context"
	"database/sql"
	"testing"
	"time"

	dialect_sql "entgo.io/ent/dialect/sql"
	"github.com/mangaweb4/mangaweb4-backend/ent"
	"github.com/mangaweb4/mangaweb4-backend/ent/enttest"
	"github.com/mangaweb4/mangaweb4-backend/grpc"
	"github.com/stretchr/testify/suite"
	_ "modernc.org/sqlite"
)

type QueryTestSuite struct {
	suite.Suite
}

func TestQueryTestSuite(t *testing.T) {
	// suite.Run(t, new(QueryTestSuite))
}

func (s *QueryTestSuite) TestReadPage() {
	db, err := sql.Open("sqlite", "file:ent?mode=memory&_fk=1&_pragma=foreign_keys(1)")
	s.Assert().Nil(err)
	s.Assert().NotNil(db)
	defer db.Close()

	client := enttest.NewClient(s.T(), enttest.WithOptions(ent.Driver(dialect_sql.OpenDB("sqlite3", db))))
	defer client.Close()

	var u *ent.User

	client.Meta.Create().SetName("[some artist]manga 1 here.zip").Save(context.Background())
	client.Meta.Create().SetName("[some artist]manga 2 here.zip").Save(context.Background())
	client.Meta.Create().SetName("[some artist]manga 3 here.zip").Save(context.Background())
	client.Meta.Create().SetName("[some artist]manga 4 here.zip").Save(context.Background())
	client.Meta.Create().SetName("[some artist]manga 5 here.zip").SetActive(false).Save(context.Background())

	tags, err := ReadPage(context.Background(), client, u, QueryParams{
		SortBy:      grpc.SortField_SORT_FIELD_NAME,
		SortOrder:   grpc.SortOrder_SORT_ORDER_ASCENDING,
		Page:        0,
		ItemPerPage: 30,
	})
	s.Assert().Nil(err)

	s.Assert().Equal(4, len(tags))

	s.Assert().Equal("[some artist]manga 1 here.zip", tags[0].Name)
	s.Assert().Equal("[some artist]manga 2 here.zip", tags[1].Name)
	s.Assert().Equal("[some artist]manga 3 here.zip", tags[2].Name)
	s.Assert().Equal("[some artist]manga 4 here.zip", tags[3].Name)
}

func (s *QueryTestSuite) TestReadPageFilterFavoriteItems() {
	db, err := sql.Open("sqlite", "file:ent?mode=memory&_fk=1&_pragma=foreign_keys(1)")
	s.Assert().Nil(err)
	s.Assert().NotNil(db)
	defer db.Close()

	client := enttest.NewClient(s.T(), enttest.WithOptions(ent.Driver(dialect_sql.OpenDB("sqlite3", db))))
	defer client.Close()

	client.Meta.Create().SetName("[some artist]manga 1 here.zip").SetFavorite(true).Save(context.Background())
	client.Meta.Create().SetName("[some artist]manga 2 here.zip").SetFavorite(true).Save(context.Background())
	client.Meta.Create().SetName("[some artist]manga 3 here.zip").Save(context.Background())
	client.Meta.Create().SetName("[some artist]manga 4 here.zip").Save(context.Background())

	var u *ent.User
	tags, err := ReadPage(context.Background(), client, u, QueryParams{
		Filter:      grpc.Filter_FILTER_FAVORITE_ITEMS,
		SortBy:      grpc.SortField_SORT_FIELD_NAME,
		SortOrder:   grpc.SortOrder_SORT_ORDER_ASCENDING,
		Page:        0,
		ItemPerPage: 30,
	})
	s.Assert().Nil(err)

	s.Assert().Equal(2, len(tags))

	s.Assert().Equal("[some artist]manga 1 here.zip", tags[0].Name)
	s.Assert().Equal("[some artist]manga 2 here.zip", tags[1].Name)
}

func (s *QueryTestSuite) TestReadPageSortByCreateTimeDesc() {
	db, err := sql.Open("sqlite", "file:ent?mode=memory&_fk=1&_pragma=foreign_keys(1)")
	s.Assert().Nil(err)
	s.Assert().NotNil(db)
	defer db.Close()

	client := enttest.NewClient(s.T(), enttest.WithOptions(ent.Driver(dialect_sql.OpenDB("sqlite3", db))))
	defer client.Close()

	client.Meta.Create().SetName("[some artist]manga 1 here.zip").SetCreateTime(time.UnixMilli(3000)).Save(context.Background())
	client.Meta.Create().SetName("[some artist]manga 2 here.zip").SetCreateTime(time.UnixMilli(4000)).Save(context.Background())
	client.Meta.Create().SetName("[some artist]manga 3 here.zip").SetCreateTime(time.UnixMilli(5000)).SetActive(false).Save(context.Background())
	client.Meta.Create().SetName("[some artist]manga 4 here.zip").SetCreateTime(time.UnixMilli(6000)).SetActive(false).Save(context.Background())

	var u *ent.User
	tags, err := ReadPage(context.Background(), client, u, QueryParams{
		Filter:      grpc.Filter_FILTER_UNKNOWN,
		SortBy:      grpc.SortField_SORT_FIELD_CREATION_TIME,
		SortOrder:   grpc.SortOrder_SORT_ORDER_DESCENDING,
		Page:        0,
		ItemPerPage: 30,
	})
	s.Assert().Nil(err)

	s.Assert().Equal(2, len(tags))

	s.Assert().Equal("[some artist]manga 2 here.zip", tags[0].Name)
	s.Assert().Equal("[some artist]manga 1 here.zip", tags[1].Name)
}

func (s *QueryTestSuite) TestReadPageSortByCreateTimeAsc() {
	db, err := sql.Open("sqlite", "file:ent?mode=memory&_fk=1&_pragma=foreign_keys(1)")
	s.Assert().Nil(err)
	s.Assert().NotNil(db)
	defer db.Close()

	client := enttest.NewClient(s.T(), enttest.WithOptions(ent.Driver(dialect_sql.OpenDB("sqlite3", db))))
	defer client.Close()

	client.Meta.Create().SetName("[some artist]manga 1 here.zip").SetCreateTime(time.UnixMilli(3000)).Save(context.Background())
	client.Meta.Create().SetName("[some artist]manga 2 here.zip").SetCreateTime(time.UnixMilli(4000)).Save(context.Background())
	client.Meta.Create().SetName("[some artist]manga 3 here.zip").SetCreateTime(time.UnixMilli(5000)).SetActive(false).Save(context.Background())
	client.Meta.Create().SetName("[some artist]manga 4 here.zip").SetCreateTime(time.UnixMilli(6000)).SetActive(false).Save(context.Background())

	var u *ent.User
	tags, err := ReadPage(context.Background(), client, u, QueryParams{
		Filter:      grpc.Filter_FILTER_UNKNOWN,
		SortBy:      grpc.SortField_SORT_FIELD_CREATION_TIME,
		SortOrder:   grpc.SortOrder_SORT_ORDER_ASCENDING,
		Page:        0,
		ItemPerPage: 30,
	})
	s.Assert().Nil(err)

	s.Assert().Equal(2, len(tags))

	s.Assert().Equal("[some artist]manga 1 here.zip", tags[0].Name)
	s.Assert().Equal("[some artist]manga 2 here.zip", tags[1].Name)
}

func (s *QueryTestSuite) TestReadPageFavoriteItemsSSortByCreateTimeDesc() {
	db, err := sql.Open("sqlite", "file:ent?mode=memory&_fk=1&_pragma=foreign_keys(1)")
	s.Assert().Nil(err)
	s.Assert().NotNil(db)
	defer db.Close()

	client := enttest.NewClient(s.T(), enttest.WithOptions(ent.Driver(dialect_sql.OpenDB("sqlite3", db))))
	defer client.Close()

	client.Meta.Create().SetName("[some artist]manga 1 here.zip").SetCreateTime(time.UnixMilli(3000)).SetFavorite(true).Save(context.Background())
	client.Meta.Create().SetName("[some artist]manga 2 here.zip").SetCreateTime(time.UnixMilli(4000)).SetFavorite(true).Save(context.Background())
	client.Meta.Create().SetName("[some artist]manga 3 here.zip").SetCreateTime(time.UnixMilli(5000)).Save(context.Background())
	client.Meta.Create().SetName("[some artist]manga 4 here.zip").SetCreateTime(time.UnixMilli(6000)).Save(context.Background())

	var u *ent.User
	tags, err := ReadPage(context.Background(), client, u, QueryParams{
		Filter:      grpc.Filter_FILTER_FAVORITE_ITEMS,
		SortBy:      grpc.SortField_SORT_FIELD_CREATION_TIME,
		SortOrder:   grpc.SortOrder_SORT_ORDER_DESCENDING,
		Page:        0,
		ItemPerPage: 30,
	})
	s.Assert().Nil(err)

	s.Assert().Equal(2, len(tags))

	s.Assert().Equal("[some artist]manga 2 here.zip", tags[0].Name)
	s.Assert().Equal("[some artist]manga 1 here.zip", tags[1].Name)
}

func (s *QueryTestSuite) TestReadPageFilterFavoriteItemsortByCreateTimeAsc() {
	db, err := sql.Open("sqlite", "file:ent?mode=memory&_fk=1&_pragma=foreign_keys(1)")
	s.Assert().Nil(err)
	s.Assert().NotNil(db)
	defer db.Close()

	client := enttest.NewClient(s.T(), enttest.WithOptions(ent.Driver(dialect_sql.OpenDB("sqlite3", db))))
	defer client.Close()

	client.Meta.Create().SetName("[some artist]manga 1 here.zip").SetCreateTime(time.UnixMilli(3000)).SetFavorite(true).Save(context.Background())
	client.Meta.Create().SetName("[some artist]manga 2 here.zip").SetCreateTime(time.UnixMilli(4000)).SetFavorite(true).Save(context.Background())
	client.Meta.Create().SetName("[some artist]manga 3 here.zip").SetCreateTime(time.UnixMilli(5000)).Save(context.Background())
	client.Meta.Create().SetName("[some artist]manga 4 here.zip").SetCreateTime(time.UnixMilli(6000)).Save(context.Background())

	var u *ent.User

	tags, err := ReadPage(context.Background(), client, u, QueryParams{
		Filter:      grpc.Filter_FILTER_FAVORITE_ITEMS,
		SortBy:      grpc.SortField_SORT_FIELD_CREATION_TIME,
		SortOrder:   grpc.SortOrder_SORT_ORDER_ASCENDING,
		Page:        0,
		ItemPerPage: 30,
	})
	s.Assert().Nil(err)

	s.Assert().Equal(2, len(tags))

	s.Assert().Equal("[some artist]manga 1 here.zip", tags[0].Name)
	s.Assert().Equal("[some artist]manga 2 here.zip", tags[1].Name)
}

func (s *QueryTestSuite) TestReadPageSearchNameFilterFavoriteItemsortByCreateTimeDesc() {
	db, err := sql.Open("sqlite", "file:ent?mode=memory&_fk=1&_pragma=foreign_keys(1)")
	s.Assert().Nil(err)
	s.Assert().NotNil(db)
	defer db.Close()

	client := enttest.NewClient(s.T(), enttest.WithOptions(ent.Driver(dialect_sql.OpenDB("sqlite3", db))))
	defer client.Close()

	client.Meta.Create().SetName("[some artist]manga 1 here.zip").SetCreateTime(time.UnixMilli(3000)).SetFavorite(true).Save(context.Background())
	client.Meta.Create().SetName("[some artist]manga 2 here.zip").SetCreateTime(time.UnixMilli(4000)).SetFavorite(true).Save(context.Background())
	client.Meta.Create().SetName("[some artist]manga 3.zip").SetCreateTime(time.UnixMilli(5000)).SetFavorite(true).Save(context.Background())
	client.Meta.Create().SetName("[some artist]manga 4.zip").SetCreateTime(time.UnixMilli(6000)).SetFavorite(true).Save(context.Background())

	var u *ent.User
	tags, err := ReadPage(context.Background(), client, u, QueryParams{
		SearchName:  "here",
		SortBy:      grpc.SortField_SORT_FIELD_CREATION_TIME,
		SortOrder:   grpc.SortOrder_SORT_ORDER_DESCENDING,
		Filter:      grpc.Filter_FILTER_FAVORITE_ITEMS,
		Page:        0,
		ItemPerPage: 30,
	})
	s.Assert().Nil(err)

	s.Assert().Equal(2, len(tags))

	s.Assert().Equal("[some artist]manga 2 here.zip", tags[0].Name)
	s.Assert().Equal("[some artist]manga 1 here.zip", tags[1].Name)
}

func (s *QueryTestSuite) TestReadPageSearchNameFilterFavoriteItemsortByCreateTimeAsc() {
	db, err := sql.Open("sqlite", "file:ent?mode=memory&_fk=1&_pragma=foreign_keys(1)")
	s.Assert().Nil(err)
	s.Assert().NotNil(db)
	defer db.Close()

	client := enttest.NewClient(s.T(), enttest.WithOptions(ent.Driver(dialect_sql.OpenDB("sqlite3", db))))
	defer client.Close()

	client.Meta.Create().SetName("[some artist]manga 1 here.zip").SetCreateTime(time.UnixMilli(3000)).SetFavorite(true).Save(context.Background())
	client.Meta.Create().SetName("[some artist]manga 2 here.zip").SetCreateTime(time.UnixMilli(4000)).SetFavorite(true).Save(context.Background())
	client.Meta.Create().SetName("[some artist]manga 3.zip").SetCreateTime(time.UnixMilli(5000)).SetFavorite(true).Save(context.Background())
	client.Meta.Create().SetName("[some artist]manga 4.zip").SetCreateTime(time.UnixMilli(6000)).SetFavorite(true).Save(context.Background())

	var u *ent.User
	tags, err := ReadPage(context.Background(), client, u, QueryParams{
		SearchName:  "here",
		SortBy:      grpc.SortField_SORT_FIELD_CREATION_TIME,
		SortOrder:   grpc.SortOrder_SORT_ORDER_ASCENDING,
		Filter:      grpc.Filter_FILTER_FAVORITE_ITEMS,
		Page:        0,
		ItemPerPage: 30,
	})
	s.Assert().Nil(err)

	s.Assert().Equal(2, len(tags))

	s.Assert().Equal("[some artist]manga 1 here.zip", tags[0].Name)
	s.Assert().Equal("[some artist]manga 2 here.zip", tags[1].Name)
}

func (s *QueryTestSuite) TestReadPageSearchTagFilterFavoriteItemsortByCreateTimeAsc() {
	db, err := sql.Open("sqlite", "file:ent?mode=memory&_fk=1&_pragma=foreign_keys(1)")
	s.Assert().Nil(err)
	s.Assert().NotNil(db)
	defer db.Close()

	client := enttest.NewClient(s.T(), enttest.WithOptions(ent.Driver(dialect_sql.OpenDB("sqlite3", db))))
	defer client.Close()

	tag1, _ := client.Tag.Create().SetName("some artist").Save(context.Background())
	tag2, _ := client.Tag.Create().SetName("artist").Save(context.Background())

	client.Meta.Create().SetName("[some artist]manga 1 here.zip").SetCreateTime(time.UnixMilli(3000)).SetFavorite(true).AddTags(tag1).Save(context.Background())
	client.Meta.Create().SetName("[some artist]manga 2 here.zip").SetCreateTime(time.UnixMilli(4000)).SetFavorite(true).AddTags(tag1).Save(context.Background())
	client.Meta.Create().SetName("[artist]manga 3.zip").SetCreateTime(time.UnixMilli(5000)).SetFavorite(true).AddTags(tag2).Save(context.Background())
	client.Meta.Create().SetName("[artist]manga 4.zip").SetCreateTime(time.UnixMilli(6000)).SetFavorite(true).AddTags(tag2).Save(context.Background())

	var u *ent.User
	tags, err := ReadPage(context.Background(), client, u, QueryParams{
		SearchTag:   "some artist",
		SortBy:      grpc.SortField_SORT_FIELD_CREATION_TIME,
		SortOrder:   grpc.SortOrder_SORT_ORDER_ASCENDING,
		Filter:      grpc.Filter_FILTER_FAVORITE_ITEMS,
		Page:        0,
		ItemPerPage: 30,
	})
	s.Assert().Nil(err)

	s.Assert().Equal(2, len(tags))

	s.Assert().Equal("[some artist]manga 1 here.zip", tags[0].Name)
	s.Assert().Equal("[some artist]manga 2 here.zip", tags[1].Name)
}

func (s *QueryTestSuite) TestReadPageSearchNameTagFilterFavoriteItemsortByCreateTimeAsc() {
	db, err := sql.Open("sqlite", "file:ent?mode=memory&_fk=1&_pragma=foreign_keys(1)")
	s.Assert().Nil(err)
	s.Assert().NotNil(db)
	defer db.Close()

	client := enttest.NewClient(s.T(), enttest.WithOptions(ent.Driver(dialect_sql.OpenDB("sqlite3", db))))
	defer client.Close()

	tag1, _ := client.Tag.Create().SetName("some artist").Save(context.Background())
	tag2, _ := client.Tag.Create().SetName("artist").Save(context.Background())

	var u *ent.User
	client.Meta.Create().SetName("[some artist]manga 1 here.zip").SetCreateTime(time.UnixMilli(3000)).SetFavorite(true).AddTags(tag1).Save(context.Background())
	client.Meta.Create().SetName("[some artist]manga 2 here.zip").SetCreateTime(time.UnixMilli(4000)).SetFavorite(true).AddTags(tag1).Save(context.Background())
	client.Meta.Create().SetName("[some artist]manga 3.zip").SetCreateTime(time.UnixMilli(5000)).SetFavorite(true).AddTags(tag1).Save(context.Background())
	client.Meta.Create().SetName("[artist]manga 4.zip").SetCreateTime(time.UnixMilli(6000)).SetFavorite(true).AddTags(tag2).Save(context.Background())

	tags, err := ReadPage(context.Background(), client, u, QueryParams{
		SearchName:  "here",
		SearchTag:   "some artist",
		SortBy:      grpc.SortField_SORT_FIELD_CREATION_TIME,
		SortOrder:   grpc.SortOrder_SORT_ORDER_ASCENDING,
		Filter:      grpc.Filter_FILTER_FAVORITE_ITEMS,
		Page:        0,
		ItemPerPage: 30,
	})
	s.Assert().Nil(err)

	s.Assert().Equal(2, len(tags))

	s.Assert().Equal("[some artist]manga 1 here.zip", tags[0].Name)
	s.Assert().Equal("[some artist]manga 2 here.zip", tags[1].Name)
}

func (s *QueryTestSuite) TestCount() {
	db, err := sql.Open("sqlite", "file:ent?mode=memory&_fk=1&_pragma=foreign_keys(1)")
	s.Assert().Nil(err)
	s.Assert().NotNil(db)
	defer db.Close()

	client := enttest.NewClient(s.T(), enttest.WithOptions(ent.Driver(dialect_sql.OpenDB("sqlite3", db))))
	defer client.Close()

	client.Meta.Create().SetName("[some artist]manga 1 here.zip").Save(context.Background())
	client.Meta.Create().SetName("[some artist]manga 2 here.zip").Save(context.Background())
	client.Meta.Create().SetName("[some artist]manga 3 here.zip").Save(context.Background())
	client.Meta.Create().SetName("[some artist]manga 4 here.zip").Save(context.Background())
	client.Meta.Create().SetName("[some artist]manga 5 here.zip").SetActive(false).Save(context.Background())

	var u *ent.User
	c, err := Count(context.Background(), client, u, QueryParams{
		SortBy:      grpc.SortField_SORT_FIELD_NAME,
		SortOrder:   grpc.SortOrder_SORT_ORDER_ASCENDING,
		Page:        0,
		ItemPerPage: 30,
	})
	s.Assert().Nil(err)
	s.Assert().Equal(4, c)
}

func (s *QueryTestSuite) TestCountFilterFavoriteItems() {
	db, err := sql.Open("sqlite", "file:ent?mode=memory&_fk=1&_pragma=foreign_keys(1)")
	s.Assert().Nil(err)
	s.Assert().NotNil(db)
	defer db.Close()

	client := enttest.NewClient(s.T(), enttest.WithOptions(ent.Driver(dialect_sql.OpenDB("sqlite3", db))))
	defer client.Close()

	client.Meta.Create().SetName("[some artist]manga 1 here.zip").SetFavorite(true).Save(context.Background())
	client.Meta.Create().SetName("[some artist]manga 2 here.zip").SetFavorite(true).Save(context.Background())
	client.Meta.Create().SetName("[some artist]manga 3 here.zip").Save(context.Background())
	client.Meta.Create().SetName("[some artist]manga 4 here.zip").Save(context.Background())

	var u *ent.User
	c, err := Count(context.Background(), client, u, QueryParams{
		Filter:      grpc.Filter_FILTER_FAVORITE_ITEMS,
		SortBy:      grpc.SortField_SORT_FIELD_NAME,
		SortOrder:   grpc.SortOrder_SORT_ORDER_ASCENDING,
		Page:        0,
		ItemPerPage: 30,
	})
	s.Assert().Nil(err)
	s.Assert().Equal(2, c)
}

func (s *QueryTestSuite) TestCountSortByCreateTimeDesc() {
	db, err := sql.Open("sqlite", "file:ent?mode=memory&_fk=1&_pragma=foreign_keys(1)")
	s.Assert().Nil(err)
	s.Assert().NotNil(db)
	defer db.Close()

	client := enttest.NewClient(s.T(), enttest.WithOptions(ent.Driver(dialect_sql.OpenDB("sqlite3", db))))
	defer client.Close()

	client.Meta.Create().SetName("[some artist]manga 1 here.zip").SetCreateTime(time.UnixMilli(3000)).Save(context.Background())
	client.Meta.Create().SetName("[some artist]manga 2 here.zip").SetCreateTime(time.UnixMilli(4000)).Save(context.Background())
	client.Meta.Create().SetName("[some artist]manga 3 here.zip").SetCreateTime(time.UnixMilli(5000)).SetActive(false).Save(context.Background())
	client.Meta.Create().SetName("[some artist]manga 4 here.zip").SetCreateTime(time.UnixMilli(6000)).SetActive(false).Save(context.Background())

	var u *ent.User
	c, err := Count(context.Background(), client, u, QueryParams{
		Filter:      grpc.Filter_FILTER_UNKNOWN,
		SortBy:      grpc.SortField_SORT_FIELD_CREATION_TIME,
		SortOrder:   grpc.SortOrder_SORT_ORDER_DESCENDING,
		Page:        0,
		ItemPerPage: 30,
	})
	s.Assert().Nil(err)
	s.Assert().Equal(2, c)
}

func (s *QueryTestSuite) TestCountSortByCreateTimeAsc() {
	db, err := sql.Open("sqlite", "file:ent?mode=memory&_fk=1&_pragma=foreign_keys(1)")
	s.Assert().Nil(err)
	s.Assert().NotNil(db)
	defer db.Close()

	client := enttest.NewClient(s.T(), enttest.WithOptions(ent.Driver(dialect_sql.OpenDB("sqlite3", db))))
	defer client.Close()

	client.Meta.Create().SetName("[some artist]manga 1 here.zip").SetCreateTime(time.UnixMilli(3000)).Save(context.Background())
	client.Meta.Create().SetName("[some artist]manga 2 here.zip").SetCreateTime(time.UnixMilli(4000)).Save(context.Background())
	client.Meta.Create().SetName("[some artist]manga 3 here.zip").SetCreateTime(time.UnixMilli(5000)).SetActive(false).Save(context.Background())
	client.Meta.Create().SetName("[some artist]manga 4 here.zip").SetCreateTime(time.UnixMilli(6000)).SetActive(false).Save(context.Background())

	var u *ent.User
	c, err := Count(context.Background(), client, u, QueryParams{
		Filter:      grpc.Filter_FILTER_UNKNOWN,
		SortBy:      grpc.SortField_SORT_FIELD_CREATION_TIME,
		SortOrder:   grpc.SortOrder_SORT_ORDER_ASCENDING,
		Page:        0,
		ItemPerPage: 30,
	})
	s.Assert().Nil(err)
	s.Assert().Equal(2, c)
}

func (s *QueryTestSuite) TestCountFilterFavoriteItemsortByCreateTimeDesc() {
	db, err := sql.Open("sqlite", "file:ent?mode=memory&_fk=1&_pragma=foreign_keys(1)")
	s.Assert().Nil(err)
	s.Assert().NotNil(db)
	defer db.Close()

	client := enttest.NewClient(s.T(), enttest.WithOptions(ent.Driver(dialect_sql.OpenDB("sqlite3", db))))
	defer client.Close()

	client.Meta.Create().SetName("[some artist]manga 1 here.zip").SetCreateTime(time.UnixMilli(3000)).SetFavorite(true).Save(context.Background())
	client.Meta.Create().SetName("[some artist]manga 2 here.zip").SetCreateTime(time.UnixMilli(4000)).SetFavorite(true).Save(context.Background())
	client.Meta.Create().SetName("[some artist]manga 3 here.zip").SetCreateTime(time.UnixMilli(5000)).Save(context.Background())
	client.Meta.Create().SetName("[some artist]manga 4 here.zip").SetCreateTime(time.UnixMilli(6000)).Save(context.Background())

	var u *ent.User
	c, err := Count(context.Background(), client, u, QueryParams{
		Filter:      grpc.Filter_FILTER_FAVORITE_ITEMS,
		SortBy:      grpc.SortField_SORT_FIELD_CREATION_TIME,
		SortOrder:   grpc.SortOrder_SORT_ORDER_DESCENDING,
		Page:        0,
		ItemPerPage: 30,
	})
	s.Assert().Nil(err)
	s.Assert().Equal(2, c)
}

func (s *QueryTestSuite) TestCountFilterFavoriteItemsortByCreateTimeAsc() {
	db, err := sql.Open("sqlite", "file:ent?mode=memory&_fk=1&_pragma=foreign_keys(1)")
	s.Assert().Nil(err)
	s.Assert().NotNil(db)
	defer db.Close()

	client := enttest.NewClient(s.T(), enttest.WithOptions(ent.Driver(dialect_sql.OpenDB("sqlite3", db))))
	defer client.Close()

	client.Meta.Create().SetName("[some artist]manga 1 here.zip").SetCreateTime(time.UnixMilli(3000)).SetFavorite(true).Save(context.Background())
	client.Meta.Create().SetName("[some artist]manga 2 here.zip").SetCreateTime(time.UnixMilli(4000)).SetFavorite(true).Save(context.Background())
	client.Meta.Create().SetName("[some artist]manga 3 here.zip").SetCreateTime(time.UnixMilli(5000)).Save(context.Background())
	client.Meta.Create().SetName("[some artist]manga 4 here.zip").SetCreateTime(time.UnixMilli(6000)).Save(context.Background())

	var u *ent.User
	c, err := Count(context.Background(), client, u, QueryParams{
		Filter:      grpc.Filter_FILTER_FAVORITE_ITEMS,
		SortBy:      grpc.SortField_SORT_FIELD_CREATION_TIME,
		SortOrder:   grpc.SortOrder_SORT_ORDER_ASCENDING,
		Page:        0,
		ItemPerPage: 30,
	})
	s.Assert().Nil(err)
	s.Assert().Equal(2, c)
}

func (s *QueryTestSuite) TestCountSearchNameFilterFavoriteItemsortByCreateTimeDesc() {
	db, err := sql.Open("sqlite", "file:ent?mode=memory&_fk=1&_pragma=foreign_keys(1)")
	s.Assert().Nil(err)
	s.Assert().NotNil(db)
	defer db.Close()

	client := enttest.NewClient(s.T(), enttest.WithOptions(ent.Driver(dialect_sql.OpenDB("sqlite3", db))))
	defer client.Close()

	client.Meta.Create().SetName("[some artist]manga 1 here.zip").SetCreateTime(time.UnixMilli(3000)).SetFavorite(true).Save(context.Background())
	client.Meta.Create().SetName("[some artist]manga 2 here.zip").SetCreateTime(time.UnixMilli(4000)).SetFavorite(true).Save(context.Background())
	client.Meta.Create().SetName("[some artist]manga 3.zip").SetCreateTime(time.UnixMilli(5000)).SetFavorite(true).Save(context.Background())
	client.Meta.Create().SetName("[some artist]manga 4.zip").SetCreateTime(time.UnixMilli(6000)).SetFavorite(true).Save(context.Background())

	var u *ent.User
	c, err := Count(context.Background(), client, u, QueryParams{
		SearchName:  "here",
		SortBy:      grpc.SortField_SORT_FIELD_CREATION_TIME,
		SortOrder:   grpc.SortOrder_SORT_ORDER_DESCENDING,
		Filter:      grpc.Filter_FILTER_FAVORITE_ITEMS,
		Page:        0,
		ItemPerPage: 30,
	})
	s.Assert().Nil(err)
	s.Assert().Equal(2, c)
}

func (s *QueryTestSuite) TestCountSearchNameFilterFavoriteItemsortByCreateTimeAsc() {
	db, err := sql.Open("sqlite", "file:ent?mode=memory&_fk=1&_pragma=foreign_keys(1)")
	s.Assert().Nil(err)
	s.Assert().NotNil(db)
	defer db.Close()

	client := enttest.NewClient(s.T(), enttest.WithOptions(ent.Driver(dialect_sql.OpenDB("sqlite3", db))))
	defer client.Close()

	client.Meta.Create().SetName("[some artist]manga 1 here.zip").SetCreateTime(time.UnixMilli(3000)).SetFavorite(true).Save(context.Background())
	client.Meta.Create().SetName("[some artist]manga 2 here.zip").SetCreateTime(time.UnixMilli(4000)).SetFavorite(true).Save(context.Background())
	client.Meta.Create().SetName("[some artist]manga 3.zip").SetCreateTime(time.UnixMilli(5000)).SetFavorite(true).Save(context.Background())
	client.Meta.Create().SetName("[some artist]manga 4.zip").SetCreateTime(time.UnixMilli(6000)).SetFavorite(true).Save(context.Background())

	var u *ent.User
	c, err := Count(context.Background(), client, u, QueryParams{
		SearchName:  "here",
		SortBy:      grpc.SortField_SORT_FIELD_CREATION_TIME,
		SortOrder:   grpc.SortOrder_SORT_ORDER_ASCENDING,
		Filter:      grpc.Filter_FILTER_FAVORITE_ITEMS,
		Page:        0,
		ItemPerPage: 30,
	})
	s.Assert().Nil(err)

	s.Assert().Equal(2, c)
}

func (s *QueryTestSuite) TestCountSearchTagFilterFavoriteItemsortByCreateTimeAsc() {
	db, err := sql.Open("sqlite", "file:ent?mode=memory&_fk=1&_pragma=foreign_keys(1)")
	s.Assert().Nil(err)
	s.Assert().NotNil(db)
	defer db.Close()

	client := enttest.NewClient(s.T(), enttest.WithOptions(ent.Driver(dialect_sql.OpenDB("sqlite3", db))))
	defer client.Close()

	tag1, _ := client.Tag.Create().SetName("some artist").Save(context.Background())
	tag2, _ := client.Tag.Create().SetName("artist").Save(context.Background())

	client.Meta.Create().SetName("[some artist]manga 1 here.zip").SetCreateTime(time.UnixMilli(3000)).SetFavorite(true).AddTags(tag1).Save(context.Background())
	client.Meta.Create().SetName("[some artist]manga 2 here.zip").SetCreateTime(time.UnixMilli(4000)).SetFavorite(true).AddTags(tag1).Save(context.Background())
	client.Meta.Create().SetName("[artist]manga 3.zip").SetCreateTime(time.UnixMilli(5000)).SetFavorite(true).AddTags(tag2).Save(context.Background())
	client.Meta.Create().SetName("[artist]manga 4.zip").SetCreateTime(time.UnixMilli(6000)).SetFavorite(true).AddTags(tag2).Save(context.Background())

	var u *ent.User
	c, err := Count(context.Background(), client, u, QueryParams{
		SearchTag:   "some artist",
		SortBy:      grpc.SortField_SORT_FIELD_CREATION_TIME,
		SortOrder:   grpc.SortOrder_SORT_ORDER_ASCENDING,
		Filter:      grpc.Filter_FILTER_FAVORITE_ITEMS,
		Page:        0,
		ItemPerPage: 30,
	})
	s.Assert().Nil(err)
	s.Assert().Equal(2, c)
}

func (s *QueryTestSuite) TestCountSearchNameTagFilterFavoriteItemsortByCreateTimeAsc() {
	db, err := sql.Open("sqlite", "file:ent?mode=memory&_fk=1&_pragma=foreign_keys(1)")
	s.Assert().Nil(err)
	s.Assert().NotNil(db)
	defer db.Close()

	client := enttest.NewClient(s.T(), enttest.WithOptions(ent.Driver(dialect_sql.OpenDB("sqlite3", db))))
	defer client.Close()

	tag1, _ := client.Tag.Create().SetName("some artist").Save(context.Background())
	tag2, _ := client.Tag.Create().SetName("artist").Save(context.Background())

	client.Meta.Create().SetName("[some artist]manga 1 here.zip").SetCreateTime(time.UnixMilli(3000)).SetFavorite(true).AddTags(tag1).Save(context.Background())
	client.Meta.Create().SetName("[some artist]manga 2 here.zip").SetCreateTime(time.UnixMilli(4000)).SetFavorite(true).AddTags(tag1).Save(context.Background())
	client.Meta.Create().SetName("[some artist]manga 3.zip").SetCreateTime(time.UnixMilli(5000)).SetFavorite(true).AddTags(tag1).Save(context.Background())
	client.Meta.Create().SetName("[artist]manga 4.zip").SetCreateTime(time.UnixMilli(6000)).SetFavorite(true).AddTags(tag2).Save(context.Background())

	var u *ent.User
	c, err := Count(context.Background(), client, u, QueryParams{
		SearchName:  "here",
		SearchTag:   "some artist",
		SortBy:      grpc.SortField_SORT_FIELD_CREATION_TIME,
		SortOrder:   grpc.SortOrder_SORT_ORDER_ASCENDING,
		Filter:      grpc.Filter_FILTER_FAVORITE_ITEMS,
		Page:        0,
		ItemPerPage: 30,
	})
	s.Assert().Nil(err)
	s.Assert().Equal(2, c)
}

func (s *QueryTestSuite) TestReadSortByPageCountAsc() {
	db, err := sql.Open("sqlite", "file:ent?mode=memory&_fk=1&_pragma=foreign_keys(1)")
	s.Assert().Nil(err)
	s.Assert().NotNil(db)
	defer db.Close()

	client := enttest.NewClient(s.T(), enttest.WithOptions(ent.Driver(dialect_sql.OpenDB("sqlite3", db))))
	defer client.Close()

	client.Meta.Create().SetName("[some artist]manga 1 here.zip").SetFileIndices([]int{1}).Save(context.Background())
	client.Meta.Create().SetName("[some artist]manga 2 here.zip").SetFileIndices([]int{1, 2}).Save(context.Background())
	client.Meta.Create().SetName("[some artist]manga 3 here.zip").SetFileIndices([]int{1, 2, 3}).Save(context.Background())
	client.Meta.Create().SetName("[some artist]manga 4 here.zip").SetFileIndices([]int{1, 2, 3, 4}).Save(context.Background())
	client.Meta.Create().SetName("[some artist]manga 5 here.zip").SetFileIndices([]int{1, 2, 3, 4}).SetActive(false).Save(context.Background())

	var u *ent.User
	tags, err := ReadPage(context.Background(), client, u, QueryParams{
		SortBy:      grpc.SortField_SORT_FIELD_NAME,
		SortOrder:   grpc.SortOrder_SORT_ORDER_ASCENDING,
		Page:        0,
		ItemPerPage: 30,
	})
	s.Assert().Nil(err)

	s.Assert().Equal(4, len(tags))

	s.Assert().Equal("[some artist]manga 1 here.zip", tags[0].Name)
	s.Assert().Equal("[some artist]manga 2 here.zip", tags[1].Name)
	s.Assert().Equal("[some artist]manga 3 here.zip", tags[2].Name)
	s.Assert().Equal("[some artist]manga 4 here.zip", tags[3].Name)
}

func (s *QueryTestSuite) TestReadSortByPageCountDesc() {
	db, err := sql.Open("sqlite", "file:ent?mode=memory&_fk=1&_pragma=foreign_keys(1)")
	s.Assert().Nil(err)
	s.Assert().NotNil(db)
	defer db.Close()

	client := enttest.NewClient(s.T(), enttest.WithOptions(ent.Driver(dialect_sql.OpenDB("sqlite3", db))))
	defer client.Close()

	client.Meta.Create().SetName("[some artist]manga 1 here.zip").SetFileIndices([]int{1, 2, 3, 4}).Save(context.Background())
	client.Meta.Create().SetName("[some artist]manga 2 here.zip").SetFileIndices([]int{1, 2, 3}).Save(context.Background())
	client.Meta.Create().SetName("[some artist]manga 3 here.zip").SetFileIndices([]int{1, 2}).Save(context.Background())
	client.Meta.Create().SetName("[some artist]manga 4 here.zip").SetFileIndices([]int{1}).Save(context.Background())
	client.Meta.Create().SetName("[some artist]manga 5 here.zip").SetFileIndices([]int{1, 2, 3, 4}).SetActive(false).Save(context.Background())

	var u *ent.User
	tags, err := ReadPage(context.Background(), client, u, QueryParams{
		SortBy:      grpc.SortField_SORT_FIELD_NAME,
		SortOrder:   grpc.SortOrder_SORT_ORDER_ASCENDING,
		Page:        0,
		ItemPerPage: 30,
	})
	s.Assert().Nil(err)

	s.Assert().Equal(4, len(tags))

	s.Assert().Equal("[some artist]manga 1 here.zip", tags[0].Name)
	s.Assert().Equal("[some artist]manga 2 here.zip", tags[1].Name)
	s.Assert().Equal("[some artist]manga 3 here.zip", tags[2].Name)
	s.Assert().Equal("[some artist]manga 4 here.zip", tags[3].Name)
}
