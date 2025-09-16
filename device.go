package astiav

//#include <libavdevice/avdevice.h>
import "C"

// https://ffmpeg.org/doxygen/8.1/group__lavd.html#ga7c90a3585267b55941ae2f7388c006b6
func RegisterAllDevices() {
	C.avdevice_register_all()
}
