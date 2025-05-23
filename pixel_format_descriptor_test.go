package astiav

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPixelFormatDescriptor(t *testing.T) {
	p := PixelFormatCuda
	d := p.Descriptor()
	require.NotNil(t, d)
	require.Equal(t, d.Name(), p.String())
	require.True(t, d.Flags().Has(PixelFormatDescriptorFlagHwAccel))
}
