package asyncstatsd

import (
	"time"
)

type noopclient struct{}

// NewNoop returns a new noop client
func NewNoop() Client {
	return new(noopclient)
}

func (c *noopclient) Count(bucket string, n int64) {
}

func (c *noopclient) Gauge(bucket string, value int64) {
}

func (c *noopclient) Increment(bucket string) {
}

func (c *noopclient) Histogram(bucket string, value int64) {
}

func (c *noopclient) Timing(bucket string, value int64) {
}

type nooptiming struct{}

func (c *noopclient) NewTiming() Timing {
	return nooptiming{}
}

func (t nooptiming) Send(bucket string) {
}

func (t nooptiming) Duration() time.Duration {
	return 0
}
