package astiav

//#cgo pkg-config: libavdevice
//#include <libavdevice/avdevice.h>
import "C"

func DeviceRegisterAll() {
	C.avdevice_register_all()
}
