package astiav

import (
	"encoding/binary"
	"fmt"
)

// https://ffmpeg.org/doxygen/8.1/group__lavc__packet__side__data.html#gga9a80bfcacc586b483a973272800edb97a2093332d8086d25a04942ede61007f6a
// https://ffmpeg.org/doxygen/8.1/group__lavu__frame.html#ggae01fa7e427274293aacdf2adc17076bca6b0b1ee4315f322922710f65d02a146b
type SkipSamples struct {
	ReasonEnd   uint8
	ReasonStart uint8
	SkipEnd     uint32
	SkipStart   uint32
}

func newSkipSamplesFromBytes(b []byte) (*SkipSamples, error) {
	if len(b) < 10 {
		return nil, fmt.Errorf("astiav: invalid length %d < 10", len(b))
	}
	return &SkipSamples{
		ReasonEnd:   b[9],
		ReasonStart: b[8],
		SkipEnd:     binary.LittleEndian.Uint32(b[4:8]),
		SkipStart:   binary.LittleEndian.Uint32(b[0:4]),
	}, nil
}

func (ss *SkipSamples) bytes() (b []byte) {
	b = binary.LittleEndian.AppendUint32(b, ss.SkipStart)
	b = binary.LittleEndian.AppendUint32(b, ss.SkipEnd)
	b = append(b, ss.ReasonStart, ss.ReasonEnd)
	return b
}
