package lockfile

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func ExampleLockfile() {
	lock, err := New("/tmp/lock.me.now.lck")
	if err != nil {
		fmt.Println("Cannot init lock. reason: %v", err)
		panic(err)
	}
	err = lock.TryLock()

	// Error handling is essential, as we only try to get the lock.
	if err != nil {
		fmt.Println("Cannot lock \"%v\", reason: %v", lock, err)
		panic(err)
	}

	defer lock.Unlock()

	fmt.Println("Do stuff under lock")
	// Output: Do stuff under lock
}

func SimpleLockTest(t *testing.T) {

	temp := os.TempDir()
	file := filepath.Join(temp, "lock.me.test")
	lock, err := New(file)

	if err != nil {
		t.Error(err)
	}

	err = lock.TryLock()
	if err != nil {
		t.Error(err)
	}

	defer func() {
		err = lock.Unlock()
		if err != nil {
			t.Error(err)
		}
	}()
}

func TestBadLock(t *testing.T) {

	temp := os.TempDir()
	file := filepath.Join(temp, "lock.me.test")
	if err := ioutil.WriteFile(file, []byte("asdf"), os.ModePerm); err != nil {
		t.Error(err)
	}

	lock, err := New(file)

	if err != nil {
		t.Error(err)
	}

	err = lock.TryLock()
	if err != nil {
		t.Error(err)
	}

	defer func() {
		err = lock.Unlock()
		if err != nil {
			t.Error(err)
		}
	}()
}

func TestConcurrentAccess(t *testing.T) {

	temp := os.TempDir()
	file := filepath.Join(temp, "lock.me.test")
	lock, err := New(file)

	if err != nil {
		t.Error(err)
	}

	err = lock.TryLock()
	if err != nil {
		t.Error(err)
	}

	defer func() {
		err = lock.Unlock()
		if err != nil {
			t.Error(err)
		}
	}()

	lock2, err := New(file)
	if err != nil {
		t.Error(err)
	}

	err = lock2.TryLock()
	if err == nil {
		t.Error("expected lock to fail")
	}
}
