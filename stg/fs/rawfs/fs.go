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

	return os.WriteFile(name, data, 0600)
}

func (impl *fsStorageImpl) ReadFile(name string) ([]byte, error) {
	if !path.IsAbs(name) {
		name = filepath.Join(impl.rootPath, name)
	}

	return os.ReadFile(name)
}
