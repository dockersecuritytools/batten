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
	"strconv"
	"strings"
)

//
// Meminfo is a parser for /proc/meminfo.
//
//
type Meminfo struct {
	MemTotal      int64
	MemFree       int64
	Buffers       int64
	Cached        int64
	SwapCached    int64
	Active        int64
	Inactive      int64
	HighTotal     int64
	HighFree      int64
	LowTotal      int64
	LowFree       int64
	SwapTotal     int64
	SwapFree      int64
	Dirty         int64
	Writeback     int64
	AnonPages     int64
	Mapped        int64
	Slab          int64
	SReclaimable  int64
	SUnreclaim    int64
	PageTables    int64
	NFS_Unstable  int64
	Bounce        int64
	WritebackTmp  int64
	CommitLimit   int64
	Committed_AS  int64
	VmallocTotal  int64
	VmallocUsed   int64
	VmallocChunk  int64
	AnonHugePages int64
}

func linesToMeminfo(lines []string) (map[string]int64, error) {
	var result map[string]int64 = make(map[string]int64)

	// first line is the header
	for i := 0; i < len(lines); i++ {

		lines[i] = strings.TrimSpace(lines[i])

		if len(lines[i]) == 0 {
			// it's empty
			continue
		}

		parts := strings.Fields(lines[i])

		if len(parts) < 2 {
			log.Println("malformed line, expected 2 parts but only got:", len(parts), "line:", lines[i])
			continue
		}

		title := strings.Replace(parts[0], ":", "", -1)
		val, err := strconv.ParseInt(parts[1], 10, 64)
		if err != nil {
			return nil, err
		}

		result[title] = val
	}
	return result, nil
}

//
// ParseMeminfo() loads a Meminfo object from a supplied path string.
//
// If the path cannot be found, it will return an error.
//
func ParseMeminfo(path string) (*Meminfo, error) {

	buf, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	lines := strings.Split(string(buf), "\n")
	meminfos, err := linesToMeminfo(lines)
	if err != nil {
		return nil, err
	}

	meminfo := &Meminfo{}
	meminfo.MemTotal      = meminfos["MemTotal"]
	meminfo.MemFree       = meminfos["MemFree"]
	meminfo.Buffers       = meminfos["Buffers"]
	meminfo.Cached        = meminfos["Cached"]
	meminfo.SwapCached    = meminfos["SwapCached"]
	meminfo.Active        = meminfos["Active"]
	meminfo.Inactive      = meminfos["Inactive"]
	meminfo.HighTotal     = meminfos["HighTotal"]
	meminfo.HighFree      = meminfos["HighFree"]
	meminfo.LowTotal      = meminfos["LowTotal"]
	meminfo.LowFree       = meminfos["LowFree"]
	meminfo.SwapTotal     = meminfos["SwapTotal"]
	meminfo.SwapFree      = meminfos["SwapFree"]
	meminfo.Dirty         = meminfos["Dirty"]
	meminfo.Writeback     = meminfos["Writeback"]
	meminfo.AnonPages     = meminfos["AnonPages"]
	meminfo.Mapped        = meminfos["Mapped"]
	meminfo.Slab          = meminfos["Slab"]
	meminfo.SReclaimable  = meminfos["SReclaimable"]
	meminfo.SUnreclaim    = meminfos["SUnreclaim"]
	meminfo.PageTables    = meminfos["PageTables"]
	meminfo.NFS_Unstable  = meminfos["NFS_Unstable"]
	meminfo.Bounce        = meminfos["Bounce"]
	meminfo.WritebackTmp  = meminfos["WritebackTmp"]
	meminfo.CommitLimit   = meminfos["CommitLimit"]
	meminfo.Committed_AS  = meminfos["Committed_AS"]
	meminfo.VmallocTotal  = meminfos["VmallocTotal"]
	meminfo.VmallocUsed   = meminfos["VmallocUsed"]
	meminfo.VmallocChunk  = meminfos["VmallocChunk"]
	meminfo.AnonHugePages = meminfos["AnonHugePages"]

	return meminfo, nil

}

//
// NewMeminfo() creates a new Meminfo object that loads from /proc/meminfo.
//
//
func NewMeminfo() (*Meminfo, error) {
	return ParseMeminfo("/proc/meminfo")
}
