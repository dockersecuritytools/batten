package batten

import (
	"io/ioutil"
	"os"
	"os/user"
	"path"
	"strconv"
	"strings"
	"syscall"
)

type FileOwnerCheck struct {
	filepath  string
	uid       uint32
	gid       uint32
	username  string
	groupname string
}

func lookupUid(username string) (uint32, error) {
	u, err := user.Lookup(username)
	if err != nil {
		return 0, err
	}
	uid, err := strconv.Atoi(u.Uid)
	return uint32(uid), err
}

func getGroups() (res map[string]uint32, err error) {
	res = make(map[string]uint32, 0)
	bytes, err := ioutil.ReadFile("/etc/group")
	if err != nil {
		return nil, err
	}
	for _, line := range strings.Split(string(bytes), "\n") {
		items := strings.Split(line, ":")
		if len(items) >= 3 {
			groupname := items[0]
			val := items[2]
			i, err := strconv.Atoi(val)
			if err == nil {
				res[groupname] = uint32(i)
			}
		}
	}

	return res, nil
}

func lookupGid(groupname string) (uint32, error) {
	groups, err := getGroups()
	if err != nil {
		return 0, err
	}
	return groups[groupname], nil
}

func isOwner(filepath string, uid uint32) (bool, error) {
	fi, err := os.Stat(filepath)
	if err != nil {
		return false, err
	}

	actualUid := fi.Sys().(*syscall.Stat_t).Uid

	return actualUid == uid, nil
}

func isGroupOwner(filepath string, gid uint32) (bool, error) {
	fi, err := os.Stat(filepath)
	if err != nil {
		return false, err
	}

	actualGid := fi.Sys().(*syscall.Stat_t).Gid

	return actualGid == gid, nil
}

func isOwnerAndGroupOwner(filepath string, uid uint32, gid uint32) (bool, error) {

	if succ, err := isOwner(filepath, uid); err != nil {
		return false, err
	} else {
		if !succ {
			return false, nil
		}
		if groupsucc, err := isGroupOwner(filepath, gid); err != nil {
			return false, err
		} else {
			return groupsucc, nil
		}
	}
	return false, nil
}

func (fo *FileOwnerCheck) IsOwnerAndGroupOwnerRecursive(uid uint32, gid uint32) (bool, error) {

	succ, err := isOwnerAndGroupOwner(fo.filepath, uid, gid)

	if !succ || err != nil {
		return succ, err
	}

	files, err := ioutil.ReadDir(fo.filepath)
	for _, file := range files {
		fullpath := path.Join(fo.filepath, file.Name())
		succ, err := isOwnerAndGroupOwner(fullpath, uid, gid)
		if !succ || err != nil {
			return succ, err
		}
	}

	return true, nil
}

func (fo *FileOwnerCheck) IsOwnerAndGroupOwner(uid uint32, gid uint32) (bool, error) {

	return isOwnerAndGroupOwner(fo.filepath, uid, gid)
}

func (fo *FileOwnerCheck) IsOwner(uid uint32) (bool, error) {
	return isOwner(fo.filepath, uid)
}

func (fo *FileOwnerCheck) IsGroupOwner(gid uint32) (bool, error) {
	return isGroupOwner(fo.filepath, gid)

}

func (fo *FileOwnerCheck) validateOwnerAndGroupOwner() (bool, error) {

	uid, err := lookupUid(fo.username)
	if err != nil {
		return false, err
	}
	gid, err := lookupGid(fo.groupname)
	if err != nil {
		return false, err
	}

	return isOwnerAndGroupOwner(fo.filepath, uid, gid)
}
