package cadvisor

import (
	"github.com/rancher/agent/utilities/config"
	"github.com/rancher/agent/utilities/utils"
	"os"
	"os/exec"
	"syscall"
)

func StartUp() error {
	args := []string{
		"cadvisor",
		"-logtostderr=true",
		"-listen_ip", config.CadvisorIP(),
		"-port", config.CadvisorPort(),
		"-housekeeping_interval", config.CadvisorInterval(),
	}
	dockerRoot := config.CadvisorDockerRoot()
	if len(dockerRoot) > 0 {
		args = append(args, []string{"-docker_root", dockerRoot}...)
	}
	cadvisorOpts := config.CadvisorOpts()
	if len(cadvisorOpts) > 0 {
		args = append(args, utils.SafeSplit(cadvisorOpts)...)
	}
	wrapper := config.CadvisorWrapper()
	if len(wrapper) > 0 {
		args = append([]string{wrapper}, args...)
	} else if _, err := os.Stat("/host/proc/1/ns/mnt"); err == nil {
		args = append([]string{"nsenter", "--mount=/host/proc/1/ns/mnt", "--"}, args...)
	}
	command := exec.Command(args[0], args[1:len(args)]...)
	command.SysProcAttr = &syscall.SysProcAttr{
		Setpgid: true,
	}
	command.Stderr = os.Stderr
	command.Stdout = os.Stdout
	command.Start()
	err := command.Wait()
	return err
}