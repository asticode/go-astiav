#include <libavutil/log.h>
#include <stdint.h>
#include <stdlib.h>

char* astiavClassItemName(AVClass* c, void* ptr) {
	return (char*)c->item_name(ptr);
}

AVClassCategory astiavClassCategory(AVClass* c, void* ptr) {
	if (c->get_category) return c->get_category(ptr);
	return c->category;
}

AVClass** astiavClassParent(AVClass* c, void* ptr) {
	if (c->parent_log_context_offset) {
		AVClass** parent = *(AVClass ***) (((uint8_t *) ptr) + c->parent_log_context_offset);
		if (parent && *parent) {
			return parent;
		}
	}
	return NULL;
}