package rawfs

import (
	"os"
	"path"
	"path/filepath"

	"github.com/sgostarter/i/stg"
	"github.com/sgostarter/libeasygo/pathutils"
)

func NewFSStorage(rootPath string) stg.FileStorage {
	if rootPath == "" {
		rootPath, _ = os.Getwd()
	}

	return &fsStorageImpl{
		rootPath: rootPath,
	}
}

type fsStorageImpl struct {
	rootPath string
}

func (impl *fsStorageImpl) WriteFile(name string, data []byte) error {
	if !path.IsAbs(name) {
		name = filepath.Join(impl.rootPath, name)
	}

	_ = pathutils.MustDirOfFileExists(name)

	if impl.trySafeWriteFile(name, data) {
		return nil
	}

	return os.WriteFile(name, data, 0600)
}

func (impl *fsStorageImpl) ReadFile(name string) ([]byte, error) {
	if !path.IsAbs(name) {
		name = filepath.Join(impl.rootPath, name)
	}

	return os.ReadFile(name)
}

func (impl *fsStorageImpl) trySafeWriteFile(name string, data []byte) (ok bool) {
	exists, err := pathutils.IsFileExists(name)
	if err != nil {
		return
	}

	if !exists {
		return
	}

	nameBak := name + ".bak"

	if o, e := pathutils.IsFileExists(nameBak); e == nil && o {
		_ = os.Remove(nameBak)
	}

	err = os.Rename(name, nameBak)
	if err != nil {
		return
	}

	err = os.WriteFile(name, data, 0600)
	if err != nil {
		return
	}

	_ = os.Remove(nameBak)

	ok = true

	return
}
