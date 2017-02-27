package bolt

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/apiarian/ipfs-pinbase/pinbase/test"
)

func TestClientService(t *testing.T) {
	filename := tempfilename(t)
	defer os.Remove(filename)

	c := NewClient(filename)
	err := c.Open()
	if err != nil {
		t.Fatalf("failed to open client: %+v", err)
	}

	ps := c.PinService()

	test.TestPinServiceHappyPath(t, ps)
}

func TestClientBackend(t *testing.T) {
	filename := tempfilename(t)
	defer os.Remove(filename)

	c := NewClient(filename)
	err := c.Open()
	if err != nil {
		t.Fatalf("failed to open client: %+v", err)
	}

	ps := c.PinService()
	pb := c.PinBackend()

	test.TestPinBackendHappyPath(t, pb, ps)
}

func TestClientFeedback(t *testing.T) {
	filename := tempfilename(t)
	defer os.Remove(filename)

	c := NewClient(filename)
	err := c.Open()
	if err != nil {
		t.Fatalf("failed to open clien: %+v", err)
	}

	ps := c.PinService()
	pb := c.PinBackend()

	test.TestPinFeedbackHappyPath(t, pb, ps)
}

func tempfilename(t *testing.T) string {
	f, err := ioutil.TempFile("", "pinbase-bolt-")
	if err != nil {
		t.Fatalf("failed to create temp file: %+v", err)
	}

	err = f.Close()
	if err != nil {
		t.Fatalf("failed to close temp file: %+v", err)
	}

	err = os.Remove(f.Name())
	if err != nil {
		t.Fatalf("failed to remove temp file: %+v", err)
	}

	return f.Name()
}
