package astiav

//#include <libavdevice/avdevice.h>
import "C"

func RegisterAllDevices() {
	C.avdevice_register_all()
}
