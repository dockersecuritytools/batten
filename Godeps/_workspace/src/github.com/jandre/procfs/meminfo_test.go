package procfs

import (
	"testing"
)

func TestParseMeminfo(t *testing.T) {
	// set the GLOBAL_SYSTEM_START
	meminfo, err :=  ParseMeminfo("./testfiles/meminfo")

	if err != nil {
		t.Fatal("Got error", err)
	}

	if meminfo == nil {
		t.Fatal("meminfo is missing")
	}

	if meminfo.MemTotal != 1011932 {
		t.Fatal("Expected 1011932 from MemTotal")
	}

	if meminfo.PageTables != 8340 {
		t.Fatal("Expected 8340 from PageTables")
	}
}

