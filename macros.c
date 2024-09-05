#include <libavutil/macros.h>
#include <stddef.h>

ptrdiff_t astiavFFAlign(int i, int align)
{
	return FFALIGN(i, align);
}