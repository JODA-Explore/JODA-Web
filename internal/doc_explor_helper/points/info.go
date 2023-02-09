package points

type Info struct {
	CountTotal                                                                     uint64
	CountNull                                                                      uint64
	CountObject                                                                    uint64
	CountArray                                                                     uint64
	CountString                                                                    uint64
	CountNumber, CountInt, CountFloat                                              uint64
	CountBoolean, CountTrue, CountFalse                                            uint64
	MaxMember, MinMember, MaxStrSize, MinStrSize, MaxInt, MinInt, MaxSize, MinSize uint64
}

type Value struct {
	Info Info
	// Indicies [backends.Num][]int
	// Ctl      [backends.Num]trie.Control
}
