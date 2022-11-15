package uploader

import (
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/jmoiron/sqlx"
)

type s3Client struct {
	client         *s3.Client
	SQLXConnection *sqlx.DB
}

func NewS3Client(c *s3.Client, sqlConn *sqlx.DB) s3Client {
	return s3Client{
		client:         c,
		SQLXConnection: sqlConn,
	}
}

// func (client s3Client) Run(ctx context.Context) error {
// 	g, c := errgroup.WithContext(ctx)
// 	g.Go(func() error { return s.Sync(c) })
// 	return g.Wait()
// }
