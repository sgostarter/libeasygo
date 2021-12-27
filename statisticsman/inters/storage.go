package inters

type Storage interface {
	Inc(key, filed string, incV int64)
}
