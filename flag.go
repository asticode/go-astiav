package astiav

type flags int

func (fs flags) add(f int) int { return int(fs) | f }

func (fs flags) del(f int) int { return int(fs) &^ f }

func (fs flags) has(f int) bool { return int(fs)&f > 0 }
