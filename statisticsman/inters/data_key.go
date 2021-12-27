package inters

type DataKey interface {
	Key() string
	From(s string) error
}
