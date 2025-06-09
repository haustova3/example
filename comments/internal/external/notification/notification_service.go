package notification

import (
	"context"
	"encoding/json"
	"example/comments/internal/logger"
	"fmt"
	"strconv"
	"time"

	"github.com/IBM/sarama"
)

type CommentNotificationRepository interface {
	GetCommentNotification(_ context.Context) ([]CommentNotification, error)
	MarkNotificationAsSend(_ context.Context, notificationID int64) error
}

type OrderNotificationService struct {
	rep      CommentNotificationRepository
	updateCh chan int64
	ticker   *time.Ticker
	topic    string
	prod     sarama.SyncProducer
}

func StartNotificationService(ctx context.Context,
	rep CommentNotificationRepository,
	brokers []string,
	topic string,
	maxCount int,
	timer int) {
	orderService := &OrderNotificationService{
		rep:      rep,
		updateCh: make(chan int64, maxCount),
		ticker:   time.NewTicker(time.Duration(timer) * time.Millisecond),
		topic:    topic,
	}
	var err error
	orderService.prod, err = orderService.newSyncProducer(brokers)
	if err != nil {
		logger.Warnw(ctx, "create producer failed", "error", err.Error())
		return
	}

	go func(s *OrderNotificationService) {
		for {
			select {
			case <-ctx.Done():
				logger.Infow(ctx, "notification service context closed")
				close(s.updateCh)
				s.prod.Close()
				s.ticker.Stop()
				return
			case <-s.ticker.C:
				go func(appCtx context.Context) {
					ctx, cancel := context.WithTimeout(appCtx, 100*time.Millisecond)
					s.SendNotification(ctx)
					cancel()
				}(ctx)
			case notificationID := <-s.updateCh:
				ctx, cancel := context.WithTimeout(ctx, 100*time.Millisecond)
				s.MarkNotificationAsSend(ctx, notificationID)
				cancel()
			}
		}
	}(orderService)
}

func (s *OrderNotificationService) SendNotification(ctx context.Context) {
	if ctx.Err() != nil {
		return
	}

	ntfs, err := s.rep.GetCommentNotification(ctx)
	if err != nil {
		logger.Warnw(ctx, "can not get notifications", "error", err.Error())
		return
	}
	if len(ntfs) == 0 {
		return
	}
	for _, val := range ntfs {
		if ctx.Err() != nil {
			return
		}
		bytes, err := json.Marshal(val)
		if err != nil {
			logger.Warnw(ctx, "marshal notification failed", "error", err.Error())
			break
		}
		msg := &sarama.ProducerMessage{
			Topic:     s.topic,
			Key:       sarama.StringEncoder(strconv.FormatInt(val.ID, 10)),
			Value:     sarama.ByteEncoder(bytes),
			Timestamp: time.Now(),
		}
		partition, offset, err := s.prod.SendMessage(msg)
		if err != nil {
			logger.Warnw(ctx, "can not send notification", "error", err.Error())
			continue
		}
		s.updateCh <- val.ID
		logger.Infow(ctx, "send notification: new comment",
			"key", val.ID,
			"partition", partition,
			"offset", offset,
			"owner_id", val.OwnerID,
			"user_id", val.CommentID)

	}
}

func (s *OrderNotificationService) newSyncProducer(brokers []string) (sarama.SyncProducer, error) {
	config := sarama.NewConfig()
	config.Producer.Partitioner = sarama.NewHashPartitioner
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Idempotent = false
	config.Producer.Retry.Max = 10
	config.Producer.Retry.Backoff = 5 * time.Millisecond
	config.Net.MaxOpenRequests = 1
	config.Producer.CompressionLevel = sarama.CompressionLevelDefault
	config.Producer.Compression = sarama.CompressionGZIP
	config.Metadata.AllowAutoTopicCreation = false
	config.Producer.Return.Successes = true
	config.Producer.Return.Errors = true
	syncProducer, err := sarama.NewSyncProducer(brokers, config)
	if err != nil {
		return nil, fmt.Errorf("NewSyncProducer failed: %w", err)
	}
	return syncProducer, nil
}

func (s *OrderNotificationService) MarkNotificationAsSend(ctx context.Context, notificationID int64) {
	err := s.rep.MarkNotificationAsSend(ctx, notificationID)
	if err != nil {
		logger.Warnw(ctx, "can not mark notification as send", "err", err.Error())
	}
}
