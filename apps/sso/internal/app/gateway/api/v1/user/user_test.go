package user

import (
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/suite"

	"github.com/kitanoyoru/kitaDriveBot/apps/sso/internal/internal/migrator"
	"github.com/kitanoyoru/kitaDriveBot/libs/test"
)

func TestIntegrationUserServiceSuite(t *testing.T) {
	testSuite := &integrationUserServiceSuite{}
	suite.Run(t, testSuite)
}

type integrationUserServiceSuite struct {
	suite.Suite
	test.DBTest

	dbPool *pgxpool.Pool
}

func (s *integrationUserServiceSuite) SetupTest() {
	var err error
	s.dbPool, err = s.CreateDB(s.T(), test.WithMigrator(migrator.Migrator))
	s.Require().NoError(err)

	s.Require().NoError(err)
}

func (s *integrationUserServiceSuite) TearDownTest() {
	if s.dbPool != nil {
		s.dbPool.Close()
	}
}
