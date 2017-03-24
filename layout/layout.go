package layout

type MinSizer interface {
	MinSize() (int, int)
}

type MaxSizer interface {
	MaxSize() (int, int)
}
