package batten

import (
	"errors"
	"io/ioutil"
	"os"
	"path"
)

type FilePermsCheck struct {
	filepath string
	// TODO set via config
	targetPerms uint32
}

func (fo *FilePermsCheck) HasPerms(targetMode os.FileMode) (bool, error) {
	fi, err := os.Stat(fo.filepath)
	if err != nil {
		return false, err
	}

	mode := fi.Mode()
	return (mode & targetMode) == targetMode, nil
}

func atLeastPerms(filepath string, targetMode os.FileMode) (bool, error) {
	fi, err := os.Stat(filepath)
	if err != nil {
		return false, err
	}

	// get just the permission bits
	mode := fi.Mode() & os.ModePerm

	// now see if the permission is less than the target permissions
	return (mode <= targetMode), nil
}

func (fo *FilePermsCheck) HasAtLeastPerms(targetMode os.FileMode) (bool, error) {
	return atLeastPerms(fo.filepath, targetMode)
}

func (fo *FilePermsCheck) HasAtLeastPermsRecursive(targetMode os.FileMode) (bool, error) {

	succ, err := atLeastPerms(fo.filepath, targetMode)

	if !succ || err != nil {
		return succ, err
	}

	files, err := ioutil.ReadDir(fo.filepath)
	for _, file := range files {
		fullpath := path.Join(fo.filepath, file.Name())
		succ, err := atLeastPerms(fullpath, targetMode)
		if !succ || err != nil {
			return succ, err
		}
	}

	return true, nil
}

func (fo *FilePermsCheck) validateFromArgs(lookForFlag string, args []string) (bool, error) {
	filepath := getArgValue(lookForFlag, args)

	if filepath != "" {
		if PathExists(filepath) {
			fo.filepath = filepath
			return fo.HasAtLeastPerms(os.FileMode(fo.targetPerms))
		} else {
			return false, errors.New("Could not find path: " + filepath)
		}
	}
	// it's ok, no tls config set
	return true, nil
}
