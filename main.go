/*
Parse strings for NSQ connections.
*/
package nsqparse

import (
	"errors"
	"fmt"
	"github.com/nsqio/go-nsq"
	"math/rand"
	"net/url"
	"strings"
	"time"
)

const (
	DefaultScheme  = "tcp"
	DefaultAddress = "localhost:4150"
)

var (
	NoTopic = errors.New("no topic")
	chars   = []rune("abcdefghijklmnopqrstuvwxyz0123456789")
)

type Details struct {
	Scheme  string
	Address string
	Topic   string
	Channel string
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

// Parse URL
func Parse(addr string) (d *Details, err error) {
	d, err = ParseNoDefaults(addr)
	if err != nil {
		return
	}

	if len(d.Scheme) < 1 {
		d.Scheme = DefaultScheme
	}

	if len(d.Address) < 1 {
		d.Address = DefaultAddress
	}

	return
}

// ParseStrict throws an error if topic is not provided
func ParseStrict(addr string) (d *Details, err error) {
	d, err = Parse(addr)
	if err != nil {
		return
	}
	if len(d.Topic) < 1 {
		err = NoTopic
	}
	return
}

// ParseNoDefaults will parse and not fill in missing fields with defaults
func ParseNoDefaults(addr string) (d *Details, err error) {
	var u *url.URL
	d = new(Details)

	u, err = url.Parse(addr)
	if err != nil {
		err = fmt.Errorf("parsing string: %w", err)
		return
	}

	d.Scheme = u.Scheme
	d.Address = u.Host
	tp := strings.Trim(u.Path, "/")
	ps := strings.Split(tp, "/")

	if len(ps) > 0 {
		d.Topic = ps[0]
		if len(ps) > 1 {
			d.Channel = ps[1]
		}
	}

	return
}

func randword(n int) string {
	r := make([]rune, n)
	for i := range r {
		r[i] = chars[rand.Intn(len(chars))]
	}
	return string(r)
}

// Consumer will make an NSQ consumer with your details
func (d *Details) Consumer(c *nsq.Config) (*nsq.Consumer, error) {
	channel := d.Channel
	if len(channel) < 1 {
		channel = randword(8) + "#ephemeral"
	}
	return nsq.NewConsumer(d.Topic, channel, c)
}

//TODO: Support for multipe addresses

// ConnectConsumer will connect this consumer using provided details
func (d *Details) ConnectConsumer(c *nsq.Consumer) error {
	//TODO: Add lookupd support
	return c.ConnectToNSQD(d.Address)
}

// Producer will make a NSQ producer with the provided details
func (d *Details) Producer(c *nsq.Config) (*nsq.Producer, error) {
	return nsq.NewProducer(d.Address, c)
}
