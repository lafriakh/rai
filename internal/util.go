package internal

import (
	"os"
)

func openOrCreateFile(fname string) (*os.File, error) {
	file, err := os.OpenFile(fname, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return nil, err
	}

	return file, nil
}

func isFileExists(fname string) bool {
	_, err := os.Stat(fname)

	return !os.IsNotExist(err)
}
