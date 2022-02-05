package astiav_test

import (
	"os"
	"testing"

	"github.com/asticode/go-astiav"
	"github.com/asticode/go-astikit"
)

var global = struct {
	closer             *astikit.Closer
	frame              *astiav.Frame
	inputFormatContext *astiav.FormatContext
	inputStream1       *astiav.Stream
	inputStream2       *astiav.Stream
	pkt                *astiav.Packet
}{
	closer: astikit.NewCloser(),
}

func TestMain(m *testing.M) {
	// Run
	m.Run()

	// Make sure to close closer
	global.closer.Close()

	// Exit
	os.Exit(0)
}
