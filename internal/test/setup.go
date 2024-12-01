package testsuite

import "github.com/jmoiron/sqlx"

type MetaData map[string]interface{}

type TestSuite struct {
	MetaData MetaData // stores shared data for tests
	DB       *sqlx.DB // stores the database instance
}

func NewTestSuite() *TestSuite {
	newTestSuite := &TestSuite{
		MetaData: make(map[string]interface{}),
	}

	// reset database data, apply migrations and seed
	newTestSuite.InitTestEnv()

	return newTestSuite
}

func (t *TestSuite) InitTestEnv() {
	db := t.ConnectTestDB()

	// reset database data
	t.ResetTestDB(db)

	// update struct field DB with new, seeded DB instance
	t.DB = db

	// perform test migrations
	t.applyTestMigrations()

	// seed the database with test data
	t.seedTestData()
}
