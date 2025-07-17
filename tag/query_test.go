package tag

import (
	"context"
	"database/sql"
	"testing"

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

func TestProviderTestSuite(t *testing.T) {
	suite.Run(t, new(QueryTestSuite))
}

func (s *QueryTestSuite) TestReadPage() {
	db, client, err := createTestDBClient(s)
	s.Assert().Nil(err)
	s.Assert().NotNil(db)
	s.Assert().NotNil(client)
	defer db.Close()
	defer client.Close()

	client.Tag.Create().SetName("Tag 1").Save(context.Background())
	client.Tag.Create().SetName("Tag 2").Save(context.Background())
	client.Tag.Create().SetName("Tag 3").Save(context.Background())

	u, err := user.GetUser(context.Background(), client, "")
	s.Assert().Nil(err)

	tags, err := ReadPage(context.Background(), client, u,
		QueryParams{
			Filter:      grpc.Filter_FILTER_UNKNOWN,
			Search:      "",
			Page:        0,
			ItemPerPage: 30,
		},
	)

	s.Assert().Nil(err)
	s.Assert().Equal(3, len(tags))
}

func createTestDBClient(s *QueryTestSuite) (db *sql.DB, client *ent.Client, err error) {
	db, err = sql.Open("sqlite", "file:ent?mode=memory&_fk=1&_pragma=foreign_keys(1)")
	if err != nil {
		return
	}

	client = enttest.NewClient(s.T(), enttest.WithOptions(ent.Driver(dialect_sql.OpenDB("sqlite3", db))))

	return
}

func (s *QueryTestSuite) TestReadPagePageCount() {
	db, client, err := createTestDBClient(s)
	s.Assert().Nil(err)
	s.Assert().NotNil(db)
	s.Assert().NotNil(client)
	defer db.Close()
	defer client.Close()

	client.Tag.Create().SetName("Tag 1").Save(context.Background())
	client.Tag.Create().SetName("Tag 2").Save(context.Background())
	client.Tag.Create().SetName("Tag 3").Save(context.Background())

	u, err := user.GetUser(context.Background(), client, "")
	s.Assert().Nil(err)

	tags, err := ReadPage(context.Background(), client, u,
		QueryParams{
			Filter:      grpc.Filter_FILTER_UNKNOWN,
			Search:      "",
			Page:        0,
			ItemPerPage: 2,
		})

	s.Assert().Nil(err)
	s.Assert().Equal(2, len(tags))

	tags, err = ReadPage(context.Background(), client, u,
		QueryParams{
			Filter:      grpc.Filter_FILTER_UNKNOWN,
			Search:      "",
			Page:        1,
			ItemPerPage: 2,
		})
	s.Assert().Nil(err)
	s.Assert().Equal(1, len(tags))
}

func (s *QueryTestSuite) TestReadPagePageWithSearch() {
	db, client, err := createTestDBClient(s)
	s.Assert().Nil(err)
	s.Assert().NotNil(db)
	s.Assert().NotNil(client)
	defer db.Close()
	defer client.Close()

	client.Tag.Create().SetName("Name 1").Save(context.Background())
	client.Tag.Create().SetName("Name 2").Save(context.Background())
	client.Tag.Create().SetName("Tag 3").Save(context.Background())

	u, err := user.GetUser(context.Background(), client, "")
	s.Assert().Nil(err)

	tags, err := ReadPage(context.Background(), client, u,
		QueryParams{
			Filter:      grpc.Filter_FILTER_UNKNOWN,
			Search:      "name",
			Page:        0,
			ItemPerPage: 30,
		})

	s.Assert().Nil(err)
	s.Assert().Equal(2, len(tags))
	s.Assert().Equal("Name 1", tags[0].Name)
	s.Assert().Equal("Name 2", tags[1].Name)
}

func (s *QueryTestSuite) TestReadPageWithSearchFilterFavoriteTags() {
	db, client, err := createTestDBClient(s)
	s.Assert().Nil(err)
	s.Assert().NotNil(db)
	s.Assert().NotNil(client)
	defer db.Close()
	defer client.Close()

	u, err := user.GetUser(context.Background(), client, "")
	s.Assert().Nil(err)

	_, err = client.Tag.Create().AddFavoriteOfUserIDs(u.ID).SetName("Name 1").Save(context.Background())
	s.Assert().Nil(err)
	_, err = client.Tag.Create().SetName("Name 2").Save(context.Background())
	s.Assert().Nil(err)
	_, err = client.Tag.Create().SetName("Tag 3").Save(context.Background())
	s.Assert().Nil(err)

	tags, err := ReadPage(context.Background(), client, u,
		QueryParams{
			Filter:      grpc.Filter_FILTER_FAVORITE_TAGS,
			Search:      "name",
			Page:        0,
			ItemPerPage: 30,
		})

	s.Assert().Nil(err)
	s.Assert().Equal(1, len(tags))
	s.Assert().Equal("Name 1", tags[0].Name)
}

func (s *QueryTestSuite) TestCount() {
	db, client, err := createTestDBClient(s)
	s.Assert().Nil(err)
	s.Assert().NotNil(db)
	s.Assert().NotNil(client)
	defer db.Close()
	defer client.Close()

	client.Tag.Create().SetName("Tag 1").Save(context.Background())
	client.Tag.Create().SetName("Tag 2").Save(context.Background())
	client.Tag.Create().SetName("Tag 3").Save(context.Background())

	u, err := user.GetUser(context.Background(), client, "")
	s.Assert().Nil(err)

	c, err := Count(context.Background(), client, u,
		QueryParams{
			Filter:      grpc.Filter_FILTER_UNKNOWN,
			Search:      "",
			Page:        0,
			ItemPerPage: 30,
		},
	)

	s.Assert().Nil(err)
	s.Assert().Equal(3, c)
}

func (s *QueryTestSuite) TestCountPageWithSearch() {
	db, client, err := createTestDBClient(s)
	s.Assert().Nil(err)
	s.Assert().NotNil(db)
	s.Assert().NotNil(client)
	defer db.Close()
	defer client.Close()

	client.Tag.Create().SetName("Name 1").Save(context.Background())
	client.Tag.Create().SetName("Name 2").Save(context.Background())
	client.Tag.Create().SetName("Tag 3").Save(context.Background())

	u, err := user.GetUser(context.Background(), client, "")
	s.Assert().Nil(err)

	c, err := Count(context.Background(), client, u,
		QueryParams{
			Filter:      grpc.Filter_FILTER_UNKNOWN,
			Search:      "name",
			Page:        0,
			ItemPerPage: 30,
		})

	s.Assert().Nil(err)
	s.Assert().Equal(2, c)
}

func (s *QueryTestSuite) TestCountWithSearchFilterFavoriteTags() {
	db, client, err := createTestDBClient(s)
	s.Assert().Nil(err)
	s.Assert().NotNil(db)
	s.Assert().NotNil(client)
	defer db.Close()
	defer client.Close()

	u, err := user.GetUser(context.Background(), client, "")
	s.Assert().Nil(err)

	_, err = client.Tag.Create().AddFavoriteOfUserIDs(u.ID).SetName("Name 1").Save(context.Background())
	s.Assert().Nil(err)

	client.Tag.Create().SetName("Name 2").Save(context.Background())
	client.Tag.Create().SetName("Tag 3").Save(context.Background())

	c, err := Count(context.Background(), client, u,
		QueryParams{
			Filter:      grpc.Filter_FILTER_FAVORITE_TAGS,
			Search:      "name",
			Page:        0,
			ItemPerPage: 30,
		})

	s.Assert().Nil(err)
	s.Assert().Equal(1, c)
}
