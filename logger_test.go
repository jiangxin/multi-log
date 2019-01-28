package log

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func testingExitFunc(code int) {
	fmt.Fprintf(os.Stderr, "will exit %d\n", code)
}

func filterTime(data string) string {
	res := []*regexp.Regexp{
		regexp.MustCompile(`(\d{4}-\d{1,2}-\d{1,2}T\d{2}:\d{2}:\d{2}[\d:+-]*)`),
	}
	for _, re := range res {
		data = re.ReplaceAllString(data, "<time>")
	}
	return data
}

func demoLoggerf(o Options) {
	o.exitFunc = testingExitFunc
	Init(o)
	i := 1

	Tracef("trace #%d", i)
	i++
	Debugf("debug #%d", i)
	i++
	Infof("info #%d", i)
	i++
	Warnf("warn #%d", i)
	i++
	Warningf("warning #%d", i)
	i++
	Errorf("error #%d", i)
	i++
	Fatalf("fatal #%d", i)
}

func demoLogger(o Options) {
	o.exitFunc = testingExitFunc
	Init(o)
	i := 1

	Trace("trace #", i)
	i++
	Debug("debug #", i)
	i++
	Info("info #", i)
	i++
	Warn("warn #", i)
	i++
	Warning("warning #", i)
	i++
	Error("error #", i)
	i++
	Fatal("fatal #", i)
}

func demoLoggerln(o Options) {
	o.exitFunc = testingExitFunc
	Init(o)
	i := 1

	Traceln("trace #", i)
	i++
	Debugln("debug #", i)
	i++
	Infoln("info #", i)
	i++
	Warnln("warn #", i)
	i++
	Warningln("warning #", i)
	i++
	Errorln("error #", i)
	i++
	Fatalln("fatal #", i)
}

func TestLoggerNoLogfile(t *testing.T) {
	var (
		assert = assert.New(t)
		buffer bytes.Buffer
		err    error
		expect string
	)

	tmpdir, err := ioutil.TempDir("", "multi-logger-")
	assert.Nil(err)
	defer func(dir string) {
		os.RemoveAll(dir)
	}(tmpdir)

	o := Options{
		Verbose:     4,
		LogLevel:    "warning",
		LogFile:     "",
		stderr:      &buffer,
		ForceColors: true,
	}

	demoLogger(o)

	expect = "" +
		"\x1b[1;34mTRACE:\x1b[0m trace #1\n" +
		"\x1b[1;34mDEBUG:\x1b[0m debug #2\n" +
		"\x1b[1;36mINFO:\x1b[0m info #3\n" +
		"\x1b[1;33mWARNING:\x1b[0m warn #4\n" +
		"\x1b[1;33mWARNING:\x1b[0m warning #5\n" +
		"\x1b[1;31mERROR:\x1b[0m error #6\n" +
		"\x1b[1;31mFATAL:\x1b[0m fatal #7\n"
	assert.Equal(expect, buffer.String())
}

func TestRelativeLoggerFile(t *testing.T) {
	var (
		assert = assert.New(t)
		home   string
		buffer bytes.Buffer
		err    error
		expect string
	)

	tmpdir, err := ioutil.TempDir("", "multi-logger-")
	assert.Nil(err)
	defer func(dir string) {
		os.RemoveAll(dir)
	}(tmpdir)

	home = os.Getenv("HOME")
	os.Setenv("HOME", tmpdir)
	defer os.Setenv("HOME", home)

	tmpfile := "~/log/log.txt"

	o := Options{
		Verbose:  4,
		LogLevel: "warning",
		LogFile:  tmpfile,
		stderr:   &buffer,
	}

	demoLogger(o)

	expect = `WARN[<time>]: warn #4
WARN[<time>]: warning #5
ERRO[<time>]: error #6
FATA[<time>]: fatal #7
`
	relTmpfile := filepath.Join(tmpdir, "log", "log.txt")
	data, err := ioutil.ReadFile(relTmpfile)
	assert.Nil(err)
	assert.Equal(expect, filterTime(string(data)))
}

func TestLoggerfDefault(t *testing.T) {
	var (
		assert = assert.New(t)
		buffer bytes.Buffer
		err    error
		expect string
	)

	tmpdir, err := ioutil.TempDir("", "multi-logger-")
	assert.Nil(err)
	defer func(dir string) {
		os.RemoveAll(dir)
	}(tmpdir)

	tmpLog := filepath.Join(tmpdir, "log.txt")

	o := Options{
		LogFile: tmpLog,
		stderr:  &buffer,
	}

	demoLoggerf(o)

	expect = `WARNING: warn #4
WARNING: warning #5
ERROR: error #6
FATAL: fatal #7
`
	assert.Equal(expect, buffer.String())

	expect = `WARN[<time>]: warn #4
WARN[<time>]: warning #5
ERRO[<time>]: error #6
FATA[<time>]: fatal #7
`
	data, err := ioutil.ReadFile(tmpLog)
	assert.Nil(err)
	assert.Equal(expect, filterTime(string(data)))
}

func TestLoggerfCustom(t *testing.T) {
	var (
		assert = assert.New(t)
		buffer bytes.Buffer
		err    error
		expect string
	)

	tmpdir, err := ioutil.TempDir("", "multi-logger-")
	assert.Nil(err)
	defer func(dir string) {
		os.RemoveAll(dir)
	}(tmpdir)

	tmpLog := filepath.Join(tmpdir, "log.txt")

	o := Options{
		Verbose:  1,
		LogLevel: "warn",
		LogFile:  tmpLog,
		stderr:   &buffer,
	}

	demoLoggerf(o)

	expect = `INFO: info #3
WARNING: warn #4
WARNING: warning #5
ERROR: error #6
FATAL: fatal #7
`
	assert.Equal(expect, buffer.String())

	expect = `WARN[<time>]: warn #4
WARN[<time>]: warning #5
ERRO[<time>]: error #6
FATA[<time>]: fatal #7
`
	data, err := ioutil.ReadFile(tmpLog)
	assert.Nil(err)
	assert.Equal(expect, filterTime(string(data)))
}

func TestLoggerDefault(t *testing.T) {
	var (
		assert = assert.New(t)
		buffer bytes.Buffer
		err    error
		expect string
	)

	tmpdir, err := ioutil.TempDir("", "multi-logger-")
	assert.Nil(err)
	defer func(dir string) {
		os.RemoveAll(dir)
	}(tmpdir)

	tmpLog := filepath.Join(tmpdir, "log.txt")

	o := Options{
		LogFile: tmpLog,
		stderr:  &buffer,
	}

	demoLogger(o)

	expect = `WARNING: warn #4
WARNING: warning #5
ERROR: error #6
FATAL: fatal #7
`
	assert.Equal(expect, buffer.String())

	expect = `WARN[<time>]: warn #4
WARN[<time>]: warning #5
ERRO[<time>]: error #6
FATA[<time>]: fatal #7
`
	data, err := ioutil.ReadFile(tmpLog)
	assert.Nil(err)
	assert.Equal(expect, filterTime(string(data)))
}

func TestLoggerCustom(t *testing.T) {
	var (
		assert = assert.New(t)
		buffer bytes.Buffer
		err    error
		expect string
	)

	tmpdir, err := ioutil.TempDir("", "multi-logger-")
	assert.Nil(err)
	defer func(dir string) {
		os.RemoveAll(dir)
	}(tmpdir)

	tmpLog := filepath.Join(tmpdir, "log.txt")

	o := Options{
		Verbose:  2,
		LogLevel: "info",
		LogFile:  tmpLog,
		stderr:   &buffer,
	}

	demoLogger(o)

	expect = `DEBUG: debug #2
INFO: info #3
WARNING: warn #4
WARNING: warning #5
ERROR: error #6
FATAL: fatal #7
`
	assert.Equal(expect, buffer.String())

	expect = `INFO[<time>]: info #3
WARN[<time>]: warn #4
WARN[<time>]: warning #5
ERRO[<time>]: error #6
FATA[<time>]: fatal #7
`
	data, err := ioutil.ReadFile(tmpLog)
	assert.Nil(err)
	assert.Equal(expect, filterTime(string(data)))
}

func TestLoggerlnDefault(t *testing.T) {
	var (
		assert = assert.New(t)
		buffer bytes.Buffer
		err    error
		expect string
	)

	tmpdir, err := ioutil.TempDir("", "multi-logger-")
	assert.Nil(err)
	defer func(dir string) {
		os.RemoveAll(dir)
	}(tmpdir)

	tmpLog := filepath.Join(tmpdir, "log.txt")

	o := Options{
		LogFile: tmpLog,
		stderr:  &buffer,
	}

	demoLoggerln(o)

	expect = `WARNING: warn # 4
WARNING: warning # 5
ERROR: error # 6
FATAL: fatal # 7
`
	assert.Equal(expect, buffer.String())

	expect = `WARN[<time>]: warn # 4
WARN[<time>]: warning # 5
ERRO[<time>]: error # 6
FATA[<time>]: fatal # 7
`
	data, err := ioutil.ReadFile(tmpLog)
	assert.Nil(err)
	assert.Equal(expect, filterTime(string(data)))
}

func TestLoggerlnCustom(t *testing.T) {
	var (
		assert = assert.New(t)
		buffer bytes.Buffer
		err    error
		expect string
	)

	tmpdir, err := ioutil.TempDir("", "multi-logger-")
	assert.Nil(err)
	defer func(dir string) {
		os.RemoveAll(dir)
	}(tmpdir)

	tmpLog := filepath.Join(tmpdir, "log.txt")

	o := Options{
		Verbose:  3,
		LogLevel: "error",
		LogFile:  tmpLog,
		stderr:   &buffer,
	}

	demoLoggerln(o)

	expect = `TRACE: trace # 1
DEBUG: debug # 2
INFO: info # 3
WARNING: warn # 4
WARNING: warning # 5
ERROR: error # 6
FATAL: fatal # 7
`
	assert.Equal(expect, buffer.String())

	expect = `ERRO[<time>]: error # 6
FATA[<time>]: fatal # 7
`
	data, err := ioutil.ReadFile(tmpLog)
	assert.Nil(err)
	assert.Equal(expect, filterTime(string(data)))
}

func TestFatalPanic(t *testing.T) {
	var (
		assert = assert.New(t)
		err    error
	)

	tmpdir, err := ioutil.TempDir("", "multi-logger-")
	assert.Nil(err)
	defer func(dir string) {
		os.RemoveAll(dir)
	}(tmpdir)

	o := Options{
		Verbose:  4,
		LogLevel: "warning",
		LogFile:  "",
	}

	Init(o)
	env := os.Getenv("TEST_LOGGER_CRASH")
	if env != "" {
		switch env {
		case "fatalf":
			Fatalf("called %s", env)
		case "fatal":
			Fatal("called ", env)
		case "fatalln":
			Fatalln("called", env)
		case "panicf":
			Panicf("called %s", env)
		case "panic":
			Panic("called ", env)
		case "panicln":
			Panicln("called", env)
		}
		return
	}

	var (
		lock = sync.Mutex{}
		wg   sync.WaitGroup
		msg  = []string{}
	)

	for _, v := range []string{"fatalf", "fatal", "fatalln", "panicf", "panic", "panicln"} {
		wg.Add(1)
		go func(v string) {
			defer wg.Done()
			args := []string{"go", "test", "-test.run=TestFatalPanic"}
			cmd := exec.Command(args[0], args[1:]...)
			env := fmt.Sprintf("TEST_LOGGER_CRASH=%s", v)
			cmd.Env = append(os.Environ(), env)
			out, err := cmd.CombinedOutput()
			if e, ok := err.(*exec.ExitError); ok && !e.Success() {
				if !strings.Contains(string(out), "called "+v) {
					lock.Lock()
					defer lock.Unlock()
					msg = append(msg, fmt.Sprintf("run '%s' with env '%s', wrong output: %s",
						strings.Join(args, " "),
						env,
						string(out),
					))

				}
			} else {
				lock.Lock()
				defer lock.Unlock()
				msg = append(msg, fmt.Sprintf("run '%s' with env '%s' should exit with status 1, but got %v",
					strings.Join(args, " "),
					env,
					err,
				))
			}
		}(v)
	}

	wg.Wait()

	if len(msg) > 0 {
		t.Fatal(strings.Join(msg, "\n"))
	}
}
