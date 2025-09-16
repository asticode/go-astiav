package astiav

//#include <libavformat/avformat.h>
import "C"

// https://ffmpeg.org/doxygen/8.1/avformat_8h.html#af09f200b4cd9bf0baa05671436eef2fb
type DispositionFlag int64

const (
	DispositionFlagDefault         = DispositionFlag(C.AV_DISPOSITION_DEFAULT)
	DispositionFlagDub             = DispositionFlag(C.AV_DISPOSITION_DUB)
	DispositionFlagOriginal        = DispositionFlag(C.AV_DISPOSITION_ORIGINAL)
	DispositionFlagComment         = DispositionFlag(C.AV_DISPOSITION_COMMENT)
	DispositionFlagLyrics          = DispositionFlag(C.AV_DISPOSITION_LYRICS)
	DispositionFlagKaraoke         = DispositionFlag(C.AV_DISPOSITION_KARAOKE)
	DispositionFlagForced          = DispositionFlag(C.AV_DISPOSITION_FORCED)
	DispositionFlagHearingImpaired = DispositionFlag(C.AV_DISPOSITION_HEARING_IMPAIRED)
	DispositionFlagVisualImpaired  = DispositionFlag(C.AV_DISPOSITION_VISUAL_IMPAIRED)
	DispositionFlagCleanEffects    = DispositionFlag(C.AV_DISPOSITION_CLEAN_EFFECTS)
	DispositionFlagAttachedPic     = DispositionFlag(C.AV_DISPOSITION_ATTACHED_PIC)
	DispositionFlagTimedThumbnails = DispositionFlag(C.AV_DISPOSITION_TIMED_THUMBNAILS)
	DispositionFlagNonDiegetic     = DispositionFlag(C.AV_DISPOSITION_NON_DIEGETIC)
	DispositionFlagCaptions        = DispositionFlag(C.AV_DISPOSITION_CAPTIONS)
	DispositionFlagDescriptions    = DispositionFlag(C.AV_DISPOSITION_DESCRIPTIONS)
	DispositionFlagMetadata        = DispositionFlag(C.AV_DISPOSITION_METADATA)
	DispositionFlagDependent       = DispositionFlag(C.AV_DISPOSITION_DEPENDENT)
	DispositionFlagStillImage      = DispositionFlag(C.AV_DISPOSITION_STILL_IMAGE)
)
