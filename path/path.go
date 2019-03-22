package path

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

// HomeDir returns home directory
func HomeDir() (string, error) {
	var (
		home string
	)

	if runtime.GOOS == "windows" {
		home = os.Getenv("USERPROFILE")
		if home == "" {
			home = os.Getenv("HOMEDRIVE") + os.Getenv("HOMEPATH")
		}
	}
	if home == "" {
		home = os.Getenv("HOME")
	}

	if home == "" {
		return "", fmt.Errorf("cannot find HOME")
	}

	return home, nil
}

// ExpendHome expends path prefix "~/" to home dir
func ExpendHome(name string) (string, error) {
	if filepath.IsAbs(name) {
		return name, nil
	}

	home, err := HomeDir()
	if err != nil {
		return "", err
	}

	if len(name) == 0 || name == "~" {
		return home, nil
	} else if len(name) > 1 && name[0] == '~' && (name[1] == '/' || name[1] == '\\') {
		return filepath.Join(home, name[2:]), nil
	}

	return filepath.Join(home, name), nil
}

// Abs returns absolute path and will expend homedir if path has "~/' prefix
func Abs(name string) (string, error) {
	if name == "" {
		return os.Getwd()
	}

	if filepath.IsAbs(name) {
		return name, nil
	}

	if len(name) > 0 && name[0] == '~' && (len(name) == 1 || name[1] == '/' || name[1] == '\\') {
		return ExpendHome(name)
	}

	return filepath.Abs(name)
}

// AbsJoin returns absolute path, and use <dir> as parent dir for relative path
func AbsJoin(dir, name string) (string, error) {
	if name == "" {
		return filepath.Abs(dir)
	}

	if filepath.IsAbs(name) {
		return name, nil
	}

	if len(name) > 0 && name[0] == '~' && (len(name) == 1 || name[1] == '/' || name[1] == '\\') {
		return ExpendHome(name)
	}

	return Abs(filepath.Join(dir, name))
}

// UnsetHome unsets HOME related environments
func UnsetHome() {
	if runtime.GOOS == "windows" {
		os.Unsetenv("USERPROFILE")
		os.Unsetenv("HOMEDRIVE")
		os.Unsetenv("HOMEPATH")
	}
	os.Unsetenv("HOME")
}

// SetHome sets proper HOME environments
func SetHome(home string) {
	if runtime.GOOS == "windows" {
		os.Setenv("USERPROFILE", home)
		if strings.Contains(home, ":\\") {
			slices := strings.SplitN(home, ":\\", 2)
			if len(slices) == 2 {
				os.Setenv("HOMEDRIVE", slices[0]+":")
				os.Setenv("HOMEPATH", "\\"+slices[1])
			}
		}
	} else {
		os.Setenv("HOME", home)
	}
}
