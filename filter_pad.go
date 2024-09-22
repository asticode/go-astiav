package astiav

// Struct attributes are internal but there are C functions to get some of them
type FilterPad struct {
	mediaType MediaType
}

func newFilterPad(mediaType MediaType) *FilterPad {
	return &FilterPad{mediaType: mediaType}
}

func (fp *FilterPad) MediaType() MediaType {
	return fp.mediaType
}
