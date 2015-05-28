package batten

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strconv"
	"strings"

	docker "github.com/fsouza/go-dockerclient"
	"github.com/jandre/procfs"
)

const (
	DockerUnixSocket = "unix:///var/run/docker.sock"
	DockerPidFile    = "/var/run/docker.pid"
)

func getDockerAPIConnection() (*docker.Client, error) {
	return docker.NewClient(DockerUnixSocket)
}

//
// PathExists return true if `filename` exists
//
func PathExists(filename string) bool {

	_, err := os.Stat(filename)

	if err != nil {
		if os.IsNotExist(err) {
			return false
		}
		panic(err)
	} else {
		return true
	}

}

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func runAuditCtl() (string, error) {
	cmdPath, err := exec.LookPath("auditctl")
	if err != nil || cmdPath == "" {
		return "", errors.New("Could not find auditctl tool. Do you have auditd installed?")
	}

	cmd := exec.Command(cmdPath, "-l")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", err
	}
	str := string(output)
	return str, nil
}

func pidOfDocker(dockerPidFile string) (int, error) {

	if dockerPidFile == "" {
		dockerPidFile = DockerPidFile
	}
	if !PathExists(dockerPidFile) {
		// no path exists
		return 0, nil
	}

	bytes, err := ioutil.ReadFile(dockerPidFile)

	if err != nil {
		// TODO: log error message
		return 0, err
	}
	pid, err := strconv.Atoi(string(bytes))

	procFile := fmt.Sprintf("/proc/%d", pid)

	if PathExists(procFile) {
		return pid, nil
	}
	return 0, nil
}

func getDockerProcess(dockerPidFile string) (*procfs.Process, error) {
	pid, err := pidOfDocker(dockerPidFile)
	if err != nil || pid <= 0 {
		return nil, err
	}

	return procfs.NewProcess(pid, true)
}

func readDockerDaemonEnviron(dockerPidFile string) (succ bool, environ map[string]string, err error) {

	pid, err := pidOfDocker(dockerPidFile)
	if err != nil {
		// TODO: log error message
		return false, environ, err
	}
	if pid <= 0 {
		// no pid was found, but not really an error
		return true, environ, nil
	}

	process, err := procfs.NewProcess(pid, true)
	if err != nil {
		return false, environ, err
	}

	environ = process.Environ
	return true, environ, nil
}

func readDockerDaemonArgs(dockerPidFile string) (succ bool, cmdLine []string, err error) {

	pid, err := pidOfDocker(dockerPidFile)
	if err != nil {
		// TODO: log error message
		return false, cmdLine, err
	}
	if pid <= 0 {
		// no pid was found, but not really an error
		return true, cmdLine, nil
	}
	process, err := procfs.NewProcess(pid, true)
	if err != nil {
		return false, cmdLine, err
	}

	cmdLine = process.Cmdline
	return true, cmdLine, nil
}

func getArgValue(lookFor string, args []string) string {
	prev := ""
	for _, arg := range args {
		arg = strings.TrimSpace(arg)
		if arg != "" {
			if prev == lookFor {
				return arg
			} else if strings.HasPrefix(arg, lookFor+"=") {
				filename := arg[len(lookFor)+1:]
				return filename
			}
			prev = arg

		}
	}

	return ""
}
