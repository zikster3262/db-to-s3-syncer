package syncer

import (
	"bytes"
	"concurrency/src/db"
	"concurrency/src/models"
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
	"golang.org/x/sync/errgroup"
)

var (
	ErrNotCantRetriveData = errors.New("can't retrive data from database")
)

type Syncer struct {
	s3w *s3.Client
	db  *sqlx.DB
}

func NewSyncer(db *sqlx.DB, s3cliet *s3.Client) Syncer {
	return Syncer{
		s3w: s3cliet,
		db:  db,
	}
}

func (s *Syncer) Sync(ctx context.Context) error {
	ticker := time.NewTicker(time.Second * 5)
	defer ticker.Stop()

	for {
		log.Info().Msg("Running Syncer.......")
		requests, err := db.GetAllRequest(s.db)
		if err != nil {
			log.Error().Msg("cant select request data")
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
		log.Error().Msg("couldn't list bucket contents")
	}

	for _, object := range listObjsResponse.Contents {
		if *object.Key == uuid {
			return true
		}
	}

	return false

}

func (s *Syncer) PutS3Object(m models.DbRequest) error {
	var bt []byte
	var err error
	bt, err = json.Marshal(m)
	if err != nil {
		log.Error().Msg("couldn't unmarshall the request")
	}

	body := bytes.NewReader(bt)

	_, err = s.s3w.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String("requests"),
		Key:    aws.String(m.Uuid),
		Body:   body,
	})

	if err != nil {
		panic("Couldn't upload file: " + err.Error())
	}
	return err
}
