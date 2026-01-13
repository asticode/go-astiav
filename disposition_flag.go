package astiav

//#include <libavformat/avformat.h>
import "C"

// https://ffmpeg.org/doxygen/8.0/avformat_8h.html#af09f200b4cd9bf0baa05671436eef2fb
type DispositionFlag int64

const (
	DispositionFlagAttachedPic     = DispositionFlag(C.AV_DISPOSITION_ATTACHED_PIC)
	DispositionFlagCaptions        = DispositionFlag(C.AV_DISPOSITION_CAPTIONS)
	DispositionFlagCleanEffects    = DispositionFlag(C.AV_DISPOSITION_CLEAN_EFFECTS)
	DispositionFlagComment         = DispositionFlag(C.AV_DISPOSITION_COMMENT)
	DispositionFlagDefault         = DispositionFlag(C.AV_DISPOSITION_DEFAULT)
	DispositionFlagDependent       = DispositionFlag(C.AV_DISPOSITION_DEPENDENT)
	DispositionFlagDescriptions    = DispositionFlag(C.AV_DISPOSITION_DESCRIPTIONS)
	DispositionFlagDub             = DispositionFlag(C.AV_DISPOSITION_DUB)
	DispositionFlagForced          = DispositionFlag(C.AV_DISPOSITION_FORCED)
	DispositionFlagHearingImpaired = DispositionFlag(C.AV_DISPOSITION_HEARING_IMPAIRED)
	DispositionFlagKaraoke         = DispositionFlag(C.AV_DISPOSITION_KARAOKE)
	DispositionFlagLyrics          = DispositionFlag(C.AV_DISPOSITION_LYRICS)
	DispositionFlagMetadata        = DispositionFlag(C.AV_DISPOSITION_METADATA)
	DispositionFlagMultilayer      = DispositionFlag(C.AV_DISPOSITION_MULTILAYER)
	DispositionFlagNonDiegetic     = DispositionFlag(C.AV_DISPOSITION_NON_DIEGETIC)
	DispositionFlagOriginal        = DispositionFlag(C.AV_DISPOSITION_ORIGINAL)
	DispositionFlagStillImage      = DispositionFlag(C.AV_DISPOSITION_STILL_IMAGE)
	DispositionFlagTimedThumbnails = DispositionFlag(C.AV_DISPOSITION_TIMED_THUMBNAILS)
	DispositionFlagVisualImpaired  = DispositionFlag(C.AV_DISPOSITION_VISUAL_IMPAIRED)
)
