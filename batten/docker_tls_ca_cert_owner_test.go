package batten

import (
	"io/ioutil"
	"os"
	"os/user"
	"strconv"
	"testing"
)

func TestLookForTLSConfigString(t *testing.T) {

	expected := "/tmp/test_certfile"
	check := makeDockerTLSCACertOwnerCheck()
	dc := check.(*DockerTLSCACertOwnerCheck)
	currUser, _ := user.Current()

	uid, _ := strconv.Atoi(currUser.Uid)
	gid, _ := strconv.Atoi(currUser.Gid)
	dc.uid = uint32(uid)
	dc.gid = uint32(gid)

	ioutil.WriteFile(expected, []byte("hello"), 0644)
	defer os.Remove(expected)

	args := []string{
		"docker",
		"--tlscacert=" + expected,
	}
	found, err := dc.validate(args)
	if !found {
		t.Fatal("Expected owner to be current user: "+expected, err)
	}
	// check that bad uids fail
	dc.uid = 0
	dc.gid = 0
	found, err = dc.validate(args)
	if found {
		t.Fatal("Expected owner to not be root:"+expected, err)
	}

}
