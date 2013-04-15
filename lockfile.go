// Handle filename based locking.
package lockfile

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"syscall"
)

type Lockfile string

var (
	ErrBusy        = errors.New("Locked by other process")
	ErrNeedAbsPath = errors.New("Lockfiles must be given as absolute path names")
)

// ugly workaround, because os.FindProcess() is crap
func findProcess(pid int) error {
	p, err := os.FindProcess(pid)
	if err != nil {
		return err
	}
	return p.Signal(os.Signal(syscall.Signal(0)))
}

// Describe a new filename located at path. It is expected to be an absolute path
func New(path string) (Lockfile, error) {
	if !filepath.IsAbs(path) {
		return Lockfile(""), ErrNeedAbsPath
	}
	return Lockfile(path), nil
}

// Try to get Lockfile lock. Returns nil, if successful and and error describing the reason, it didn't work out.
// Please note, that existing lockfiles containing pids of dead processes and lockfiles containing no pid at all
// are deleted.
func (l Lockfile) TryLock() error {
	name := string(l)

	// This has been checked by New already. If we trigger here, 
        // the caller didn't use New and re-implemented it's functionality badly. 
        // So panic, that he might find this easily during testing.
	if !filepath.IsAbs(string(name)) {
		panic(ErrNeedAbsPath)
	}

	tmplock, err := ioutil.TempFile(filepath.Dir(name), "")
	if err != nil {
		return err
	} else {
		defer tmplock.Close()
		defer os.Remove(tmplock.Name())
	}

	_, err = tmplock.WriteString(fmt.Sprintf("%d\n", os.Getpid()))
	if err != nil {
		return err
	}

	// return value intentionally ignored, as ignoring it is part of the algorithm
	_ = os.Link(tmplock.Name(), name)

	fiTmp, err := os.Lstat(tmplock.Name())
	if err != nil {
		return err
	}
	fiLock, err := os.Lstat(name)
	if err != nil {
		return err
	}

	// Success
	if os.SameFile(fiTmp, fiLock) {
		return nil
	}

	// Ok, see, if we have a stale lockfile here
	content, err := ioutil.ReadFile(name)

	var pid int
	_, err = fmt.Sscanln(string(content), &pid)
	if err != nil {
		return err
	}

	// try hard for pids. If no pid, the lockfile is junk anyway and we delete it.
	if pid != 0 {
		err = findProcess(pid)
		if err == nil {
			return ErrBusy
		}
		if errno, ok := err.(syscall.Errno); ok && errno != syscall.ESRCH {
			return ErrBusy
		}
	}

	// clean stale/invalid lockfile
	err = os.Remove(name)
	if err != nil {
		return err
	}

	// now that we cleaned up the stale lockfile, let's recurse
	return l.TryLock()
}

// Release a lock again. Returns any error that happend during release
func (l Lockfile) Unlock() error {
	return os.Remove(string(l))
}
