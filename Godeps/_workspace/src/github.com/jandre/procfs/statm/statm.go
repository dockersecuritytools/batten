//
// statm.Statm describes data in /proc/<pid>/statm.
//
// Use statm.New() to create a new stat.Statm object
// from data in a path.
//
package statm

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
//

import (
	"github.com/jandre/procfs/util"
	"io/ioutil"
	"strings"
)

//
// Abstraction for /proc/<pid>/stat
//
type Statm struct {
	Size     int64 // total program size (pages)(same as VmSize in status)
	Resident int64 //size of memory portions (pages)(same as VmRSS in status)
	Shared   int   // number of pages that are shared(i.e. backed by a file)
	Trs      int   // number of pages that are 'code'(not including libs; broken, includes data segment)
	Lrs      int   //number of pages of library(always 0 on 2.6)
	Drs      int   //number of pages of data/stack(including libs; broken, includes library text)
	Dt       int   //number of dirty pages(always 0 on 2.6)
}

func New(path string) (*Statm, error) {

	buf, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	lines := strings.Split(string(buf), " ")
	stat := &Statm{}
	err = util.ParseStringsIntoStruct(stat, lines)
	return stat, err
}
