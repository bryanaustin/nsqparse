package nsqparse

import (
	"testing"
)

func TestParsing(t *testing.T) {
	t.Parallel()
	if d, err := Parse("tcp://nsq.server:1234/coolthings/mine"); !checkerror(t, err) {
		compare(t, &Details{
			Scheme:"tcp",
			Address:"nsq.server:1234",
			Topic:"coolthings",
			Channel:"mine",
		}, d)
	}
}

func TestParsingSansChannel(t *testing.T) {
	t.Parallel()
	if d, err := Parse("nsqd://server:555/stuff"); !checkerror(t, err) {
		compare(t, &Details{
			Scheme:"nsqd",
			Address:"server:555",
			Topic:"stuff",
			Channel:"",
		}, d)
	}
}

func TestParsingDefaults(t *testing.T) {
	t.Parallel()
	if d, err := Parse("/solo"); !checkerror(t, err) {
		compare(t, &Details{
			Scheme:DefaultScheme,
			Address:DefaultAddress,
			Topic:"solo",
			Channel:"",
		}, d)
	}
}

func checkerror(t *testing.T, err error) bool {
	t.Helper()
	if err != nil {
		t.Errorf("Expected to not get an error, got %q", err)
	}
	return t.Failed()
}

func compare(t *testing.T, x, y *Details) {
	t.Helper()
	if x.Scheme != y.Scheme {
		t.Errorf("Expected scheme to be %q, it was %q", x.Scheme, y.Scheme)
	}
	if x.Address != y.Address {
		t.Errorf("Expected address to be %q, it was %q", x.Address, y.Address)
	}
	if x.Topic != y.Topic {
		t.Errorf("Expected topic to be %q, it was %q", x.Topic, y.Topic)
	}
	if x.Channel != y.Channel {
		t.Errorf("Expected channel to be %q, it was %q", x.Channel, y.Channel)
	}
}