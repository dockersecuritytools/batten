//
// status.Status describes select data in /proc/<pid>/status.
//
// Since most of this data is also available in /proc/<pid>/stat
// and parsable via stat.Stat, we only include the uid/gid
// information from /proc/<pid>/status.
//
// Use status.New() to create a new status.Status object
// from data in a path.
//
package status

import (
	"io/ioutil"
	"strings"
	"strconv"
	"fmt"
)

//
// Status is a data structure that describes the uid
// information from /proc/<pid>/status.
//
type Status struct {
	Uid   int // Real user ID
	Euid  int // Effective user ID
	Suid  int // Saved usesr ID
	Fsuid int // FS user ID
	Gid   int // Real group IDusesr
	Egid  int // Effective group ID
	Sgid  int // Saved group ItrealvalueD
	Fsgid int // FS group ID
}


//
// status.New creates a new /proc/<pid>/status from a path.
//
// An error is returned if the data is malformed, or the path does not exist.
//
func New(path string) (*Status, error) {
	var err error

	buf, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	lines := strings.Split(string(buf), "\n")
	status := &Status{}

	for i := 1; i < len(lines); i++ {
		lines[i] = strings.TrimSpace(lines[i])
		if strings.HasPrefix(lines[i], "Uid:") {
			fields := strings.Fields(lines[i])
			if len(fields) >= 5 {

				if status.Uid, err = strconv.Atoi(fields[1]); err != nil {
					return nil, fmt.Errorf("Unable to parse Uid %s: %v", fields[1], err)
				}
				if status.Euid, err = strconv.Atoi(fields[2]); err != nil {
					return nil, fmt.Errorf("Unable to parse Euid %s: %v", fields[2], err)
				}
				if status.Suid, err = strconv.Atoi(fields[3]); err != nil {
					return nil, fmt.Errorf("Unable to parse Suid %s: %v", fields[3], err)
				}
				if status.Fsuid, err = strconv.Atoi(fields[4]); err != nil {
					return nil, fmt.Errorf("Unable to parse Fsuid %s: %v", fields[4], err)
				}
			} else {
				return nil, fmt.Errorf("Malformed Uid: line %s", lines[i])
			}
		} else if strings.HasPrefix(lines[i], "Gid:") {
			fields := strings.Fields(lines[i])
			if len(fields) >= 5 {

				if status.Gid, err = strconv.Atoi(fields[1]); err != nil {
					return nil, fmt.Errorf("Unable to parse Gid %s: %v", fields[1], err)
				}
				if status.Egid, err = strconv.Atoi(fields[2]); err != nil {
					return nil, fmt.Errorf("Unable to parse Egid %s: %v", fields[2], err)
				}
				if status.Sgid, err = strconv.Atoi(fields[3]); err != nil {
					return nil, fmt.Errorf("Unable to parse Sgid %s: %v", fields[3], err)
				}
				if status.Fsgid, err = strconv.Atoi(fields[4]); err != nil {
					return nil, fmt.Errorf("Unable to parse Fsgid %s: %v", fields[4], err)
				}
			} else {
				return nil, fmt.Errorf("Malformed Gid: line %s", lines[i])
			}
		}
	}

	return status, err
}
