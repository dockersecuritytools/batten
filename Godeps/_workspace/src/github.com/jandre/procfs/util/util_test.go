package util

import (
	"testing"
)

func TestSystemStart(t *testing.T) {
	systemStart()
}

func TestJiffiesToEpoch(t * testing.T) {
	val := JiffiesToEpoch(17350497)
	if (val < GLOBAL_SYSTEM_START) {
		t.Error("error", val)
	}
}
