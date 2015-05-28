package procfs

// Copyright Jen Andre (jandre@gmail.com)
//
// Permission is hereby granted, free of charge, to any person obtaining
// a copy of this software and associated documentation files (the
// "Software"), to deal in the Software without restriction, including
// without limitation the rights to use, copy, modify, merge, publish,
// distribute, sublicense, and/or sell copies of the Software, and to
// permit persons to whom the Software is furnished to do so, subject to
// the following conditions:
//
// The above copyright notice and this permission notice shall be
// included in all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
// EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF
// MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
// NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE
// LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION
// OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION
// WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

import (
	"io/ioutil"
	"log"
	"os"
	"path"
	"strconv"
	"strings"

	"github.com/jandre/procfs/limits"
	"github.com/jandre/procfs/stat"
	"github.com/jandre/procfs/statm"
)

//
// Process describes a /proc/<pid> entry
//
type Process struct {
	Pid       int               // Process ID
	Environ   map[string]string // Environment variables
	Cmdline   []string          // Command line of process (argv array)
	Cwd       string            // Process current working directory
	Exe       string            // Symlink to executed command.
	Root      string            // Per-process root (e.g. chroot)
	prefix    string            // directory path, e.g. /proc/<pid>
	stat      *stat.Stat        // Status information from /proc/pid/stat - see Stat()
	statm     *statm.Statm      // Memory Status information from /proc/pid/statm - see Statm()
	limits    *limits.Limits    // Per process rlimit settings from /proc/pid/limits - see Limits()
	loginuid  *int              // Maybe loginuid from /proc/pid/loginuid - see Loginuid()
	sessionid *int              // Maybe sessionid from /proc/pid/sessionid- see Sessionid()

}

//
// Read /proc information for `pid`
//
func NewProcess(pid int, lazy bool) (*Process, error) {
	prefix := path.Join("/proc", strconv.Itoa(pid))
	return NewProcessFromPath(pid, prefix, lazy)
}

//
// Read a process entry from a directory path
//
// if lazy = true, then preload the stat, limits, and other files.
//
func NewProcessFromPath(pid int, prefix string, lazy bool) (*Process, error) {
	var err error

	if _, err = os.Stat(prefix); err != nil {
		// error reading pid
		return nil, err
	}

	process := &Process{prefix: prefix, Pid: pid}
	if process.Cmdline, err = readCmdLine(prefix); err != nil {
		return nil, err
	}

	if process.Exe, err = readLink(prefix, "exe"); err != nil {
		process.Exe = ""
	}

	if process.Cwd, err = readLink(prefix, "cwd"); err != nil {
		return nil, err
	}

	if process.Root, err = readLink(prefix, "root"); err != nil {
		return nil, err
	}

	process.readEnviron()

	if !lazy {
		// preload all of the subdirs
		process.Stat()
		process.Limits()
		process.Loginuid()
		process.Sessionid()
	}

	return process, nil
}

const NO_VALUE = 4294967295

func clearEmpty(strings []string) []string {
	var filtered []string
	for _, v := range strings {
		if len(v) != 0 {
			filtered = append(filtered, v)
		}
	}
	return filtered
}

//
// read integer from file
//
func readFileInteger(prefix string, file string) (*int, error) {
	filename := path.Join(prefix, file)
	str, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	if val, err := strconv.Atoi(string(str)); err != nil {
		return nil, err
	} else {
		return &val, nil
	}
}

//

//
// read link
//
func readLink(prefix string, file string) (string, error) {
	filename := path.Join(prefix, file)
	str, err := os.Readlink(filename)
	if err != nil {
		return "", err
	}
	return string(str), nil
}

//
// Read a /proc/<pid>/cmdline file and break up into array
// of argv
//
func readCmdLine(prefix string) ([]string, error) {
	filename := path.Join(prefix, "cmdline")
	str, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return clearEmpty(strings.Split(string(str), "\x00")), nil
}

//
// Parser for /proc/<pid>/stat
//
func (p *Process) Stat() (*stat.Stat, error) {
	var err error
	if p.stat == nil {
		if p.stat, err = stat.New(path.Join(p.prefix, "stat")); err != nil {
			return nil, err
		}
	}

	return p.stat, nil
}

//
// Parser for /proc/<pid>/limits
//
func (p *Process) Limits() (*limits.Limits, error) {
	var err error
	if p.limits == nil {
		if p.limits, err = limits.New(path.Join(p.prefix, "limits")); err != nil {
			return nil, err
		}
	}

	return p.limits, nil
}

//
// Parses contents of /proc/pid/loginuid (if present)
//
func (p *Process) Loginuid() int {
	var err error
	if p.loginuid == nil {
		p.loginuid, err = readFileInteger(p.prefix, "loginuid")
		if err != nil {
			log.Println("Warning: unable to read loginuid so returning nil:", err)
		}

	}

	if p.loginuid != nil {
		return *p.loginuid
	}
	return NO_VALUE

}

//
// Parses contents of /proc/pid/sessionid (if present)
//
func (p *Process) Sessionid() int {
	var err error
	if p.sessionid == nil {
		p.sessionid, err = readFileInteger(p.prefix, "sessionid")
		if err != nil {
			log.Println("Warning: unable to read sessionid so returning nil:", err)
		}

	}

	if p.sessionid != nil {
		return *p.sessionid
	}
	return NO_VALUE
}

func (p *Process) readEnviron() {
	p.Environ = make(map[string]string, 0)
	bytes, err := ioutil.ReadFile(path.Join(p.prefix, "environ"))
	if err != nil {
		return
	}
	for _, item := range strings.Split(string(bytes), "\x00") {
		fields := strings.SplitN(item, "=", 2)
		if (len(fields) == 2) {
			p.Environ[fields[0]] = fields[1]
		}
	}
}
