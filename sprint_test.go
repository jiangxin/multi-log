package log

import (
	"bytes"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func demoNotes(o Options) {
	Init(o)
	i := 1

	Notef("note #%d", i)
	i++
	Note("note #", i)
	i++
	Noteln("note #", i)
	i++
	Printf("note #%d", i)
	i++
	Print("note #", i)
	i++
	Println("note #", i)
}

func TestNoteNotSaveLogfile(t *testing.T) {
	var (
		assert = assert.New(t)
		buffer bytes.Buffer
		err    error
		expect string
	)

	tmpdir, err := ioutil.TempDir("", "multi-logger-")
	if err != nil {
		panic(err)
	}
	defer func(dir string) {
		os.RemoveAll(dir)
	}(tmpdir)

	tmpLog := filepath.Join(tmpdir, "log.txt")

	o := Options{
		Quiet:    false,
		LogLevel: "error",
		LogFile:  tmpLog,
		stderr:   &buffer,
	}

	demoNotes(o)

	expect = `NOTE: note #1
NOTE: note #2
NOTE: note # 3
NOTE: note #4
NOTE: note #5
NOTE: note # 6` + "\n"
	assert.Equal(expect, buffer.String())

	expect = ""
	data, err := ioutil.ReadFile(tmpLog)
	assert.Nil(err)
	assert.Equal(expect, filterTime(string(data)))
}

func TestNoteQuiet(t *testing.T) {
	var (
		assert = assert.New(t)
		buffer bytes.Buffer
		err    error
		expect string
	)

	tmpdir, err := ioutil.TempDir("", "multi-logger-")
	if err != nil {
		panic(err)
	}
	defer func(dir string) {
		os.RemoveAll(dir)
	}(tmpdir)

	tmpLog := filepath.Join(tmpdir, "log.txt")

	o := Options{
		Quiet:    true,
		Verbose:  1,
		LogLevel: "error",
		LogFile:  tmpLog,
		stderr:   &buffer,
	}

	demoNotes(o)

	expect = ""
	assert.Equal(expect, buffer.String())

	expect = ""
	data, err := ioutil.ReadFile(tmpLog)
	assert.Nil(err)
	assert.Equal(expect, filterTime(string(data)))
}
