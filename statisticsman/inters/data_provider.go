package inters

type DataScannerCB func(rKey, k string, v int64, err error) error

type DataProvider interface {
	Scan(k string, cb DataScannerCB) error
	Exists(k string) (exists bool, err error)
	Delete(k string) error

	ScanEx(k string, cb DataScannerCB, reset bool) error
}
