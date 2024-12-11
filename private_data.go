package astiav

import "unsafe"

type PrivateData struct {
	c unsafe.Pointer
}

func newPrivateDataFromC(c unsafe.Pointer) *PrivateData {
	if c == nil {
		return nil
	}
	return &PrivateData{c: c}
}

func (pd *PrivateData) Options() *Options {
	return newOptionsFromC(pd.c)
}
