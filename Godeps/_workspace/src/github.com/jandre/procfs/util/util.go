package util

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
	"strconv"
	"strings"
	"time"
)

/*
#include <unistd.h>
*/
import "C"

//
// systemStart() attempts to detect the system start time. 
//
// It will attempt to fetch the system start time by parsing /proc/stat
// look for the btime line, and converting that to an int64 value (the start
// time in Epoch seconds).
//
// This value must be used for certain /proc time values that are specified
// in 'jiffies', and is used with the _SC_CLK_TCK retrieved from sysconf(3)
// to calculate epoch timestamps.
//
//
func systemStart() int64 {
	str, err := ioutil.ReadFile("/proc/stat")
	if err != nil {
		log.Fatal("Unable to read btime from /proc/stat - is this Linux?", err)
	}
	lines := strings.Split(string(str), "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "btime") {
			parts := strings.Split(line, " ")
			// format is btime 1388417200
			val, err := strconv.ParseInt(parts[1], 10, 32)
			if err != nil {
				log.Fatal("Unable to convert btime value in /proc/stat to int64", parts[1], err)
			}
			return int64(val)
		}
	}

	log.Fatal("No btime found in /proc/stat.  This value is needed to calculate timestamps")

	return 0
}

var GLOBAL_SYSTEM_START int64 = systemStart()

//
// jiffiesToTime converts jiffies to a Time object
// using the GLOBAL_SYSTEM_START time and the value of the
// _SC_CLK_TICK value obtained from sysconf(3).
//
// To get the # of seconds elapsed since system start, we do jiffies / _SC_CLK_TCK.
//
// We then add the system start time (GLOBAL_SYSTEM_START) to get the epoch
// timestamp in seconds.
//
func jiffiesToTime(jiffies int64) time.Time {
	ticks := C.sysconf(C._SC_CLK_TCK)
	return time.Unix(GLOBAL_SYSTEM_START + jiffies/int64(ticks), 0)
}

