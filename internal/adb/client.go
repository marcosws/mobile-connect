package adb

type ADBClient interface {
	Devices() (string, error)
	GetProp(serial, prop string) (string, error)
}
