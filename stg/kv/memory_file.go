package kv

import (
	"os"
	"sync"

	"github.com/sgostarter/i/commerr"
	"github.com/sgostarter/libeasygo/pathutils"
	"gopkg.in/yaml.v3"
)

func NewMemoryFileStorage(fileName string) (StorageTiny, error) {
	impl := &memoryFileStorageImpl{
		fileName: fileName,
		m:        make(map[string]string),
	}

	if err := impl.init(); err != nil {
		return nil, err
	}

	return impl, nil
}

type memoryFileStorageImpl struct {
	fileName string

	mLock sync.Mutex
	m     map[string]string
}

//
//
//

func (impl *memoryFileStorageImpl) init() error {
	return impl.load()
}

func (impl *memoryFileStorageImpl) load() error {
	if impl.fileName == "" {
		return commerr.ErrInvalidArgument
	}

	_ = pathutils.MustDirOfFileExists(impl.fileName)

	d, err := os.ReadFile(impl.fileName)
	if err != nil {
		if _, ok := err.(*os.PathError); ok {
			err = nil
		}

		return err
	}

	var m map[string]string
	err = yaml.Unmarshal(d, &m)
	if err != nil {
		return err
	}

	impl.m = m

	if impl.m == nil {
		impl.m = make(map[string]string)
	}

	return nil
}

func (impl *memoryFileStorageImpl) save() error {
	d, err := yaml.Marshal(impl.m)
	if err != nil {
		return err
	}

	err = os.WriteFile(impl.fileName, d, os.ModePerm)
	if err != nil {
		return err
	}

	return nil
}

//
//
//

func (impl *memoryFileStorageImpl) Set(key string, v interface{}) (err error) {
	d, err := yaml.Marshal(v)
	if err != nil {
		return
	}

	impl.mLock.Lock()
	defer impl.mLock.Unlock()

	impl.m[key] = string(d)

	return impl.save()
}

func (impl *memoryFileStorageImpl) Get(key string, v interface{}) (ok bool, err error) {
	impl.mLock.Lock()
	defer impl.mLock.Unlock()

	d, ok := impl.m[key]
	if !ok {
		return
	}

	err = yaml.Unmarshal([]byte(d), v)
	if err != nil {
		return
	}

	ok = true

	return
}

func (impl *memoryFileStorageImpl) Del(key string) error {
	impl.mLock.Lock()
	defer impl.mLock.Unlock()

	delete(impl.m, key)

	return impl.save()
}

func (impl *memoryFileStorageImpl) getDataList() (values [][]byte, _ error) {
	impl.mLock.Lock()
	defer impl.mLock.Unlock()

	for _, v := range impl.m {
		values = append(values, []byte(v))
	}

	return
}

func (impl *memoryFileStorageImpl) getDataMap() (values map[string][]byte, err error) {
	values = make(map[string][]byte)

	impl.mLock.Lock()
	defer impl.mLock.Unlock()

	for k, v := range impl.m {
		values[k] = []byte(v)
	}

	return
}

func (impl *memoryFileStorageImpl) GetList(itemGen func() interface{}) (items []interface{}, err error) {
	if itemGen() == nil {
		err = commerr.ErrInvalidArgument

		return
	}

	values, err := impl.getDataList()
	if err != nil {
		return
	}

	for _, value := range values {
		item := itemGen()
		if item == nil {
			err = commerr.ErrNotFound

			break
		}

		err = yaml.Unmarshal(value, item)
		if err != nil {
			continue
		}

		items = append(items, item)
	}

	return
}

func (impl *memoryFileStorageImpl) GetMap(itemGen func() interface{}) (items map[string]interface{}, err error) {
	if itemGen() == nil {
		err = commerr.ErrInvalidArgument

		return
	}

	values, err := impl.getDataMap()
	if err != nil {
		return
	}

	items = make(map[string]interface{})

	for key, value := range values {
		item := itemGen()
		if item == nil {
			err = commerr.ErrNotFound

			break
		}

		err = yaml.Unmarshal(value, item)
		if err != nil {
			continue
		}

		items[key] = item
	}

	return
}
