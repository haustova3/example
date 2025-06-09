package repository

import (
	"context"
	"example/comments/internal/model"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
	"log"
	"strings"
	"testing"
	"time"
)

const (
	dbUser     = "comments-user-1"
	dbPassword = "comments-password-1"
	dbName     = "comments-db"
)

type RepositoryIntegrationTestSuite struct {
	suite.Suite
	repositoryContainer *postgres.PostgresContainer
	repository          *Repository
	rwPool              *pgxpool.Pool
}

func TestIntegration(t *testing.T) {
	s := new(RepositoryIntegrationTestSuite)
	suite.Run(t, s)
}

func (s *RepositoryIntegrationTestSuite) SetupSuite() {
	ctx := context.Background()
	container, err := postgres.Run(ctx,
		"postgres:16",
		postgres.WithDatabase(dbName),
		postgres.WithUsername(dbUser),
		postgres.WithPassword(dbPassword),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).WithStartupTimeout(time.Minute)),
	)
	s.Suite.Require().NoError(err, "create container postgres")
	s.repositoryContainer = container
	postgresAddress, err := s.repositoryContainer.ConnectionString(ctx, "sslmode=disable")
	s.Suite.Require().NoError(err, "can not get postgres host")
	postgresConf, err := pgxpool.ParseConfig(postgresAddress)
	if err != nil {
		log.Fatalf("unable to parse master repository config: %v\n", err)
	}
	s.rwPool, err = pgxpool.NewWithConfig(ctx, postgresConf)
	if err != nil {
		log.Fatalf("unable to create pgx pool master: %v\n", err)
	}

	db := stdlib.OpenDBFromPool(s.rwPool)
	if err := goose.Up(db, "../../migrations"); err != nil {
		log.Fatalf("Failed to apply migrations: %v", err)
	}

	s.repository = NewRepository(s.rwPool)
}

func (s *RepositoryIntegrationTestSuite) TestSaveAndGetCommentSuccess() {
	ctx := context.Background()
	com := model.Comment{
		UserID:         456,
		ProductID:      123,
		ProductOwnerID: 789,
		Text:           "Отличный товар",
	}
	comID, err := s.repository.SaveComment(ctx, com)
	s.Suite.Require().NoError(err, "Can not save comment")
	comments, err := s.repository.GetComments(ctx, 123)
	s.Suite.Require().NoError(err, "Can not get comments")
	s.Suite.Require().Equal(1, len(comments), "Len comments mismatch")
	s.Suite.Require().Equal(comID, comments[0].ID, "Comment ID mismatch")
	s.Suite.Require().Equal(com.UserID, comments[0].UserID, "UserID mismatch")
	s.Suite.Require().Equal(com.ProductID, comments[0].ProductID, "ProductID mismatch")
	s.Suite.Require().Equal(com.Text, comments[0].Text, "Text mismatch")
	ntfs, err := s.repository.GetCommentNotification(ctx)
	s.Suite.Require().NoError(err, "Can not get notifications")
	s.Suite.Require().Equal(1, len(ntfs), "Len Notifications mismatch(1)")
	s.Suite.Require().Equal(comID, ntfs[0].CommentID, "Notification Comment ID mismatch")
	s.Suite.Require().Equal(int64(789), ntfs[0].OwnerID, "Notification OwnerID mismatch")
	err = s.repository.MarkNotificationAsSend(ctx, ntfs[0].ID)
	s.Suite.Require().NoError(err, "Can not mark notification as send")
	ntfs, err = s.repository.GetCommentNotification(ctx)
	s.Suite.Require().NoError(err, "Can not get notifications after sent")
	s.Suite.Require().Equal(0, len(ntfs), "Len notifications mismatch(0)")
}

func (s *RepositoryIntegrationTestSuite) TestSaveCommentFailedIDNotPositive() {
	ctx := context.Background()
	com := model.Comment{
		UserID:         0,
		ProductID:      123,
		ProductOwnerID: 789,
		Text:           "Отличный товар",
	}
	_, err := s.repository.SaveComment(ctx, com)
	s.Suite.Require().Error(err, "Saved zero userID")
	com.UserID = 456
	com.ProductOwnerID = 0
	_, err = s.repository.SaveComment(ctx, com)
	s.Suite.Require().Error(err, "Saved zero ProductOwnerID")
	com.ProductOwnerID = 789
	com.ProductID = 0
	_, err = s.repository.SaveComment(ctx, com)
	s.Suite.Require().Error(err, "Saved zero ProductID")
}

func (s *RepositoryIntegrationTestSuite) TestGetEmptyCommentList() {
	ctx := context.Background()
	comments, err := s.repository.GetComments(ctx, 321)
	s.Suite.Require().NoError(err, "Can not get comments by productID")
	s.Suite.Require().Equal(0, len(comments), "Len comments mismatch")
}

func (s *RepositoryIntegrationTestSuite) TestSaveCommentFailedWrongTextLen() {
	ctx := context.Background()
	com := model.Comment{
		UserID:         456,
		ProductID:      123,
		ProductOwnerID: 789,
		Text:           "Отл",
	}
	_, err := s.repository.SaveComment(ctx, com)
	s.Suite.Require().Error(err, "Len text < 5")
	strBuilder := strings.Builder{}
	for i := 0; i < 512; i++ {
		strBuilder.WriteRune('s')
	}
	com.Text = strBuilder.String()
	_, err = s.repository.SaveComment(ctx, com)
	s.Suite.Require().Error(err, "len text > 256")
}
