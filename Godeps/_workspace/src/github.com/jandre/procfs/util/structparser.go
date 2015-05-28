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
// WITH THE SOFTWARE OR THE USE OR O

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"
)

//
// ParseStringsIntoStruct expects a pointer to a struct as its
// first argument. It assigns each element from the lines slice
// sequentially to the struct members, parsing each according to
// type. It currently accepts fields of type int, int64, string
// and time.Time (it assumes that values of the latter kind
// are formatted as a clock-tick count since the system start).
//
// Extra lines are ignored.
//
func ParseStringsIntoStruct(vi interface{}, strs []string) error {
	v := reflect.ValueOf(vi).Elem()
	typeOf := v.Type()

	for i := 0; i < v.NumField(); i++ {
		if i > len(strs) {
			break
		}
		str := strings.TrimSpace(strs[i])
		interf := v.Field(i).Addr().Interface()
		if err := parseField(interf, str); err != nil {
			return fmt.Errorf("cannot parse field %s=%q: %v", typeOf.Field(i).Name, str, err)
		}
	}
	return nil
}

func parseField(field interface{}, line string) error {
	switch field := field.(type) {
	case *int:
		val, err := strconv.Atoi(line)
		if err != nil {
			return err
		}
		*field = val
	case *int64:
		val, err := strconv.ParseInt(line, 10, 64)
		if err != nil {
			return err
		}
		*field = val
	case *uint64:
		val, err := strconv.ParseUint(line, 10, 64)
		if err != nil {
			return err
		}
		*field = val
	case *string:
		*field = line
	case *time.Time:
		jiffies, err := strconv.ParseInt(line, 10, 64)
		if err != nil {
			return err
		}
		*field = jiffiesToTime(jiffies)
	default:
		return fmt.Errorf("unsupported field type %T", field)
	}
	return nil
}
