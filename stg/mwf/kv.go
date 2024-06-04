package mwf

import (
	"encoding/json"
	"sync"

	"github.com/sgostarter/i/commerr"
	"github.com/sgostarter/i/stg"
	"github.com/sgostarter/libeasygo/stg/kv"
)

func NewKV(file string) kv.StorageTiny2 {
	return NewKVEx(file, nil)
}

func NewKVEx(file string, storage stg.FileStorage) kv.StorageTiny2 {
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
	return impl.Sets([]string{key}, v)
}

func (impl *kvImpl) Get(key string, v interface{}) (ok bool, err error) {
	vs, err := impl.Gets([]string{key}, v)
	if err != nil {
		return
	}

	ok = vs[0] != nil

	return
}

func (impl *kvImpl) Del(key string) error {
	return impl.Dels([]string{key})
}

func (impl *kvImpl) Sets(keys []string, vs ...interface{}) error {
	if len(keys) != len(vs) {
		return commerr.ErrInvalidArgument
	}

	ds := make([][]byte, 0, len(keys))

	for _, v := range vs {
		d, err := json.Marshal(v)
		if err != nil {
			return err
		}

		ds = append(ds, d)
	}

	return impl.d.Change(func(v map[string]string) (newV map[string]string, err error) {
		newV = v

		if newV == nil {
			newV = make(map[string]string)
		}

		for idx := 0; idx < len(keys); idx++ {
			newV[keys[idx]] = string(ds[idx])
		}

		return
	})
}

func (impl *kvImpl) Gets(keys []string, vsi ...interface{}) (vs []interface{}, err error) {
	var ds []string

	impl.d.Read(func(v map[string]string) {
		for _, key := range keys {
			ds = append(ds, v[key])
		}
	})

	vs = make([]interface{}, len(keys))

	for idx := 0; idx < len(ds); idx++ {
		if ds[idx] == "" {
			vs[idx] = nil

			continue
		}

		if idx >= len(vsi) || vsi[idx] == nil {
			vs[idx] = ds[idx]

			continue
		}

		err = json.Unmarshal([]byte(ds[idx]), vsi[idx])
		if err != nil {
			return
		}

		vs[idx] = vsi[idx]
	}

	return
}

func (impl *kvImpl) Dels(keys []string) error {
	return impl.d.Change(func(v map[string]string) (newV map[string]string, err error) {
		newV = v

		if newV == nil {
			newV = make(map[string]string)
		}

		for _, key := range keys {
			delete(newV, key)
		}

		return
	})
}
