#include <libavutil/intreadwrite.h>

uint32_t astiavRL32(uint8_t *i) {
	return AV_RL32(i);
}
uint32_t astiavRL32WithOffset(uint8_t *i, int o) {
	return AV_RL32(i+o);
}