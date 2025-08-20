package log

import (
	"fmt"
	"log"

	"log/slog"
	"testing"
	"time"
)

func TestLog(t *testing.T) {
	Init(LogOptions{
		Handle: func(time time.Time, msg, level string, source *slog.Source, attrs []slog.Attr) {
			fmt.Printf("%s %s %s %+v[%d] %+v\n", time.Format("2006-01-02 15:04:05"), level, msg, source.Function, source.Line, attrs)
		},
		Level: "info",
	})


	//log.Fatalf("hello world111")
	slog.Error("hello world222")
	log.Fatal("hello world333")


	// var log *slog.Logger
	// log = slog.With("func", "TestLog")
	// log.Info("hello world")
	// log.Info("hello world2")

	// log = slog.Default()
	// log = log.With("func2", "TestLog")
	// log.Info("hello world3")
}

func TestLog2(t *testing.T) {
	Init(LogOptions{})
	slog.With("func", "TestLog2").Info("hello world")
	fmt.Println("hello world")
}

func TestLog3(t *testing.T) {
	slog.With("func", "TestLog").Info("hello world")
}
