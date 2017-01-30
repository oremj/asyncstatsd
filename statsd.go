package asyncstatsd

import (
	"time"

	"github.com/alexcesaro/statsd"
)

type Client struct {
	statsd *statsd.Client

	queue *RunQueue
}

// New returns a new client
func New(c *statsd.Client, queueSize int) *Client {
	return &Client{
		statsd: c,
		queue:  NewRunQueue(queueSize),
	}
}

func (c *Client) Count(bucket string, n interface{}) {
	c.queue.Queue(func() {
		c.statsd.Count(bucket, n)
	})
}

func (c *Client) Gauge(bucket string, value interface{}) {
	c.queue.Queue(func() {
		c.statsd.Gauge(bucket, value)
	})
}

func (c *Client) Increment(bucket string) {
	c.queue.Queue(func() {
		c.statsd.Increment(bucket)
	})
}

func (c *Client) Histogram(bucket string, value interface{}) {
	c.queue.Queue(func() {
		c.statsd.Histogram(bucket, value)
	})
}

func (c *Client) Timing(bucket string, value interface{}) {
	c.queue.Queue(func() {
		c.statsd.Timing(bucket, value)
	})
}

func (c *Client) Clone(opts ...statsd.Option) *Client {
	c.statsd.Clone()
	return &Client{
		statsd: c.statsd.Clone(opts...),
		queue:  c.queue,
	}
}

type Timing struct {
	timing statsd.Timing
	c      *Client
}

// NewTiming returns a new wrapped timing struct
func (c *Client) NewTiming() Timing {
	return Timing{
		timing: c.statsd.NewTiming(),
		c:      c,
	}
}

func (t Timing) Send(bucket string) {
	t.c.Timing(bucket, int(t.Duration()/time.Millisecond))
}

func (t Timing) Duration() time.Duration {
	return t.timing.Duration()
}
