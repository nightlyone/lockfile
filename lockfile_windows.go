package lockfile

import (
	"syscall"
)

func isRunning(pid int) (bool, error) {
	procHnd, err := syscall.OpenProcess(syscall.PROCESS_QUERY_INFORMATION, true, uint32(pid))
	if err != nil {
		if scerr, ok := err.(syscall.Errno); ok {
			if uintptr(scerr) == 87 {
				return false, nil //I only have a vague idea why this error occurs. I'm pretty sure it only occurs when the process isn't running #WindowsIsAPain
			}
		}
	}

	var code uint32
	err = syscall.GetExitCodeProcess(procHnd, &code)
	if err != nil {
		return false, err
	}

	return code == 259, nil
}
