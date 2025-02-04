package indexing

import (
	"flag"
	"os"
	"testing"

	"github.com/inconshreveable/log15"
)

func init() {
	autoIndexingEnabled = func() bool { return true }
}

func TestMain(m *testing.M) {
	flag.Parse()
	if !testing.Verbose() {
		log15.Root().SetHandler(log15.DiscardHandler())
	}
	os.Exit(m.Run())
}
