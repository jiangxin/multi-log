package path

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExpendHome(t *testing.T) {
	var (
		home   string
		tmpdir string
		name   string
		err    error
		assert = assert.New(t)
	)

	tmpdir, err = ioutil.TempDir("", "goconfig")
	if err != nil {
		panic(err)
	}
	defer func(dir string) {
		os.RemoveAll(dir)
	}(tmpdir)

	home, err = HomeDir()
	assert.Nil(err)
	defer func(home string) {
		SetHome(home)
	}(home)

	UnsetHome()
	name, err = HomeDir()
	assert.NotNil(err)
	assert.Equal("", name)

	name, err = ExpendHome("")
	assert.NotNil(err)
	assert.Equal("", name)

	SetHome(tmpdir)

	name, err = HomeDir()
	assert.Equal(tmpdir, name)

	name, err = ExpendHome("")
	assert.Nil(err)
	assert.Equal(tmpdir, name)

	name, err = ExpendHome("a")
	assert.Nil(err)
	assert.Equal(filepath.Join(tmpdir, "a"), name)

	name, err = ExpendHome("~a")
	assert.Nil(err)
	assert.Equal(filepath.Join(tmpdir, "~a"), name)

	name, err = ExpendHome("~")
	assert.Nil(err)
	assert.Equal(tmpdir, name)

	name, err = ExpendHome("~/")
	assert.Nil(err)
	assert.Equal(tmpdir, name)

	name, err = ExpendHome("~/a")
	assert.Nil(err)
	assert.Equal(filepath.Join(tmpdir, "a"), name)

	name, err = ExpendHome("ab")
	assert.Nil(err)
	assert.Equal(filepath.Join(tmpdir, "ab"), name)

	inputdir := "/"
	if runtime.GOOS == "windows" {
		inputdir = "c:\\"
	}
	name, err = ExpendHome(inputdir)
	assert.Nil(err)
	assert.Equal(inputdir, name)

	inputdir = "/a"
	if runtime.GOOS == "windows" {
		inputdir = "c:\\a"
	}
	name, err = ExpendHome(inputdir)
	assert.Nil(err)
	assert.Equal(inputdir, name)

}

func TestAbs(t *testing.T) {
	var (
		home   string
		tmpdir string
		name   string
		err    error
		assert = assert.New(t)
	)

	tmpdir, err = ioutil.TempDir("", "goconfig")
	if err != nil {
		panic(err)
	}
	defer func(dir string) {
		os.RemoveAll(dir)
	}(tmpdir)

	home, err = HomeDir()
	assert.Nil(err)
	defer func(home string) {
		SetHome(home)
	}(home)

	UnsetHome()
	name, err = Abs("~/")
	assert.NotNil(err)
	assert.Equal("", name)

	SetHome(tmpdir)
	cwd, err := os.Getwd()
	assert.Nil(err)

	name, err = Abs("")
	assert.Nil(err, fmt.Sprintf("err should be nil, but got: %s", err))
	assert.Equal(cwd, name)

	name, err = Abs("a")
	assert.Nil(err)
	assert.Equal(filepath.Join(cwd, "a"), name)

	name, err = Abs("~a")
	assert.Nil(err)
	assert.Equal(filepath.Join(cwd, "~a"), name)

	name, err = Abs("~")
	assert.Nil(err)
	assert.Equal(tmpdir, name)

	name, err = Abs("~/")
	assert.Nil(err)
	assert.Equal(tmpdir, name)

	name, err = Abs("~/a")
	assert.Nil(err)
	assert.Equal(filepath.Join(tmpdir, "a"), name)

	name, err = Abs("ab")
	assert.Nil(err)
	assert.Equal(filepath.Join(cwd, "ab"), name)

	inputdir := "/"
	if runtime.GOOS == "windows" {
		inputdir = "c:\\"
	}
	name, err = Abs(inputdir)
	assert.Nil(err)
	assert.Equal(inputdir, name)

	inputdir = "/a"
	if runtime.GOOS == "windows" {
		inputdir = "c:\\a"
	}
	name, err = Abs(inputdir)
	assert.Nil(err)
	assert.Equal(inputdir, name)
}

func TestAbsJoin(t *testing.T) {
	var (
		home   string
		tmpdir string
		name   string
		err    error
		assert = assert.New(t)
	)

	tmpdir, err = ioutil.TempDir("", "goconfig")
	if err != nil {
		panic(err)
	}
	defer func(dir string) {
		os.RemoveAll(dir)
	}(tmpdir)

	home, err = HomeDir()
	assert.Nil(err)
	defer func(home string) {
		SetHome(home)
	}(home)

	SetHome(tmpdir)

	cwd := "/some/dir"
	if runtime.GOOS == "windows" {
		cwd = "c:\\some\\dir"
	}

	name, err = AbsJoin(cwd, "")
	assert.Nil(err)
	assert.Equal(cwd, name)

	name, err = AbsJoin(cwd, "a")
	assert.Nil(err)
	assert.Equal(filepath.Join(cwd, "a"), name)

	name, err = AbsJoin(cwd, "~a")
	assert.Nil(err)
	assert.Equal(filepath.Join(cwd, "~a"), name)

	name, err = AbsJoin(cwd, "~")
	assert.Nil(err)
	assert.Equal(tmpdir, name)

	name, err = AbsJoin(cwd, "~/")
	assert.Nil(err)
	assert.Equal(tmpdir, name)

	name, err = AbsJoin(cwd, "~/a")
	assert.Nil(err)
	assert.Equal(filepath.Join(tmpdir, "a"), name)

	name, err = AbsJoin(cwd, "ab")
	assert.Nil(err)
	assert.Equal(filepath.Join(cwd, "ab"), name)

	inputdir := "/"
	if runtime.GOOS == "windows" {
		inputdir = "c:\\"
	}
	name, err = AbsJoin(cwd, inputdir)
	assert.Nil(err)
	assert.Equal(inputdir, name)

	inputdir = "/a"
	if runtime.GOOS == "windows" {
		inputdir = "c:\\a"
	}
	name, err = AbsJoin(cwd, inputdir)
	assert.Nil(err)
	assert.Equal(inputdir, name)
}
