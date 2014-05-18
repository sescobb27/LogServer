package main

import (
        "bufio"
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

func BenchmarkWriteByte(b *testing.B) {
        msg := []byte("Hello World!!!\n")
        for i := 0; i < b.N; i++ {
                b.StopTimer()
                file, err := os.OpenFile("./logfile0.log", os.O_RDWR|os.O_APPEND, 0777)
                assertNoError(err)

                writer := bufio.NewWriter(file)
                b.StartTimer()
                for _, v := range msg {
                        err := writer.WriteByte(v)
                        if err != nil {
                                file.Close()
                                assertNoError(err)
                                return
                        }
                }
                writer.Flush()
                file.Close()
        }
}

func BenchmarkWriteString(b *testing.B) {
        msg := "Hello World!!!\n"
        for i := 0; i < b.N; i++ {
                b.StopTimer()
                file, err := os.OpenFile("./logfile1.log", os.O_RDWR|os.O_APPEND, 0777)
                assertNoError(err)

                writer := bufio.NewWriter(file)
                b.StartTimer()
                _, err = writer.WriteString(msg)
                if err != nil {
                        file.Close()
                        assertNoError(err)
                        return
                }
                writer.Flush()
                file.Close()
        }
}

func BenchmarkWrite(b *testing.B) {
        msg := []byte("Hello World!!!\n")
        for i := 0; i < b.N; i++ {
                b.StopTimer()
                file, err := os.OpenFile("./logfile2.log", os.O_RDWR|os.O_APPEND, 0777)
                assertNoError(err)

                writer := bufio.NewWriter(file)
                b.StartTimer()
                _, err = writer.Write(msg)
                if err != nil {
                        file.Close()
                        assertNoError(err)
                        return
                }
                writer.Flush()
                file.Close()
        }
}
