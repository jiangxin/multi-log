Multi-log
=========

[![Travis Build Status](https://travis-ci.org/jiangxin/multi-log.svg?branch=master)](https://travis-ci.org/jiangxin/multi-log)

Multi-log is based on logrus, and provides logging in two directions at the same time.
One for logging on console and one for file logging.

## Usage

### Log on console with builtin logger

    import "github.com/jiangxin/multi-log"

    func main() {
        log.Trace("trace ...")
        log.Debug("debug ...")
        log.Info("info ...")
        log.Warn("warn ...")
        log.Error("error ...")
        log.Fatal("fatal ...")
    }

### To log on console and file, call Init() first

    import "github.com/jiangxin/multi-log"

    func main() {
        log.Init(log.Options{
                Quiet: false,
                Verbose: 2,
                LogFile: "/var/log/my-app.log",
                LogLevel: "warn",
        })

        log.Tracef("trace %s", "...")
        log.Debugf("debug %s", "...")
        log.Infof("info %s", "...")
        log.Warningf("info %s", "...")
        log.Errorf("info %s", "...")
        log.Fatalf("info %s", "...")
    }

### Logging with fields

    import (
        "github.com/jiangxin/multi-log"
        "time"
    )

    func main() {
        logger := log.WithFields(map[string]interface{}{
            "size":   "10MB",
            "period": 2 * time.Minute,
        })

        logger.Traceln("trace", "...")
        logger.Debugln("debug", "...")
        logger.Infoln("info", "...")
        logger.Warnln("warn", "...")
        logger.Errorln("error", "...")
        logger.Panicln("panic", "...")
    }

### Always show notes on console, unless quiet

    import (
        "github.com/jiangxin/multi-log"
    )

    func main() {
        log.Init(log.Options{
                Quiet: false,
        })

        log.Notef("note %s", "...")
        log.Note("note ", "...")
        log.Noteln("note", "...")

        log.Printf("print %s", "...")
        log.Print("print ", "...")
        log.Println("print", "...")
    }





