package queue

import (
	"context"

	"github.com/redis/go-redis/v9"
)

type TaskQueue interface {
	Publish(ctx context.Context, jobID string) error
	Subscribe(ctx context.Context) <-chan string
}

type redisQueue struct {
	client  *redis.Client
	channel string
}

func NewRedisQueue(addr string, password string, db int) TaskQueue {
	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	return &redisQueue{
		client:  rdb,
		channel: "notification_jobs",
	}
}

// Publish sends a JOB ID to the Redis channel
func (r *redisQueue) Publish(ctx context.Context, jobID string) error {
	return r.client.Publish(ctx, r.channel, jobID).Err()
}

// Subscribe returns a Go channel that receives Job IDs as they arrive
func (r *redisQueue) Subscribe(ctx context.Context) <-chan string {
	out := make(chan string)
	pubsub := r.client.Subscribe(ctx, r.channel)

	go func() {
		defer pubsub.Close()
		defer close(out)

		ch := pubsub.Channel()
		for msg := range ch {
			out <- msg.Payload
		}
	}()

	return out
}
