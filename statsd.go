package asyncstatsd

import (
	"net"
	"strconv"
	"time"

	"github.com/sirupsen/logrus"
)

type statsdclient struct {
	serverAddr net.Addr
	conn       net.PacketConn
}

// New returns a new statsdclient
func New(addr string) (Client, error) {
	conn, err := net.ListenPacket("udp", ":0")
	if err != nil {
		return nil, err
	}
	serverAddr, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		return nil, err
	}

	return &statsdclient{
		conn:       conn,
		serverAddr: serverAddr,
	}, nil
}

func (c *statsdclient) send(msg []byte) {
	n, err := c.conn.WriteTo(msg, c.serverAddr)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error":      err,
			"msg":        msg,
			"bytes_sent": n,
		}).Error("Could not send msg.")
	}
}

func (c *statsdclient) sendMetric(name, value, metricType string) {
	msg := name + ":" + value + "|" + metricType
	c.send([]byte(msg))
}

func (c *statsdclient) Count(bucket string, value int64) {
	c.sendMetric(bucket, strconv.FormatInt(value, 10), "c")
}

func (c *statsdclient) Gauge(bucket string, value int64) {
	c.sendMetric(bucket, strconv.FormatInt(value, 10), "g")
}

func (c *statsdclient) Increment(bucket string) {
	c.Count(bucket, 1)
}

func (c *statsdclient) Histogram(bucket string, value int64) {
	c.sendMetric(bucket, strconv.FormatInt(value, 10), "h")
}

func (c *statsdclient) Timing(bucket string, value int64) {
	c.sendMetric(bucket, strconv.FormatInt(value, 10), "h")
}

type statsdtiming struct {
	start time.Time
	c     Client
}

// Newstatsdtiming returns a new wrapped timing struct
func (c *statsdclient) NewTiming() Timing {
	return statsdtiming{
		start: time.Now(),
		c:     c,
	}
}

func (t statsdtiming) Send(bucket string) {
	t.c.Timing(bucket, int64(t.Duration()/time.Millisecond))
}

func (t statsdtiming) Duration() time.Duration {
	return time.Now().Sub(t.start)
}
