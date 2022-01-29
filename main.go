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
	"net"
	"strings"
	"time"
)

const (
	DefaultScheme  = "tcp"
	DefaultAddress = "localhost:4150"
	DefaultPort = ":4150"
	MagicErrorWordsOfPortMissing = "missing port in address"
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
	} else {
		_, _, err := net.SplitHostPort(d.Address)
		if err != nil && strings.Contains(err.Error(), MagicErrorWordsOfPortMissing) {
			d.Address = d.Address + DefaultPort
		}
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


	if len(u.Host) < 1 {
		// Host not parsed
		if len(u.Opaque) > 0 {
			// Uses opaque syntax
			tp := strings.TrimRight(u.Opaque, "/")
			ps := strings.Split(tp, "/")

			d.Address = u.Scheme + ":" + ps[0]
			d.Scheme = ""

			if len(ps) > 1 {
				d.Topic = ps[1]
				if len(ps) > 2 {
					d.Channel = ps[2]
				}
			}
		} else {
			// Probably means the first part of the path is the host
			d.Scheme = u.Scheme
			tp := strings.TrimRight(u.Path, "/")
			ps := strings.Split(tp, "/")

			d.Address = ps[0]
			if len(ps) > 1 {
				d.Topic = ps[1]
				if len(ps) > 2 {
					d.Channel = ps[2]
				}
			}
		}
	} else {
		d.Scheme = u.Scheme
		d.Address = u.Host
		tp := strings.Trim(u.Path, "/")
		ps := strings.Split(tp, "/")

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

//TODO: Support for multiple addresses

// ConnectConsumer will connect this consumer using provided details
func (d *Details) ConnectConsumer(c *nsq.Consumer) error {
	//TODO: Add lookupd support
	return c.ConnectToNSQD(d.Address)
}

// Producer will make a NSQ producer with the provided details
func (d *Details) Producer(c *nsq.Config) (*nsq.Producer, error) {
	return nsq.NewProducer(d.Address, c)
}
