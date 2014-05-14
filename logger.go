package main

import (
        "os"
)

type Logger struct {
        log_files []string
}

func assertNoError(err error) {
        if err != nil {
                panic(err)
        }
}

func verifyPath(path string) error {
        _, err := os.Stat(path)
        if err != nil {
                if os.IsNotExist(err) {
                        file, err := os.Create(path)
                        if err != nil {
                                return err
                        }
                        defer file.Close()
                } else {
                        return err
                }
        }
        return nil
}

func newLogger() *Logger {
        logger := &Logger{log_files: make([]string, 0)}
        return logger
}

func (l Logger) exist(fn func(string) bool) bool {
        for _, log_file := range l.log_files {
                if fn(log_file) {
                        return true
                }
        }
        return false
}

func (l *Logger) AddLogFile(path string, empty *struct{}) error {
        if err := verifyPath(path); err != nil {
                return err
        }

        fn := func(log_file string) bool {
                if path == log_file {
                        return true
                }
                return false
        }
        if !(*l).exist(fn) {
                l.log_files = append(l.log_files, path)
        }
        return nil
}

func (l *Logger) AddLogsFile(paths []string, empty *struct{}) error {
        var err error
        for _, path := range paths {
                if err = l.AddLogFile(path, empty); err != nil {
                        return err
                }
        }
        return nil
}
