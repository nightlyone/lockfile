package lockfile

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"path/filepath"
	"strconv"
	"testing"
)

func TestLock(t *testing.T) {
	path, err := filepath.Abs("test_lockfile.pid")
	if err != nil {
		panic(err)
	}

	lf, err := New(path)
	if err != nil {
		t.Fail()
		fmt.Println("Error making lockfile: ", err)
		return
	}

	err = lf.TryLock()
	if err != nil {
		t.Fail()
		fmt.Println("Error locking lockfile: ", err)
		return
	}

	err = lf.Unlock()
	if err != nil {
		t.Fail()
		fmt.Println("Error unlocking lockfile: ", err)
		return
	}
}

func TestDeadPID(t *testing.T) {
	path, err := filepath.Abs("test_lockfile.pid")
	if err != nil {
		panic(err)
	}

	pid := GetDeadPID()

	ioutil.WriteFile(path, []byte(strconv.Itoa(pid)+"\n"), 0666)
}

func GetDeadPID() int {
	for {
		pid := rand.Int() % 4096 //I have no idea how windows handles large PIDs, or if they even exist. Limit it to 4096 to be safe.
		running, err := isRunning(pid)
		if err != nil {
			fmt.Println("Error checking PID: ", err)
			continue
		}

		if !running {
			return pid
		}
	}
}

func TestBusy(t *testing.T) {
	path, err := filepath.Abs("test_lockfile.pid")
	if err != nil {
		t.Fail()
		fmt.Println(err)
		return
	}

	lf1, err := New(path)
	if err != nil {
		t.Fail()
		fmt.Println(err)
		return
	}

	err = lf1.TryLock()
	if err != nil {
		t.Fail()
		fmt.Println(err)
		return
	}

	lf2, err := New(path)
	if err != nil {
		t.Fail()
		fmt.Println(err)
		return
	}

	err = lf2.TryLock()
	if err == nil {
		t.Fail()
		fmt.Println("No error locking already locked lockfile!")
		return
	} else if err != ErrBusy {
		t.Fail()
		fmt.Println(err)
		return
	}

	err = lf1.Unlock()
	if err != nil {
		t.Fail()
		fmt.Println(err)
		return
	}
}
