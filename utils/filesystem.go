package utils

import "os"

// PathExists Checks if a file path exists on the filesystem
// Args:
//
//	path   - string, the file path of the source file
//
// Returns:
//
//	out    - bool, whether the filepath was found to exist on the filesystem
func PathExists(path string) (bool, error) {
	// pull file data
	_, err := os.Stat(path)
	// if no error then file exists
	if err == nil {
		return true, nil
	}
	// if err check if it is the expect "file does not exist" error
	if os.IsNotExist(err) {
		return false, nil
	}
	// if file was found but error was thrown return error
	return true, err
}
