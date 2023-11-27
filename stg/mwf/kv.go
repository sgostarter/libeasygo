package mwf

import (
	"encoding/json"
	"sync"

	"github.com/sgostarter/i/commerr"
	"github.com/sgostarter/i/stg"
	"github.com/sgostarter/libeasygo/stg/kv"
)

func NewKV(file string) kv.StorageTiny {
	return NewKVEx(file, nil)
}

func NewKVEx(file string, storage stg.FileStorage) kv.StorageTiny {
	return &kvImpl{
		d: NewMemWithFile[map[string]string, Serial, Lock](make(map[string]string), &JSONSerial{}, &sync.RWMutex{}, file, storage),
	}
}

type kvImpl struct {
	d *MemWithFile[map[string]string, Serial, Lock]
}

func (impl *kvImpl) GetList(itemGen func(key string) interface{}) (items []interface{}, err error) {
	if itemGen == nil {
		err = commerr.ErrInvalidArgument

		return
	}

	impl.d.Read(func(values map[string]string) {
		for key, value := range values {
			item := itemGen(key)
			if item == nil {
				continue
			}

			err = json.Unmarshal([]byte(value), &item)
			if err != nil {
				continue
			}

			items = append(items, item)
		}
	})

	return
}

func (impl *kvImpl) GetMap(itemGen func(key string) interface{}) (items map[string]interface{}, err error) {
	if itemGen == nil {
		err = commerr.ErrInvalidArgument

		return
	}

	items = make(map[string]interface{})

	impl.d.Read(func(values map[string]string) {
		for key, value := range values {
			item := itemGen(key)
			if item == nil {
				continue
			}

			err = json.Unmarshal([]byte(value), &item)
			if err != nil {
				continue
			}

			items[key] = item
		}
	})

	return
}

func (impl *kvImpl) Set(key string, v interface{}) error {
	d, err := json.Marshal(v)
	if err != nil {
		return err
	}

	return impl.d.Change(func(v map[string]string) (newV map[string]string, err error) {
		newV = v

		if newV == nil {
			newV = make(map[string]string)
		}

		newV[key] = string(d)

		return
	})
}

func (impl *kvImpl) Get(key string, v interface{}) (ok bool, err error) {
	var d string

	impl.d.Read(func(v map[string]string) {
		d, ok = v[key]
	})

	if !ok {
		return
	}

	err = json.Unmarshal([]byte(d), v)
	if err != nil {
		return
	}

	ok = true

	return
}

func (impl *kvImpl) Del(key string) error {
	return impl.d.Change(func(v map[string]string) (newV map[string]string, err error) {
		newV = v

		if newV == nil {
			newV = make(map[string]string)
		}

		delete(newV, key)

		return
	})
}
