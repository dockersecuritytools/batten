package procfs

import (
	"testing"
	"log"
	"fmt"
)

func TestAllProc(t * testing.T) {
	procs, err := Processes(false)
	if err != nil {
		t.Fatal(err)
	}
	if len(procs) <= 0 {
		t.Fatal("procs length must be > 0")
	}

	log.Println("Pid 1", procs[0])
	for i, p := range procs {
		fmt.Printf("%d PID: %d - CMDLINE: %s - CWD: %s - EXE: %s\n", i, p.Pid, p.Cmdline, p.Cwd, p.Exe)
		// noop
	}

	// for i := 0; i < len(procs); i++ {
		// p := procs[i];
		// log.Println("%d PID: %d - CMDLINE: %s - CWD: %s - EXE: %s\n", p.Pid, p.Cmdline, p.Cwd, p.Exe)
	// }

}
