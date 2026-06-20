package adb

type ADBClient interface {
	Devices() (string, error)
	GetProp(serial, prop string) (string, error)
	Shell(serial string, command string) (string, error)
	InstallAPK(serial string, apkPath string) (string, error)
	UninstallAPK(serial string, packageName string) (string, error)
	ListPackages(serial string) ([]string, error)
	LaunchApp(deviceID string, pkg string) error
	StopApp(deviceID string, pkg string) error
}
