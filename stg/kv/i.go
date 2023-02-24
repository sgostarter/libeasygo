package kv

type Storage interface {
	Set(key string, v interface{}) error
	Get(key string, v interface{}) (ok bool, err error)
	Del(key string) error
}

type StorageTiny interface {
	Storage
	GetList(itemGen func() interface{}) (items []interface{}, err error)
	GetMap(itemGen func() interface{}) (items map[string]interface{}, err error)
}
