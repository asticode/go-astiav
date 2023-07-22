package astiav

//#cgo pkg-config: libavdevice
//#include <libavdevice/avdevice.h>
import "C"

func RegisterAllDevices() {
	C.avdevice_register_all()
}
