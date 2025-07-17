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
	"github.com/mangaweb4/mangaweb4-backend/user"
	"github.com/stretchr/testify/suite"
	_ "modernc.org/sqlite"
)

type QueryTestSuite struct {
	suite.Suite
}

func TestQueryTestSuite(t *testing.T) {
	suite.Run(t, new(QueryTestSuite))
}

func createTestDBClient(s *QueryTestSuite) (db *sql.DB, client *ent.Client, err error) {
	db, err = sql.Open("sqlite", "file:ent?mode=memory&_fk=1&_pragma=foreign_keys(1)")
	if err != nil {
		return
	}

	client = enttest.NewClient(s.T(), enttest.WithOptions(ent.Driver(dialect_sql.OpenDB("sqlite3", db))))

	return
}

func (s *QueryTestSuite) TestReadPage() {
	db, client, err := createTestDBClient(s)
	s.Assert().Nil(err)
	s.Assert().NotNil(db)
	s.Assert().NotNil(client)
	defer db.Close()
	defer client.Close()

	u, err := user.GetUser(context.Background(), client, "")
	s.Assert().Nil(err)

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
	db, client, err := createTestDBClient(s)
	s.Assert().Nil(err)
	s.Assert().NotNil(db)
	s.Assert().NotNil(client)
	defer db.Close()
	defer client.Close()

	u, err := user.GetUser(context.Background(), client, "")
	s.Assert().Nil(err)

	client.Meta.Create().SetName("[some artist]manga 1 here.zip").AddFavoriteOfUser(u).Save(context.Background())
	client.Meta.Create().SetName("[some artist]manga 2 here.zip").AddFavoriteOfUser(u).Save(context.Background())
	client.Meta.Create().SetName("[some artist]manga 3 here.zip").Save(context.Background())
	client.Meta.Create().SetName("[some artist]manga 4 here.zip").Save(context.Background())

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
	db, client, err := createTestDBClient(s)
	s.Assert().Nil(err)
	s.Assert().NotNil(db)
	s.Assert().NotNil(client)
	defer db.Close()
	defer client.Close()

	client.Meta.Create().SetName("[some artist]manga 1 here.zip").SetCreateTime(time.UnixMilli(3000)).Save(context.Background())
	client.Meta.Create().SetName("[some artist]manga 2 here.zip").SetCreateTime(time.UnixMilli(4000)).Save(context.Background())
	client.Meta.Create().SetName("[some artist]manga 3 here.zip").SetCreateTime(time.UnixMilli(5000)).SetActive(false).Save(context.Background())
	client.Meta.Create().SetName("[some artist]manga 4 here.zip").SetCreateTime(time.UnixMilli(6000)).SetActive(false).Save(context.Background())

	u, err := user.GetUser(context.Background(), client, "")
	s.Assert().Nil(err)
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
	db, client, err := createTestDBClient(s)
	s.Assert().Nil(err)
	s.Assert().NotNil(db)
	s.Assert().NotNil(client)
	defer db.Close()
	defer client.Close()

	client.Meta.Create().SetName("[some artist]manga 1 here.zip").SetCreateTime(time.UnixMilli(3000)).Save(context.Background())
	client.Meta.Create().SetName("[some artist]manga 2 here.zip").SetCreateTime(time.UnixMilli(4000)).Save(context.Background())
	client.Meta.Create().SetName("[some artist]manga 3 here.zip").SetCreateTime(time.UnixMilli(5000)).SetActive(false).Save(context.Background())
	client.Meta.Create().SetName("[some artist]manga 4 here.zip").SetCreateTime(time.UnixMilli(6000)).SetActive(false).Save(context.Background())

	u, err := user.GetUser(context.Background(), client, "")
	s.Assert().Nil(err)
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
	db, client, err := createTestDBClient(s)
	s.Assert().Nil(err)
	s.Assert().NotNil(db)
	s.Assert().NotNil(client)
	defer db.Close()
	defer client.Close()

	u, err := user.GetUser(context.Background(), client, "")
	s.Assert().Nil(err)

	client.Meta.Create().SetName("[some artist]manga 1 here.zip").SetCreateTime(time.UnixMilli(3000)).AddFavoriteOfUser(u).Save(context.Background())
	client.Meta.Create().SetName("[some artist]manga 2 here.zip").SetCreateTime(time.UnixMilli(4000)).AddFavoriteOfUser(u).Save(context.Background())
	client.Meta.Create().SetName("[some artist]manga 3 here.zip").SetCreateTime(time.UnixMilli(5000)).Save(context.Background())
	client.Meta.Create().SetName("[some artist]manga 4 here.zip").SetCreateTime(time.UnixMilli(6000)).Save(context.Background())

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
	db, client, err := createTestDBClient(s)
	s.Assert().Nil(err)
	s.Assert().NotNil(db)
	s.Assert().NotNil(client)
	defer db.Close()
	defer client.Close()

	u, err := user.GetUser(context.Background(), client, "")
	s.Assert().Nil(err)

	client.Meta.Create().SetName("[some artist]manga 1 here.zip").SetCreateTime(time.UnixMilli(3000)).AddFavoriteOfUser(u).Save(context.Background())
	client.Meta.Create().SetName("[some artist]manga 2 here.zip").SetCreateTime(time.UnixMilli(4000)).AddFavoriteOfUser(u).Save(context.Background())
	client.Meta.Create().SetName("[some artist]manga 3 here.zip").SetCreateTime(time.UnixMilli(5000)).Save(context.Background())
	client.Meta.Create().SetName("[some artist]manga 4 here.zip").SetCreateTime(time.UnixMilli(6000)).Save(context.Background())

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
	db, client, err := createTestDBClient(s)
	s.Assert().Nil(err)
	s.Assert().NotNil(db)
	s.Assert().NotNil(client)
	defer db.Close()
	defer client.Close()

	u, err := user.GetUser(context.Background(), client, "")
	s.Assert().Nil(err)

	client.Meta.Create().SetName("[some artist]manga 1 here.zip").SetCreateTime(time.UnixMilli(3000)).AddFavoriteOfUser(u).Save(context.Background())
	client.Meta.Create().SetName("[some artist]manga 2 here.zip").SetCreateTime(time.UnixMilli(4000)).AddFavoriteOfUser(u).Save(context.Background())
	client.Meta.Create().SetName("[some artist]manga 3.zip").SetCreateTime(time.UnixMilli(5000)).AddFavoriteOfUser(u).Save(context.Background())
	client.Meta.Create().SetName("[some artist]manga 4.zip").SetCreateTime(time.UnixMilli(6000)).AddFavoriteOfUser(u).Save(context.Background())

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
	db, client, err := createTestDBClient(s)
	s.Assert().Nil(err)
	s.Assert().NotNil(db)
	s.Assert().NotNil(client)
	defer db.Close()
	defer client.Close()

	u, err := user.GetUser(context.Background(), client, "")
	s.Assert().Nil(err)

	client.Meta.Create().SetName("[some artist]manga 1 here.zip").SetCreateTime(time.UnixMilli(3000)).AddFavoriteOfUser(u).Save(context.Background())
	client.Meta.Create().SetName("[some artist]manga 2 here.zip").SetCreateTime(time.UnixMilli(4000)).AddFavoriteOfUser(u).Save(context.Background())
	client.Meta.Create().SetName("[some artist]manga 3.zip").SetCreateTime(time.UnixMilli(5000)).AddFavoriteOfUser(u).Save(context.Background())
	client.Meta.Create().SetName("[some artist]manga 4.zip").SetCreateTime(time.UnixMilli(6000)).AddFavoriteOfUser(u).Save(context.Background())

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
	db, client, err := createTestDBClient(s)
	s.Assert().Nil(err)
	s.Assert().NotNil(db)
	s.Assert().NotNil(client)
	defer db.Close()
	defer client.Close()

	tag1, _ := client.Tag.Create().SetName("some artist").Save(context.Background())
	tag2, _ := client.Tag.Create().SetName("artist").Save(context.Background())

	u, err := user.GetUser(context.Background(), client, "")
	s.Assert().Nil(err)

	client.Meta.Create().SetName("[some artist]manga 1 here.zip").SetCreateTime(time.UnixMilli(3000)).AddFavoriteOfUser(u).AddTags(tag1).Save(context.Background())
	client.Meta.Create().SetName("[some artist]manga 2 here.zip").SetCreateTime(time.UnixMilli(4000)).AddFavoriteOfUser(u).AddTags(tag1).Save(context.Background())
	client.Meta.Create().SetName("[artist]manga 3.zip").SetCreateTime(time.UnixMilli(5000)).AddFavoriteOfUser(u).AddTags(tag2).Save(context.Background())
	client.Meta.Create().SetName("[artist]manga 4.zip").SetCreateTime(time.UnixMilli(6000)).AddFavoriteOfUser(u).AddTags(tag2).Save(context.Background())

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
	db, client, err := createTestDBClient(s)
	s.Assert().Nil(err)
	s.Assert().NotNil(db)
	s.Assert().NotNil(client)
	defer db.Close()
	defer client.Close()

	tag1, _ := client.Tag.Create().SetName("some artist").Save(context.Background())
	tag2, _ := client.Tag.Create().SetName("artist").Save(context.Background())

	u, err := user.GetUser(context.Background(), client, "")
	s.Assert().Nil(err)
	client.Meta.Create().SetName("[some artist]manga 1 here.zip").SetCreateTime(time.UnixMilli(3000)).AddFavoriteOfUser(u).AddTags(tag1).Save(context.Background())
	client.Meta.Create().SetName("[some artist]manga 2 here.zip").SetCreateTime(time.UnixMilli(4000)).AddFavoriteOfUser(u).AddTags(tag1).Save(context.Background())
	client.Meta.Create().SetName("[some artist]manga 3.zip").SetCreateTime(time.UnixMilli(5000)).AddFavoriteOfUser(u).AddTags(tag1).Save(context.Background())
	client.Meta.Create().SetName("[artist]manga 4.zip").SetCreateTime(time.UnixMilli(6000)).AddFavoriteOfUser(u).AddTags(tag2).Save(context.Background())

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
	db, client, err := createTestDBClient(s)
	s.Assert().Nil(err)
	s.Assert().NotNil(db)
	s.Assert().NotNil(client)
	defer db.Close()
	defer client.Close()

	client.Meta.Create().SetName("[some artist]manga 1 here.zip").Save(context.Background())
	client.Meta.Create().SetName("[some artist]manga 2 here.zip").Save(context.Background())
	client.Meta.Create().SetName("[some artist]manga 3 here.zip").Save(context.Background())
	client.Meta.Create().SetName("[some artist]manga 4 here.zip").Save(context.Background())
	client.Meta.Create().SetName("[some artist]manga 5 here.zip").SetActive(false).Save(context.Background())

	u, err := user.GetUser(context.Background(), client, "")
	s.Assert().Nil(err)
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
	db, client, err := createTestDBClient(s)
	s.Assert().Nil(err)
	s.Assert().NotNil(db)
	s.Assert().NotNil(client)
	defer db.Close()
	defer client.Close()

	u, err := user.GetUser(context.Background(), client, "")
	s.Assert().Nil(err)

	client.Meta.Create().SetName("[some artist]manga 1 here.zip").AddFavoriteOfUser(u).Save(context.Background())
	client.Meta.Create().SetName("[some artist]manga 2 here.zip").AddFavoriteOfUser(u).Save(context.Background())
	client.Meta.Create().SetName("[some artist]manga 3 here.zip").Save(context.Background())
	client.Meta.Create().SetName("[some artist]manga 4 here.zip").Save(context.Background())

	s.Assert().Nil(err)
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
	db, client, err := createTestDBClient(s)
	s.Assert().Nil(err)
	s.Assert().NotNil(db)
	s.Assert().NotNil(client)
	defer db.Close()
	defer client.Close()

	client.Meta.Create().SetName("[some artist]manga 1 here.zip").SetCreateTime(time.UnixMilli(3000)).Save(context.Background())
	client.Meta.Create().SetName("[some artist]manga 2 here.zip").SetCreateTime(time.UnixMilli(4000)).Save(context.Background())
	client.Meta.Create().SetName("[some artist]manga 3 here.zip").SetCreateTime(time.UnixMilli(5000)).SetActive(false).Save(context.Background())
	client.Meta.Create().SetName("[some artist]manga 4 here.zip").SetCreateTime(time.UnixMilli(6000)).SetActive(false).Save(context.Background())

	u, err := user.GetUser(context.Background(), client, "")
	s.Assert().Nil(err)
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
	db, client, err := createTestDBClient(s)
	s.Assert().Nil(err)
	s.Assert().NotNil(db)
	s.Assert().NotNil(client)
	defer db.Close()
	defer client.Close()

	client.Meta.Create().SetName("[some artist]manga 1 here.zip").SetCreateTime(time.UnixMilli(3000)).Save(context.Background())
	client.Meta.Create().SetName("[some artist]manga 2 here.zip").SetCreateTime(time.UnixMilli(4000)).Save(context.Background())
	client.Meta.Create().SetName("[some artist]manga 3 here.zip").SetCreateTime(time.UnixMilli(5000)).SetActive(false).Save(context.Background())
	client.Meta.Create().SetName("[some artist]manga 4 here.zip").SetCreateTime(time.UnixMilli(6000)).SetActive(false).Save(context.Background())

	u, err := user.GetUser(context.Background(), client, "")
	s.Assert().Nil(err)
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
	db, client, err := createTestDBClient(s)
	s.Assert().Nil(err)
	s.Assert().NotNil(db)
	s.Assert().NotNil(client)
	defer db.Close()
	defer client.Close()

	u, err := user.GetUser(context.Background(), client, "")
	s.Assert().Nil(err)

	client.Meta.Create().SetName("[some artist]manga 1 here.zip").SetCreateTime(time.UnixMilli(3000)).AddFavoriteOfUser(u).Save(context.Background())
	client.Meta.Create().SetName("[some artist]manga 2 here.zip").SetCreateTime(time.UnixMilli(4000)).AddFavoriteOfUser(u).Save(context.Background())
	client.Meta.Create().SetName("[some artist]manga 3 here.zip").SetCreateTime(time.UnixMilli(5000)).Save(context.Background())
	client.Meta.Create().SetName("[some artist]manga 4 here.zip").SetCreateTime(time.UnixMilli(6000)).Save(context.Background())

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
	db, client, err := createTestDBClient(s)
	s.Assert().Nil(err)
	s.Assert().NotNil(db)
	s.Assert().NotNil(client)
	defer db.Close()
	defer client.Close()

	u, err := user.GetUser(context.Background(), client, "")
	s.Assert().Nil(err)

	client.Meta.Create().SetName("[some artist]manga 1 here.zip").SetCreateTime(time.UnixMilli(3000)).AddFavoriteOfUser(u).Save(context.Background())
	client.Meta.Create().SetName("[some artist]manga 2 here.zip").SetCreateTime(time.UnixMilli(4000)).AddFavoriteOfUser(u).Save(context.Background())
	client.Meta.Create().SetName("[some artist]manga 3 here.zip").SetCreateTime(time.UnixMilli(5000)).Save(context.Background())
	client.Meta.Create().SetName("[some artist]manga 4 here.zip").SetCreateTime(time.UnixMilli(6000)).Save(context.Background())

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
	db, client, err := createTestDBClient(s)
	s.Assert().Nil(err)
	s.Assert().NotNil(db)
	s.Assert().NotNil(client)
	defer db.Close()
	defer client.Close()

	u, err := user.GetUser(context.Background(), client, "")
	s.Assert().Nil(err)

	client.Meta.Create().SetName("[some artist]manga 1 here.zip").SetCreateTime(time.UnixMilli(3000)).AddFavoriteOfUser(u).Save(context.Background())
	client.Meta.Create().SetName("[some artist]manga 2 here.zip").SetCreateTime(time.UnixMilli(4000)).AddFavoriteOfUser(u).Save(context.Background())
	client.Meta.Create().SetName("[some artist]manga 3.zip").SetCreateTime(time.UnixMilli(5000)).AddFavoriteOfUser(u).Save(context.Background())
	client.Meta.Create().SetName("[some artist]manga 4.zip").SetCreateTime(time.UnixMilli(6000)).AddFavoriteOfUser(u).Save(context.Background())

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
	db, client, err := createTestDBClient(s)
	s.Assert().Nil(err)
	s.Assert().NotNil(db)
	s.Assert().NotNil(client)
	defer db.Close()
	defer client.Close()

	u, err := user.GetUser(context.Background(), client, "")
	s.Assert().Nil(err)

	client.Meta.Create().SetName("[some artist]manga 1 here.zip").SetCreateTime(time.UnixMilli(3000)).AddFavoriteOfUser(u).Save(context.Background())
	client.Meta.Create().SetName("[some artist]manga 2 here.zip").SetCreateTime(time.UnixMilli(4000)).AddFavoriteOfUser(u).Save(context.Background())
	client.Meta.Create().SetName("[some artist]manga 3.zip").SetCreateTime(time.UnixMilli(5000)).AddFavoriteOfUser(u).Save(context.Background())
	client.Meta.Create().SetName("[some artist]manga 4.zip").SetCreateTime(time.UnixMilli(6000)).AddFavoriteOfUser(u).Save(context.Background())

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
	db, client, err := createTestDBClient(s)
	s.Assert().Nil(err)
	s.Assert().NotNil(db)
	s.Assert().NotNil(client)
	defer db.Close()
	defer client.Close()

	tag1, _ := client.Tag.Create().SetName("some artist").Save(context.Background())
	tag2, _ := client.Tag.Create().SetName("artist").Save(context.Background())

	u, err := user.GetUser(context.Background(), client, "")
	s.Assert().Nil(err)

	client.Meta.Create().SetName("[some artist]manga 1 here.zip").SetCreateTime(time.UnixMilli(3000)).AddFavoriteOfUser(u).AddTags(tag1).Save(context.Background())
	client.Meta.Create().SetName("[some artist]manga 2 here.zip").SetCreateTime(time.UnixMilli(4000)).AddFavoriteOfUser(u).AddTags(tag1).Save(context.Background())
	client.Meta.Create().SetName("[artist]manga 3.zip").SetCreateTime(time.UnixMilli(5000)).AddFavoriteOfUser(u).AddTags(tag2).Save(context.Background())
	client.Meta.Create().SetName("[artist]manga 4.zip").SetCreateTime(time.UnixMilli(6000)).AddFavoriteOfUser(u).AddTags(tag2).Save(context.Background())

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
	db, client, err := createTestDBClient(s)
	s.Assert().Nil(err)
	s.Assert().NotNil(db)
	s.Assert().NotNil(client)
	defer db.Close()
	defer client.Close()

	tag1, _ := client.Tag.Create().SetName("some artist").Save(context.Background())
	tag2, _ := client.Tag.Create().SetName("artist").Save(context.Background())

	u, err := user.GetUser(context.Background(), client, "")
	s.Assert().Nil(err)

	client.Meta.Create().SetName("[some artist]manga 1 here.zip").SetCreateTime(time.UnixMilli(3000)).AddFavoriteOfUser(u).AddTags(tag1).Save(context.Background())
	client.Meta.Create().SetName("[some artist]manga 2 here.zip").SetCreateTime(time.UnixMilli(4000)).AddFavoriteOfUser(u).AddTags(tag1).Save(context.Background())
	client.Meta.Create().SetName("[some artist]manga 3.zip").SetCreateTime(time.UnixMilli(5000)).AddFavoriteOfUser(u).AddTags(tag1).Save(context.Background())
	client.Meta.Create().SetName("[artist]manga 4.zip").SetCreateTime(time.UnixMilli(6000)).AddFavoriteOfUser(u).AddTags(tag2).Save(context.Background())

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
	db, client, err := createTestDBClient(s)
	s.Assert().Nil(err)
	s.Assert().NotNil(db)
	s.Assert().NotNil(client)
	defer db.Close()
	defer client.Close()

	client.Meta.Create().SetName("[some artist]manga 1 here.zip").SetFileIndices([]int{1}).Save(context.Background())
	client.Meta.Create().SetName("[some artist]manga 2 here.zip").SetFileIndices([]int{1, 2}).Save(context.Background())
	client.Meta.Create().SetName("[some artist]manga 3 here.zip").SetFileIndices([]int{1, 2, 3}).Save(context.Background())
	client.Meta.Create().SetName("[some artist]manga 4 here.zip").SetFileIndices([]int{1, 2, 3, 4}).Save(context.Background())
	client.Meta.Create().SetName("[some artist]manga 5 here.zip").SetFileIndices([]int{1, 2, 3, 4}).SetActive(false).Save(context.Background())

	u, err := user.GetUser(context.Background(), client, "")
	s.Assert().Nil(err)
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
	db, client, err := createTestDBClient(s)
	s.Assert().Nil(err)
	s.Assert().NotNil(db)
	s.Assert().NotNil(client)
	defer db.Close()
	defer client.Close()

	client.Meta.Create().SetName("[some artist]manga 1 here.zip").SetFileIndices([]int{1, 2, 3, 4}).Save(context.Background())
	client.Meta.Create().SetName("[some artist]manga 2 here.zip").SetFileIndices([]int{1, 2, 3}).Save(context.Background())
	client.Meta.Create().SetName("[some artist]manga 3 here.zip").SetFileIndices([]int{1, 2}).Save(context.Background())
	client.Meta.Create().SetName("[some artist]manga 4 here.zip").SetFileIndices([]int{1}).Save(context.Background())
	client.Meta.Create().SetName("[some artist]manga 5 here.zip").SetFileIndices([]int{1, 2, 3, 4}).SetActive(false).Save(context.Background())

	u, err := user.GetUser(context.Background(), client, "")
	s.Assert().Nil(err)
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
