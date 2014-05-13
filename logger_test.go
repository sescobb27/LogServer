package main

import (
        "os"
        "testing"
)

func failTest(t *testing.T, expected, got interface{}) {
        t.Errorf("Error: expected: %v, got: %v", expected, got)
}

func TestAddLogFile(t *testing.T) {
        var empty = new(struct{})
        log := newLogger()
        path := "./logfile.log"
        err := log.AddLogFile(path, empty)
        if err != nil {
                failTest(t, nil, err.Error())
        }
        _, err = os.Stat(path)
        if err != nil {
                failTest(t, nil, err.Error())
        }
        if len(log.log_files) != 1 {
                failTest(t, 1, len(log.log_files))
        }
}

func TestAddMultipleLogFiles(t *testing.T) {
        var empty = new(struct{})
        log := newLogger()
        paths := []string{"./logfile0.log",
                "./logfile1.log",
                "./logfile2.log",
                "./logfile3.log"}
        err := log.AddLogsFile(paths, empty)
        if err != nil {
                failTest(t, nil, err.Error())
        }
        for _, path := range paths {
                _, err = os.Stat(path)
                if err != nil {
                        failTest(t, nil, err.Error())
                }
        }
        if len(log.log_files) != 4 {
                failTest(t, 4, len(log.log_files))
        }
}

func TestAddNoPermissionLogFile(t *testing.T) {
        var empty = new(struct{})
        log := newLogger()
        path := "/logfile.log"
        err := log.AddLogFile(path, empty)
        if err == nil {
                failTest(t, nil, err.Error())
        }
        if len(log.log_files) != 0 {
                failTest(t, 0, len(log.log_files))
        }
}

func TestAddMultipleNoPermissionLogFiles(t *testing.T) {
        var empty = new(struct{})
        log := newLogger()
        paths := []string{"/logfile0.log",
                "/logfile1.log",
                "/logfile2.log",
                "/logfile3.log"}
        err := log.AddLogsFile(paths, empty)
        if err == nil {
                failTest(t, nil, err.Error())
        }
        if len(log.log_files) != 0 {
                failTest(t, 0, len(log.log_files))
        }
}
