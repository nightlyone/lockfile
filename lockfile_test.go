package lockfile

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func ExampleLockfile() {

	file := filepath.Join(os.TempDir(), "lock.me.now.lck")
	lock, err := New(file)
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

	//create the lock
	file := filepath.Join(os.TempDir(), "lock.me.test")
	lock, err := New(file)
	if err != nil {
		t.Error(err)
	}

	//aquire the lock
	err = lock.TryLock()
	if err != nil {
		t.Error(err)
	}

	//close the lock
	defer func() {
		err = lock.Unlock()
		if err != nil {
			t.Error(err)
		}
	}()
}

func TestBadLock(t *testing.T) {

	file := filepath.Join(os.TempDir(), "lock.me.test")

	//write junk to the lock file
	if err := ioutil.WriteFile(file, []byte("asdf"), os.ModePerm); err != nil {
		t.Error(err)
	}

	//create the lock obj
	lock, err := New(file)
	if err != nil {
		t.Error(err)
	}

	//we expect this to detect its a bogus lock and aquire the lock anyway
	err = lock.TryLock()
	if err != nil {
		t.Error(err)
	}

	//succesfully close the lock
	defer func() {
		err = lock.Unlock()
		if err != nil {
			t.Error(err)
		}
	}()
}

func TestConcurrentAccess(t *testing.T) {

	//create the lock obj
	file := filepath.Join(os.TempDir(), "lock.me.test")
	lock, err := New(file)
	if err != nil {
		t.Error(err)
	}

	//aquire the lock
	err = lock.TryLock()
	if err != nil {
		t.Error(err)
	}

	//defer close the lock
	defer func() {
		err = lock.Unlock()
		if err != nil {
			t.Error(err)
		}
	}()

	//create another lock pointed to the same file
	lock2, err := New(file)
	if err != nil {
		t.Error(err)
	}

	//try to aquire the same lock, we expect this to fail
	err = lock2.TryLock()
	if err != ErrBusy {
		t.Error("expected lock to fail with ErrBusy")
	}
}
