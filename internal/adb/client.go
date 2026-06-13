package adb

type ADBClient interface {
	Devices() (string, error)
	GetProp(serial, prop string) (string, error)
	Shell(serial string, command string) (string, error)
}
