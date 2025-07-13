package meta

import (
	"context"
	"database/sql"
	"testing"
	"time"

	dialect_sql "entgo.io/ent/dialect/sql"
	"github.com/mangaweb4/mangaweb4-backend/ent"
	"github.com/mangaweb4/mangaweb4-backend/ent/enttest"
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
		SortBy:      SortFieldName,
		SortOrder:   SortOrderAscending,
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

func (s *QueryTestSuite) TestReadPageFilterFavoriteItem() {
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
		Filter:      FilterFavoriteItem,
		SortBy:      SortFieldName,
		SortOrder:   SortOrderAscending,
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
		Filter:      FilterNone,
		SortBy:      SortFieldCreateTime,
		SortOrder:   SortOrderDescending,
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
		Filter:      FilterNone,
		SortBy:      SortFieldCreateTime,
		SortOrder:   SortOrderAscending,
		Page:        0,
		ItemPerPage: 30,
	})
	s.Assert().Nil(err)

	s.Assert().Equal(2, len(tags))

	s.Assert().Equal("[some artist]manga 1 here.zip", tags[0].Name)
	s.Assert().Equal("[some artist]manga 2 here.zip", tags[1].Name)
}

func (s *QueryTestSuite) TestReadPageFilterFavoriteItemSortByCreateTimeDesc() {
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
		Filter:      FilterFavoriteItem,
		SortBy:      SortFieldCreateTime,
		SortOrder:   SortOrderDescending,
		Page:        0,
		ItemPerPage: 30,
	})
	s.Assert().Nil(err)

	s.Assert().Equal(2, len(tags))

	s.Assert().Equal("[some artist]manga 2 here.zip", tags[0].Name)
	s.Assert().Equal("[some artist]manga 1 here.zip", tags[1].Name)
}

func (s *QueryTestSuite) TestReadPageFilterFavoriteItemSortByCreateTimeAsc() {
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
		Filter:      FilterFavoriteItem,
		SortBy:      SortFieldCreateTime,
		SortOrder:   SortOrderAscending,
		Page:        0,
		ItemPerPage: 30,
	})
	s.Assert().Nil(err)

	s.Assert().Equal(2, len(tags))

	s.Assert().Equal("[some artist]manga 1 here.zip", tags[0].Name)
	s.Assert().Equal("[some artist]manga 2 here.zip", tags[1].Name)
}

func (s *QueryTestSuite) TestReadPageSearchNameFilterFavoriteItemSortByCreateTimeDesc() {
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
		SortBy:      SortFieldCreateTime,
		SortOrder:   SortOrderDescending,
		Filter:      FilterFavoriteItem,
		Page:        0,
		ItemPerPage: 30,
	})
	s.Assert().Nil(err)

	s.Assert().Equal(2, len(tags))

	s.Assert().Equal("[some artist]manga 2 here.zip", tags[0].Name)
	s.Assert().Equal("[some artist]manga 1 here.zip", tags[1].Name)
}

func (s *QueryTestSuite) TestReadPageSearchNameFilterFavoriteItemSortByCreateTimeAsc() {
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
		SortBy:      SortFieldCreateTime,
		SortOrder:   SortOrderAscending,
		Filter:      FilterFavoriteItem,
		Page:        0,
		ItemPerPage: 30,
	})
	s.Assert().Nil(err)

	s.Assert().Equal(2, len(tags))

	s.Assert().Equal("[some artist]manga 1 here.zip", tags[0].Name)
	s.Assert().Equal("[some artist]manga 2 here.zip", tags[1].Name)
}

func (s *QueryTestSuite) TestReadPageSearchTagFilterFavoriteItemSortByCreateTimeAsc() {
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
		SortBy:      SortFieldCreateTime,
		SortOrder:   SortOrderAscending,
		Filter:      FilterFavoriteItem,
		Page:        0,
		ItemPerPage: 30,
	})
	s.Assert().Nil(err)

	s.Assert().Equal(2, len(tags))

	s.Assert().Equal("[some artist]manga 1 here.zip", tags[0].Name)
	s.Assert().Equal("[some artist]manga 2 here.zip", tags[1].Name)
}

func (s *QueryTestSuite) TestReadPageSearchNameTagFilterFavoriteItemSortByCreateTimeAsc() {
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
		SortBy:      SortFieldCreateTime,
		SortOrder:   SortOrderAscending,
		Filter:      FilterFavoriteItem,
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
		SortBy:      SortFieldName,
		SortOrder:   SortOrderAscending,
		Page:        0,
		ItemPerPage: 30,
	})
	s.Assert().Nil(err)
	s.Assert().Equal(4, c)
}

func (s *QueryTestSuite) TestCountFilterFavoriteItem() {
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
		Filter:      FilterFavoriteItem,
		SortBy:      SortFieldName,
		SortOrder:   SortOrderAscending,
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
		Filter:      FilterNone,
		SortBy:      SortFieldCreateTime,
		SortOrder:   SortOrderDescending,
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
		Filter:      FilterNone,
		SortBy:      SortFieldCreateTime,
		SortOrder:   SortOrderAscending,
		Page:        0,
		ItemPerPage: 30,
	})
	s.Assert().Nil(err)
	s.Assert().Equal(2, c)
}

func (s *QueryTestSuite) TestCountFilterFavoriteItemSortByCreateTimeDesc() {
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
		Filter:      FilterFavoriteItem,
		SortBy:      SortFieldCreateTime,
		SortOrder:   SortOrderDescending,
		Page:        0,
		ItemPerPage: 30,
	})
	s.Assert().Nil(err)
	s.Assert().Equal(2, c)
}

func (s *QueryTestSuite) TestCountFilterFavoriteItemSortByCreateTimeAsc() {
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
		Filter:      FilterFavoriteItem,
		SortBy:      SortFieldCreateTime,
		SortOrder:   SortOrderAscending,
		Page:        0,
		ItemPerPage: 30,
	})
	s.Assert().Nil(err)
	s.Assert().Equal(2, c)
}

func (s *QueryTestSuite) TestCountSearchNameFilterFavoriteItemSortByCreateTimeDesc() {
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
		SortBy:      SortFieldCreateTime,
		SortOrder:   SortOrderDescending,
		Filter:      FilterFavoriteItem,
		Page:        0,
		ItemPerPage: 30,
	})
	s.Assert().Nil(err)
	s.Assert().Equal(2, c)
}

func (s *QueryTestSuite) TestCountSearchNameFilterFavoriteItemSortByCreateTimeAsc() {
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
		SortBy:      SortFieldCreateTime,
		SortOrder:   SortOrderAscending,
		Filter:      FilterFavoriteItem,
		Page:        0,
		ItemPerPage: 30,
	})
	s.Assert().Nil(err)

	s.Assert().Equal(2, c)
}

func (s *QueryTestSuite) TestCountSearchTagFilterFavoriteItemSortByCreateTimeAsc() {
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
		SortBy:      SortFieldCreateTime,
		SortOrder:   SortOrderAscending,
		Filter:      FilterFavoriteItem,
		Page:        0,
		ItemPerPage: 30,
	})
	s.Assert().Nil(err)
	s.Assert().Equal(2, c)
}

func (s *QueryTestSuite) TestCountSearchNameTagFilterFavoriteItemSortByCreateTimeAsc() {
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
		SortBy:      SortFieldCreateTime,
		SortOrder:   SortOrderAscending,
		Filter:      FilterFavoriteItem,
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
		SortBy:      SortFieldName,
		SortOrder:   SortOrderAscending,
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
		SortBy:      SortFieldName,
		SortOrder:   SortOrderAscending,
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
