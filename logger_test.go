package main

import (
        "bufio"
        "os"
        "strings"
        "testing"
        "time"
)

var (
        lorem = []string{
                `Lorem ipsum dolor sit amet, consectetur adipiscing elit.
Nunc vel purus id arcu tristique iaculis a quis neque.
Nulla ipsum tortor, facilisis tempor cursus et, tincidunt vel nunc.
Praesent adipiscing mi mi, non sodales urna luctus vitae.
Quisque sit amet nunc ut eros lobortis tincidunt.
Fusce ligula eros, tincidunt nec ultrices ut, varius euismod eros.
Proin consectetur dui id eros rutrum, a luctus enim molestie.
Cum sociis natoque penatibus et magnis dis parturient montes, nascetur ridiculus mus.
Donec facilisis dictum metus, non pretium diam imperdiet a.
Interdum et malesuada fames ac ante ipsum primis in faucibus.
In euismod metus quam, eget aliquet est consequat id.
Phasellus molestie tristique nisl ut blandit.
Integer fermentum sed sapien vitae tristique.
Sed elit lorem, rhoncus lobortis orci id, vestibulum laoreet sem.
Pellentesque tristique ipsum tortor, sit amet tincidunt nunc ornare viverra.`,
                // ========================================================
                `Aliquam malesuada suscipit velit id vehicula.
Nam eget euismod turpis.
Vestibulum ante ipsum primis in faucibus orci luctus et ultrices posuere cubilia Curae; Morbi posuere vulputate eleifend.
Donec diam diam, blandit quis sollicitudin vitae, dignissim nec neque.
Maecenas accumsan placerat condimentum.
Nulla eu nunc feugiat, tincidunt justo nec, vehicula risus.`,
                // ========================================================
                `Nam leo justo, posuere eu risus ac, placerat tempus purus.
Integer sollicitudin nisl neque, in blandit nibh dignissim condimentum.
Nulla eleifend placerat nibh.
Suspendisse sed lectus nec felis posuere pulvinar.
Nunc porttitor nibh et diam congue, quis condimentum sem gravida.
Vivamus dapibus, nibh sit amet fringilla convallis, nulla felis ultrices est, at convallis nulla mauris dapibus diam.
Proin laoreet rhoncus sem, eu vulputate dolor rhoncus eu.
Nunc bibendum ipsum pretium, semper libero accumsan, ullamcorper felis.
Phasellus rutrum a enim ut eleifend.
Curabitur faucibus fringilla ullamcorper.
Nullam lacinia est risus, et auctor metus vestibulum nec.
Fusce at mauris erat.
In quis odio eget magna commodo lacinia vel sed felis.
Sed iaculis, mauris vitae consectetur gravida, nisi nibh consectetur quam, eget commodo lectus lorem aliquam enim.`,
                // ========================================================
                `Quisque rutrum mauris a ligula interdum, ac aliquam quam posuere.
Proin scelerisque scelerisque egestas.
Nullam quis lobortis enim.
Curabitur id feugiat tortor.
In rhoncus dictum tellus, vel egestas nulla faucibus ut.
Phasellus molestie convallis vulputate.
Praesent commodo ullamcorper nisl in interdum.
Aliquam a blandit nibh, sit amet aliquam massa.
Nulla quis laoreet enim.
Aenean vel nibh nunc.
Cras vulputate vestibulum metus vel ullamcorper.`}
)

func failTest(t *testing.T, expected, got interface{}) {
        t.Errorf("Error: expected: %v, got: %v", expected, got)
}

func TestAddLogFile(t *testing.T) {
        log := newLogger()
        path := "./logfile.log"
        err := log.AddLogFile(path)
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
        log := newLogger()
        paths := []string{"./logfile0.log",
                "./logfile1.log",
                "./logfile2.log",
                "./logfile3.log"}
        err := log.AddLogsFile(paths)
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
        log := newLogger()
        path := "/logfile.log"
        err := log.AddLogFile(path)
        if err == nil {
                failTest(t, nil, err.Error())
        }
        if len(log.log_files) != 0 {
                failTest(t, 0, len(log.log_files))
        }
}

func TestAddMultipleNoPermissionLogFiles(t *testing.T) {
        log := newLogger()
        paths := []string{"/logfile0.log",
                "/logfile1.log",
                "/logfile2.log",
                "/logfile3.log"}
        err := log.AddLogsFile(paths)
        if err == nil {
                failTest(t, nil, err.Error())
        }
        if len(log.log_files) != 0 {
                failTest(t, 0, len(log.log_files))
        }
}

func assertContent(path, paragraph string, t *testing.T) {
        file, err := os.Open(path)
        if err != nil {
                failTest(t, nil, err.Error())
        }
        defer file.Close()
        reader := bufio.NewReader(file)
        var line []byte
        for _, paragraph_line := range strings.Split(paragraph, "\n") {
                line, _, err = reader.ReadLine()
                if err != nil {
                        failTest(t, nil, err.Error())
                }
                if paragraph_line != string(line) {
                        failTest(t, paragraph_line, line)
                }
        }
        os.Remove(path)
}

func TestWriteOnLogFiles(t *testing.T) {
        log := newLogger()
        paths := []string{"./logfile0.log",
                "./logfile1.log",
                "./logfile2.log",
                "./logfile3.log"}
        err := log.AddLogsFile(paths)
        if err != nil {
                failTest(t, nil, err.Error())
        }
        for i, path := range paths {
                err = log.AsyncWrite(path, []byte(lorem[i]))
                if err != nil {
                        failTest(t, nil, err.Error())
                }
                for len(log.log_files[i].msg) > 0 {
                        time.Sleep(time.Second * 1)
                }
                assertContent(path, lorem[i], t)
        }
}

func BenchmarkWriteByte(b *testing.B) {
        msg := []byte("Hello World!!!\n")
        for i := 0; i < b.N; i++ {
                b.StopTimer()
                file, err := os.OpenFile("./logfile0.log", os.O_RDWR|os.O_APPEND, 0660)
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
                file, err := os.OpenFile("./logfile1.log", os.O_RDWR|os.O_APPEND, 0660)
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
                file, err := os.OpenFile("./logfile2.log", os.O_RDWR|os.O_APPEND, 0660)
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
