package main

import (
        "bufio"
        "errors"
        "os"
)

var (
        NoSuchFile error = errors.New("Error: No such file in Logger")
)

type File struct {
        path    string
        msg     chan []byte
}

type Logger struct {
        log_files []*File
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
        logger := &Logger{log_files: make([]*File, 0, 20)}
        return logger
}

func (l Logger) exist(path string) (*File, bool) {
        for _, log_file := range l.log_files {
                if path == log_file.path {
                        return log_file, true
                }
        }
        return nil, false
}

func (l *Logger) AddLogFile(path string) error {
        if err := verifyPath(path); err != nil {
                return err
        }
        if _, ok := (*l).exist(path); !ok {
                log_file := &File{path: path, msg: make(chan []byte, 10)}
                l.log_files = append(l.log_files, log_file)
                go log_file.listen()
        }
        return nil
}

func (l *Logger) AddLogsFile(paths []string) error {
        var err error
        for _, path := range paths {
                if err = l.AddLogFile(path); err != nil {
                        return err
                }
        }
        return nil
}

func (l *Logger) AsyncWrite(path string, msg []byte) error {
        if log_file, ok := (*l).exist(path); ok {
                log_file.msg <- msg
                return nil
        }
        return NoSuchFile
}

func (f *File) listen() {
        for bytes := range f.msg {
                file, err := os.OpenFile(f.path, os.O_RDWR|os.O_APPEND, 0660)
                assertNoError(err)

                writer := bufio.NewWriter(file)

                _, err = writer.Write(bytes)
                if err != nil {
                        file.Close()
                        assertNoError(err)
                        continue
                }
                writer.Flush()
                file.Close()
        }
}
