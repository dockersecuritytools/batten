package procfs
//
// License information:
//
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
	"strconv"
)

//
// Take a list of FileInfo from /proc listing and only return 
// pid directories
//
func filterByPidDirectories(files []os.FileInfo) []int {
	var result []int

	for _, file := range files {
		if file.IsDir() {
			if pid,  err := strconv.Atoi(file.Name()); err == nil {
				result = append(result, pid)
			}
		}
	}
	return result
}

//
// Load all processes from /proc
//
// If lazy = true, do not load ancillary information (/proc/<pid>/stat, 
// /proc/<pid>/statm, etc) immediately - load it on demand 
//
func Processes(lazy bool) (map[int]*Process, error) {
	processes := make(map[int]*Process)
	done := make(chan *Process)
	files, err := ioutil.ReadDir("/proc")
	if err != nil {
		return nil, err
	}

	pids := filterByPidDirectories(files)

	fetch := func(pid int) {
		proc, err := NewProcess(pid, lazy)
		if err != nil {
			// TODO: bubble errors up if requested
			log.Println("Failure loading process pid: ", pid, err)
			done <- nil
		} else {
			done <- proc
		}
	}

	todo := len(pids)

	// create a goroutine that asynchronously processes all /proc/<pid> entries
	for _, pid := range pids {
		go fetch(pid)
	}

	// 
	// fetch all processes until we're done
	//
	for ;todo > 0; {
		proc := <-done
		todo--
		if proc != nil {
			processes[proc.Pid] = proc
		}
	}

	return processes, nil
}
