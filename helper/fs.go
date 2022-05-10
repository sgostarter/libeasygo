package helper

import "os"

func PathExists(path string) (exists bool, err error) {
	_, err = os.Stat(path)
	if err == nil {
		exists = true

		return
	}

	if os.IsNotExist(err) {
		err = nil
	}

	return
}

func IsDir(path string) bool {
	s, err := os.Stat(path)
	if err != nil {
		return false
	}

	return s.IsDir()
}

func IsFile(path string) bool {
	s, err := os.Stat(path)
	if err != nil {
		return false
	}

	return !s.IsDir()
}
