package asyncstatsd

type Client interface {
	Count(bucket string, value int64)
	Gauge(bucket string, value int64)
	Increment(bucket string)
	Histogram(bucket string, value int64)
	Timing(bucket string, value int64)

	NewTiming() Timing
}
