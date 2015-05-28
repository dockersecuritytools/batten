package procfs

import (
	"log"
	"testing"
)

func TestParsingInit(t *testing.T) {
	// set the GLOBAL_SYSTEM_START
	process, err := NewProcess(1, false)

	if err != nil {
		t.Fatal("Got error", err)
	}

	if process == nil {
		t.Fatal("process is missing")
	}

	if process.Cwd != "/" {
		t.Fatal("Expected / for cwd")
	}
}

func TestParsingCmdline(t *testing.T) {
	prefix := "./testfiles/"
	cmd, err := readCmdLine(prefix)

	if err != nil {
		t.Fatal("Got error", err)
	}

	log.Println("Got cmdline:", cmd, len(cmd))
	if len(cmd) != 9 {
		t.Fatal("Expected string length to be 9 for cmdline, got", len(cmd))
	}

	if cmd[0] != "dhclient3" {
		t.Fatal("Expected dhclient3 for cmdline argv0, got", cmd[0])
	}
}

