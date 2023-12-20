package mwf

import (
	"os"

	"github.com/sgostarter/i/stg"
	"github.com/sgostarter/libeasygo/stg/fs/rawfs"
)

type Serial interface {
	Marshal(t any) ([]byte, error)
	Unmarshal(d []byte, t any) error
}

type Lock interface {
	RLock()
	RUnlock()

	Lock()
	Unlock()
}

type MemWithFile[T any, S Serial, L Lock] struct {
	memD   T
	serial S
	lock   L

	fileName string
	storage  stg.FileStorage
}

func NewMemWithFile[T any, S Serial, L Lock](d T, serial S, lock L, fileName string, storage stg.FileStorage) *MemWithFile[T, S, L] {
	if storage == nil && fileName != "" {
		storage = rawfs.NewFSStorage("")
	}

	mwf := &MemWithFile[T, S, L]{
		memD:     d,
		serial:   serial,
		lock:     lock,
		fileName: fileName,
		storage:  storage,
	}

	_ = mwf.load()

	return mwf
}

func (mwf *MemWithFile[T, S, L]) Read(proc func(memD T)) {
	mwf.lock.RLock()
	defer mwf.lock.RUnlock()

	proc(mwf.memD)
}

func (mwf *MemWithFile[T, S, L]) Change(proc func(memD T) (newMemD T, err error)) error {
	mwf.lock.Lock()
	defer mwf.lock.Unlock()

	newMemD, err := proc(mwf.memD)
	if err != nil {
		return err
	}

	mwf.memD = newMemD

	return mwf.save()
}

func (mwf *MemWithFile[T, S, L]) load() error {
	if mwf.fileName == "" {
		return nil
	}

	d, err := mwf.storage.ReadFile(mwf.fileName)
	if err != nil {
		if _, ok := err.(*os.PathError); ok {
			err = nil
		}

		return err
	}

	var m T

	err = mwf.serial.Unmarshal(d, &m)
	if err != nil {
		return err
	}

	mwf.memD = m

	return nil
}

func (mwf *MemWithFile[T, S, L]) save() error {
	if mwf.fileName == "" {
		return nil
	}

	d, err := mwf.serial.Marshal(mwf.memD)
	if err != nil {
		return err
	}

	err = mwf.storage.WriteFile(mwf.fileName, d)
	if err != nil {
		return err
	}

	return nil
}
