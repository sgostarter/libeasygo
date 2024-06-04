package kv

type Storage interface {
	Set(key string, v interface{}) error
	Get(key string, v interface{}) (ok bool, err error)
	Del(key string) error
}

type StorageCollect interface {
	GetList(itemGen func(key string) interface{}) (items []interface{}, err error)
	GetMap(itemGen func(key string) interface{}) (items map[string]interface{}, err error)
}

type StorageTiny interface {
	Storage
	StorageCollect
}

type Storage2 interface {
	Storage

	Sets(keys []string, vs ...interface{}) error
	Gets(keys []string, vsi ...interface{}) (vs []interface{}, err error)
	Dels(keys []string) error
}

type StorageTiny2 interface {
	Storage2
	StorageCollect
}
