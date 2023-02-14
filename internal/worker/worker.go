package worker

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gocraft/work"
	"github.com/gomodule/redigo/redis"
	"github.com/w1lldone/go-khotmil-bot/internal/models"
)

type Context struct {
}

var redisPool = &redis.Pool{
	MaxActive: 5,
	MaxIdle:   5,
	Wait:      true,
	Dial: func() (redis.Conn, error) {
		conn, err := redis.Dial("tcp", os.Getenv("REDIS_HOST"), redis.DialPassword(os.Getenv("REDIS_PASSWORD")))
		if err != nil {
			log.Fatal("error connecting to Redis ", err)
			return nil, err
		}

		return conn, nil
	},
}

func InitWorker() {
	var pool = work.NewWorkerPool(Context{}, 1, "go_khotmil_bot", redisPool)

	pool.Middleware((*Context).Log)

	pool.PeriodicallyEnqueue("@every 1m", "reminder")
	pool.JobWithOptions("reminder", work.JobOptions{MaxConcurrency: 1}, (*Context).SendReminder)

	fmt.Println("Starting worker pool.")
	pool.Start()
	// Wait for a signal to quit:
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGINT)
	<-signalChan

	// Stop the pool
	fmt.Println("Stopping worker pool.")
	pool.Stop()

}

func (c *Context) SendReminder(job *work.Job) error {
	var schedules []models.Schedule
	err := models.DB.Where("finished_at IS NULL").Where("deadline <= ?", time.Now()).Preload("Member").Find(&schedules).Error
	if err != nil {
		return err
	}

	for _, s := range schedules {
		fmt.Printf("Found unfinished schedule for: %s \n", s.Member.Name)
	}

	fmt.Print(time.Now().Format(time.RFC3339Nano), " ")
	fmt.Printf("Finished executing job: %s \n", job.Name)
	return nil
}

func (c *Context) Log(job *work.Job, next work.NextMiddlewareFunc) error {
	fmt.Print(time.Now().Format(time.RFC3339Nano), " ")
	fmt.Printf("Starting job: %s \n", job.Name)
	return next()
}
