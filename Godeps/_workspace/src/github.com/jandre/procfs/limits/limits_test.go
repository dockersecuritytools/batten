package limits

import (
	"testing"
)

func TestParsingLimit(t *testing.T) {
	s, err := New("./testfiles/limits")

	if err != nil {
		t.Fatal("Got error", err)
	}

	if s == nil {
		t.Fatal("stat is missing")
	}

	if s.NicePriority == nil || s.NicePriority.SoftValue != 0 || s.NicePriority.HardValue != 0 {
		t.Fatal("Wrong values for NicePriority, expected soft=0 and hard=4096", s.NicePriority)
	}

	if s.OpenFiles == nil || s.OpenFiles.SoftValue != 1024 || s.OpenFiles.HardValue != 4096 {
		t.Fatal("Wrong values for OpenFiles, expected soft=1024 and hard=4096", s.OpenFiles)
	}

}
