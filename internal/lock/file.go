package lock

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"syscall"
)

var ErrHeld = errors.New("dbinstall host lock is already held")

type File struct{ file *os.File }

func Acquire(stateDirectory string) (*File, error) {
	if err := os.MkdirAll(stateDirectory, 0o750); err != nil {
		return nil, fmt.Errorf("create state directory: %w", err)
	}
	file, err := os.OpenFile(filepath.Join(stateDirectory, "host.lock"), os.O_CREATE|os.O_RDWR, 0o640)
	if err != nil {
		return nil, fmt.Errorf("open host lock: %w", err)
	}
	if err := syscall.Flock(int(file.Fd()), syscall.LOCK_EX|syscall.LOCK_NB); err != nil {
		file.Close()
		if errors.Is(err, syscall.EWOULDBLOCK) {
			return nil, ErrHeld
		}
		return nil, fmt.Errorf("acquire host lock: %w", err)
	}
	return &File{file: file}, nil
}

func (f *File) Close() error {
	if f == nil || f.file == nil {
		return nil
	}
	unlock := syscall.Flock(int(f.file.Fd()), syscall.LOCK_UN)
	closeErr := f.file.Close()
	if unlock != nil {
		return fmt.Errorf("release host lock: %w", unlock)
	}
	return closeErr
}
