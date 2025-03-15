package astiav

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func testFrameSideData(sd *FrameSideData, t *testing.T) {
	rois1 := []RegionOfInterest{
		{
			Bottom:             1,
			Left:               2,
			QuantisationOffset: NewRational(3, 4),
			Right:              5,
			Top:                6,
		},
		{
			Bottom:             7,
			Left:               8,
			QuantisationOffset: NewRational(9, 10),
			Right:              11,
			Top:                12,
		},
	}
	require.NoError(t, sd.RegionsOfInterest().Add(rois1))
	rois2, ok := sd.RegionsOfInterest().Get()
	require.True(t, ok)
	require.Equal(t, rois1, rois2)
}

func TestFrameSideData(t *testing.T) {
	f := AllocFrame()
	require.NotNil(t, f)
	defer f.Free()
	sd := f.SideData()

	rois1, ok := sd.RegionsOfInterest().Get()
	require.False(t, ok)
	require.Nil(t, rois1)
	rois1 = []RegionOfInterest{
		{
			Bottom:             1,
			Left:               2,
			QuantisationOffset: NewRational(3, 4),
			Right:              5,
			Top:                6,
		},
		{
			Bottom:             7,
			Left:               8,
			QuantisationOffset: NewRational(9, 10),
			Right:              11,
			Top:                12,
		},
	}
	require.NoError(t, sd.RegionsOfInterest().Add(rois1))
	rois2, ok := sd.RegionsOfInterest().Get()
	require.True(t, ok)
	require.Equal(t, rois1, rois2)
}
