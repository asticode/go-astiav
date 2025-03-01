#include <stdatomic.h>

int astiavAtomicLoadInt(atomic_int* i)
{
    return atomic_load(i);
}

void astiavAtomicStoreInt(atomic_int* i, int v)
{
    return atomic_store(i, v);
}