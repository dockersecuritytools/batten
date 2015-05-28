//
// limits.Limits describes data in /proc/<pid>/limits.
//
// Use limits.New() to create a new limits.Limits object
// from data in a limits file.
//
package limits

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
	"bytes"
	"io/ioutil"
	"log"
	"regexp"
	"strconv"
	"strings"
)

const UNLIMITED = -1

type Unit string

const (
	Bytes     Unit = "bytes"
	Seconds        = "seconds"
	Processes      = "processes"
	Files          = "files"
	Signals        = "signals"
	Locks          = "locks"
	Us             = "us"
	Unknown        = "unknown"
)

type Limit struct {
	SoftValue int
	HardValue int
	units     Unit
}

func makeUnit(str string) Unit {
	switch str {
	case "bytes":
		return Bytes
	case "seconds":
		return Seconds
	case "processes":
		return Processes
	case "signals":
		return Signals
	case "files":
		return Files
	case "us":
		return Us
	}
	return Unknown
}

//
// parseLimit parses a limit value from a string
// if the value is 'unlimited', it will be set to
// UNLIMITED (-1)
func parseLimit(str string) (int, error) {
	if str == "unlimited" {
		return UNLIMITED, nil
	}
	return strconv.Atoi(str)
}

var splitBy2Whitespace = regexp.MustCompile("\\s\\s+")

// logic taken from: https://github.com/etgryphon/stringUp
// no license is given, but i hope this is ok :(
var camelingRegex = regexp.MustCompile("[0-9A-Za-z]+")

//
// toCamelCase converts a spaced string to CamelCase
//
// e.g., "my string is" becomes "MyStringIs"
//
func toCamelCase(str string) string {
	byteSrc := []byte(str)
	chunks := camelingRegex.FindAll(byteSrc, -1)
	for idx, val := range chunks {
		if idx >= 0 {
			chunks[idx] = bytes.Title(val)
		}
	}
	return string(bytes.Join(chunks, nil))
}

//
// linesToLimits turns every line in in a /proc/pid/limits file
// to a Limit object.   
//
func linesToLimits(lines []string) (map[string]*Limit, error) {
	var result map[string]*Limit = make(map[string]*Limit)
	var units Unit

	// first line is the header
	for i := 1; i < len(lines); i++ {

		lines[i] = strings.TrimSpace(lines[i])

		if len(lines[i]) == 0 {
			// it's empty
			continue
		}

		parts := splitBy2Whitespace.Split(lines[i], UNLIMITED)

		if len(parts) < 3 {
			log.Println("malformed line, expected 4 parts but only got:", len(parts), "line:", lines[i])
			continue
		}

		title := strings.Replace(parts[0], "Max ", "", UNLIMITED)

		title = toCamelCase(title)

		soft, err := parseLimit(parts[1])
		if err != nil {
			return nil, err
		}
		hard, err := parseLimit(parts[2])
		if err != nil {
			return nil, err
		}

		if len(parts) > 3 {
			units = makeUnit(parts[3])
		} else {
			units = Unknown
		}

		result[title] = &Limit{
			SoftValue: soft,
			HardValue: hard,
			units:     units,
		}
	}
	return result, nil
}

//
// Limits is abstraction for /proc/<pid>/limit
//
// Each Limit pointer pointers to a Limit object, 
// which has the HardValue, SoftValue, and Units.
//
// If no limit is found, it will be nil. 
//
type Limits struct {
	CpuTime          *Limit
	FileSize         *Limit
	DataSize         *Limit
	StackSize        *Limit
	CoreFileSize     *Limit
	ResidentSet      *Limit
	Processes        *Limit
	OpenFiles        *Limit
	LockedMemory     *Limit
	AddressSpace     *Limit
	FileLocks        *Limit
	PendingSignals   *Limit
	MsgqueueSize     *Limit
	NicePriority     *Limit
	RealtimePriority *Limit
	RealtimeTimeout  *Limit
}

//
// Create a Limit instance from a /proc/<pid>/limits path
//
func New(path string) (*Limits, error) {

	buf, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	lines := strings.Split(string(buf), "\n")
	limits, err := linesToLimits(lines)
	if err != nil {
		return nil, err
	}

	limit := &Limits{}

	limit.CpuTime = limits["CpuTime"]
	limit.FileSize = limits["FileSize"]
	limit.DataSize = limits["DataSize"]
	limit.StackSize = limits["StackSize"]
	limit.CoreFileSize = limits["CoreFileSize"]
	limit.ResidentSet = limits["ResidentSet"]
	limit.Processes = limits["Processes"]
	limit.OpenFiles = limits["OpenFiles"]
	limit.LockedMemory = limits["LockedMemory"]
	limit.AddressSpace = limits["AddressSpace"]
	limit.FileLocks = limits["FileLocks"]
	limit.PendingSignals = limits["PendingSignals"]
	limit.MsgqueueSize = limits["MsgqueueSize"]
	limit.NicePriority = limits["NicePriority"]
	limit.RealtimePriority = limits["RealtimePriority"]
	limit.RealtimeTimeout = limits["RealtimeTimeout"]

	return limit, nil
}
