package main

import (
	"concurrency/src/api"
	"concurrency/src/awss3"
	"concurrency/src/db"
	"concurrency/src/rabbitmq"
	"concurrency/src/runner"
	"concurrency/src/syncer"
	"concurrency/src/utils"
	"context"
	"net/http"
	"os"
	"syscall"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/jmoiron/sqlx"
)

var (
	sqlxDB   *sqlx.DB
	s3Client *s3.Client
	sy       syncer.Syncer
	queue    = "go-rabbit"
)

func main() {

	Initialize()
}

func Initialize() error {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	rabbitCh, err := rabbitmq.ConnectToRabbit()
	utils.FailOnError(err, "cant connect to RabbitMQ")

	q, err := rabbitmq.CreateRabbitMQueue(rabbitCh, queue)
	utils.FailOnError(err, "cant create RabbitMQ queue")

	s3Client = awss3.SetS3Config()
	sqlxDB = db.OpenSQLx()
	sy = syncer.NewSyncer(sqlxDB, s3Client, rabbitCh, q)
	router, err := NewServer(ctx)
	if err != nil {
		utils.FailOnError(err, err.Error())
	}
	_ = awss3.CreateBucket(s3Client, "requests")

	runners := []runner.Runner{
		runner.NewSignal(os.Interrupt, syscall.SIGTERM),
		router,
		sy,
	}

	err = runner.RunParallel(ctx, runners...)
	switch err {
	case context.Canceled, runner.SignalReceived, nil:
	default:
		return err
	}
	return nil
}

func NewServer(ctx context.Context) (runner.Runner, error) {

	rt := api.Register(sqlxDB, ctx)

	return runner.NewServer(
		&http.Server{
			Handler:      rt,
			Addr:         ":8080",
			ReadTimeout:  time.Second * 20,
			WriteTimeout: time.Second * 20,
		}, time.Second*10,
	), nil

}
