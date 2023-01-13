package log

import (
    "fmt"
    "os"
    "time"
)

var (
    Verbose = false
    uptime  = time.Now()
)

func D(format string, a ...any) {
    if !Verbose {
        return
    }
    prefix := fmtDuration(time.Since(uptime))
    _, _ = fmt.Fprintf(os.Stdout, prefix+format, a...)
}
func I(format string, a ...any) {
    prefix := fmtDuration(time.Since(uptime))
    _, _ = fmt.Fprintf(os.Stdout, prefix+format, a...)
}

func E(format string, a ...any) {
    _, _ = fmt.Fprintf(os.Stderr, format, a...)
}
func fmtDuration(d time.Duration) string {
    hour := int(d.Seconds() / 3600)
    minute := int(d.Seconds()/60) % 60
    second := int(d.Seconds()) % 60

    return fmt.Sprintf("[%.2d:%.2d:%.2d]", hour, minute, second)
}
