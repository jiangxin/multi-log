package log

import (
	"bytes"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func demoWithFieldsLogger(o Options) {
	o.exitFunc = testingExitFunc
	Init(o)
	i := 1

	logger := WithFields(map[string]interface{}{
		"size":   "10MB",
		"period": 2 * time.Minute,
	})
	logger.Errorf("with-field #%s", i)
	logger.Error("with-field #", i)
	logger.Errorln("with-field #", i)
}

func TestLoggerWithFields(t *testing.T) {
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
		LogFile: tmpLog,
		stderr:  &buffer,
	}

	demoWithFieldsLogger(o)

	assert.Nil(err)
	expect = `ERROR: with-field #%!s(int=1)                       (period=2m0s size=10MB)
ERROR: with-field #1                                (period=2m0s size=10MB)
ERROR: with-field # 1                               (period=2m0s size=10MB)` + "\n"
	assert.Equal(expect, buffer.String())

	expect = `ERRO[<time>]: with-field #%!s(int=1)                       (period=2m0s size=10MB)
ERRO[<time>]: with-field #1                                (period=2m0s size=10MB)
ERRO[<time>]: with-field # 1                               (period=2m0s size=10MB)` + "\n"
	data, err := ioutil.ReadFile(tmpLog)
	assert.Nil(err)
	assert.Equal(expect, filterTime(string(data)))
}
