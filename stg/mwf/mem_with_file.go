package mwf

import (
	"fmt"
	"os"
	"time"

	"github.com/sgostarter/i/stg"
	"github.com/sgostarter/libeasygo/pathutils"
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

type EventObserver[T any] interface {
	BeforeLoad()
	AfterLoad(memD T, err error)
	BeforeSave()
	AfterSave(memD T, err error)
}

type MemWithFile[T any, S Serial, L Lock] struct {
	memD   T
	serial S
	lock   L

	fileName string
	storage  stg.FileStorage
	ob       EventObserver[T]
}

func NewMemWithFile[T any, S Serial, L Lock](d T, serial S, lock L, fileName string, storage stg.FileStorage) *MemWithFile[T, S, L] {
	return NewMemWithFileEx(d, serial, lock, fileName, storage, nil)
}

func NewMemWithFileEx[T any, S Serial, L Lock](d T, serial S, lock L, fileName string, storage stg.FileStorage, ob EventObserver[T]) *MemWithFile[T, S, L] {
	if storage == nil && fileName != "" {
		storage = rawfs.NewFSStorage("")
	}

	mwf := &MemWithFile[T, S, L]{
		memD:     d,
		serial:   serial,
		lock:     lock,
		fileName: fileName,
		storage:  storage,
		ob:       ob,
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

	ok, err := pathutils.IsFileExists(mwf.tmpFile())
	if err == nil && ok {
		_ = os.Rename(mwf.fileName, fmt.Sprintf("%s.r.%d", mwf.fileName, time.Now().UnixMilli()))

		_ = os.Rename(mwf.tmpFile(), mwf.fileName)
	}

	if mwf.ob != nil {
		mwf.ob.BeforeLoad()
	}

	d, err := mwf.storage.ReadFile(mwf.fileName)
	if err != nil {
		if _, ok := err.(*os.PathError); ok {
			err = nil
		}

		if mwf.ob != nil {
			mwf.ob.AfterLoad(mwf.memD, err)
		}

		return err
	}

	var m T

	err = mwf.serial.Unmarshal(d, &m)
	if err != nil {
		if mwf.ob != nil {
			mwf.ob.AfterLoad(mwf.memD, err)
		}

		return err
	}

	mwf.memD = m

	if mwf.ob != nil {
		mwf.ob.AfterLoad(mwf.memD, nil)
	}

	return nil
}

func (mwf *MemWithFile[T, S, L]) tmpFile() string {
	return mwf.fileName + ".tmp"
}

func (mwf *MemWithFile[T, S, L]) save() error {
	if mwf.fileName == "" {
		return nil
	}

	if mwf.ob != nil {
		mwf.ob.BeforeSave()
	}

	d, err := mwf.serial.Marshal(mwf.memD)
	if err != nil {
		if mwf.ob != nil {
			mwf.ob.AfterSave(mwf.memD, err)
		}

		return err
	}

	err = os.Rename(mwf.fileName, mwf.tmpFile())
	if err != nil {
		if mwf.ob != nil {
			mwf.ob.AfterSave(mwf.memD, err)
		}

		return nil
	}

	err = mwf.storage.WriteFile(mwf.fileName, d)
	if err != nil {
		_ = os.Rename(mwf.tmpFile(), mwf.fileName)

		if mwf.ob != nil {
			mwf.ob.AfterSave(mwf.memD, err)
		}

		return err
	}

	_ = os.Remove(mwf.tmpFile())

	if mwf.ob != nil {
		mwf.ob.AfterSave(mwf.memD, nil)
	}

	return nil
}
