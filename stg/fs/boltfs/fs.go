package boltfs

import (
	"path"
	"strings"

	"github.com/boltdb/bolt"
	"github.com/sgostarter/i/commerr"
	"github.com/sgostarter/i/stg"
	"github.com/sgostarter/libeasygo/pathutils"
)

func NewFileStorage(path string) (storage stg.FileStorage, err error) {
	impl := &boltFileStorageImpl{
		path: path,
	}

	err = impl.init()
	if err != nil {
		return
	}

	storage = impl

	return
}

type boltFileStorageImpl struct {
	path string
	db   *bolt.DB
}

func (impl *boltFileStorageImpl) init() (err error) {
	_ = pathutils.MustDirOfFileExists(impl.path)

	db, err := bolt.Open(impl.path, 0600, nil)
	if err != nil {
		return
	}

	impl.db = db

	return
}

func (impl *boltFileStorageImpl) splitFilePath(name string) (string, string) {
	dir, file := path.Split(name)

	dir = strings.ReplaceAll(dir, "\\", "/")

	return dir, file
}

func (impl *boltFileStorageImpl) WriteFile(name string, d []byte) error {
	return impl.db.Update(func(tx *bolt.Tx) (err error) {
		dir, file := impl.splitFilePath(name)
		bucket, err := tx.CreateBucketIfNotExists([]byte(dir))
		if err != nil {
			return
		}

		err = bucket.Put([]byte(file), d)

		return
	})
}

func (impl *boltFileStorageImpl) ReadFile(name string) (d []byte, err error) {
	err = impl.db.View(func(tx *bolt.Tx) error {
		dir, file := impl.splitFilePath(name)
		bucket := tx.Bucket([]byte(dir))
		if bucket == nil {
			return commerr.ErrNotFound
		}

		d = bucket.Get([]byte(file))

		return nil
	})

	return
}
