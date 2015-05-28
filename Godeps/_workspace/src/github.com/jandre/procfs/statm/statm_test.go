package statm

import (
	// "procfs/util"
	"log"
	"testing"
)

func TestParsingStatm(t *testing.T) {
	// set the GLOBAL_SYSTEM_START
	s, err := New("./testfiles/statm")

	if err != nil {
		t.Fatal("Got error", err)
	}

	if s == nil {
		t.Fatal("statm is missing")
	}
	log.Println("statm", s)

	if s.Size != 134008 {
		t.Fatal("Expected size to be 134008")
	}
	if s.Resident != 72921 {
		t.Fatal("Expected Resident to be 72921")
	}
}
