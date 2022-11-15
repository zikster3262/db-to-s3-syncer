package syncer

import (
	"bytes"
	"concurrency/src/db"
	"concurrency/src/models"
	"concurrency/src/utils"
	"context"
	"errors"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
	"github.com/streadway/amqp"
	"golang.org/x/sync/errgroup"
)

var (
	ErrNotCantRetriveData = errors.New("can't retrive data from database")
)

type Syncer struct {
	s3w *s3.Client
	db  *sqlx.DB
	rmq *amqp.Channel
	q   amqp.Queue
}

func NewSyncer(db *sqlx.DB, s3cliet *s3.Client, rmq *amqp.Channel, q amqp.Queue) Syncer {
	return Syncer{
		s3w: s3cliet,
		db:  db,
		rmq: rmq,
		q:   q,
	}
}

func (s *Syncer) Sync(ctx context.Context) error {
	ticker := time.NewTicker(time.Second * 5)
	defer ticker.Stop()

	for {
		log.Info().Msg("Running Syncer.......")
		requests, err := db.GetAllRequest(s.db)
		if err != nil {
			utils.FailOnError(err, "cant select request data")
		}

		for i := 0; i < len(requests); i++ {
			ex := s.GetS3Object(requests[i].Uuid)
			if !ex {
				s.PutS3Object(requests[i])
			}
		}

		select {
		case <-ticker.C:
		case <-ctx.Done():
			return ctx.Err()
		}

	}
}

func (s Syncer) Run(ctx context.Context) error {
	g, c := errgroup.WithContext(ctx)
	g.Go(func() error { return s.Sync(c) })
	return g.Wait()
}

func (s *Syncer) GetS3Object(uuid string) bool {

	listObjsResponse, err := s.s3w.ListObjectsV2(context.TODO(), &s3.ListObjectsV2Input{
		Bucket: aws.String("requests"),
	})

	if err != nil {
		utils.FailOnError(err, "couldn't list bucket contents")
	}

	for _, object := range listObjsResponse.Contents {
		if *object.Key == uuid {
			return true
		}
	}

	return false

}

func (s *Syncer) PutS3Object(m models.DbRequest) error {
	bt := utils.StructToJson(m)
	body := bytes.NewReader(bt)

	_, err := s.s3w.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String("requests"),
		Key:    aws.String(m.Uuid),
		Body:   body,
	})

	if err != nil {
		utils.FailOnError(err, "Couldn't upload file: "+err.Error())
	}
	return err
}
