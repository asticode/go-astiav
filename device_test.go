package astiav

import (
	"fmt"
	"testing"
)

func TestDevice(t *testing.T) {
	RegisterAllDevices()
}

func TestInputAudioDeviceNext(t *testing.T) {
	RegisterAllDevices()
	var d *InputFormat
	for {
		d = InputAudioDeviceNext(d)
		if d != nil {
			fmt.Println(d.Name(), d.LongName())
		} else {
			break
		}
	}
}

func TestOutputAudioDeviceNext(t *testing.T) {
	RegisterAllDevices()
	var d *OutputFormat
	for {
		d = OutputAudioDeviceNext(d)
		if d != nil {
			fmt.Println(d.Name(), d.LongName())
		} else {
			break
		}
	}
}

func TestInputVideoDeviceNext(t *testing.T) {
	RegisterAllDevices()
	var d *InputFormat
	for {
		d = InputVideoDeviceNext(d)
		if d != nil {
			fmt.Println(d.Name(), d.LongName())
		} else {
			break
		}
	}
}

func TestOutputVideoDeviceNext(t *testing.T) {
	RegisterAllDevices()
	var d *OutputFormat
	for {
		d = OutputVideoDeviceNext(d)
		if d != nil {
			fmt.Println(d.Name(), d.LongName())
		} else {
			break
		}
	}
}
