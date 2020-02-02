package server

import (
	"os"
	"testing"

	"github.com/apex/log"
)

func TestMain(m *testing.M) {
	log.SetLevel(log.FatalLevel)

	os.Exit(m.Run())
}
