package clog

import (
	"fmt"
	"log"
	"runtime/debug"
	"strings"
)

const (
	LogALL int = iota
	LogTRACE
	LogDEBUG
	LogPRINT
	LogINFO
	LogWARN
	LogERROR
	LogFATAL
	LogCritical
	LogOFF
)

// const DebugLevel = LogCritical
// const DebugLevel = LogPRINT
const DebugLevel = LogPRINT

const errorWrapperStart = "--------------- This is an error message ---------------"
const errorwrapperEnd = "--------------------------------------------------------"

// const DebugLevel = LogWARN
func strBefore(input string, split string) string {
	val := strings.LastIndex(input, split)
	if val > 0 {
		return strings.TrimSpace(input[:val])
	} else {
		return input
	}
}

func strAfter(input string, split string) string {
	val := strings.LastIndex(input, split)
	if val > 0 {
		return strings.TrimSpace(input[val:])
	} else {
		return input
	}
}

func callInfo() string {
	stack := debug.Stack()
	stackLines := strings.Split(string(stack), "\n")
	routine0 := strings.Replace(stackLines[0], "[running]:", "", -1)
	routine1 := strings.Replace(routine0, "goroutine", "", -1)
	routine := strings.TrimSpace(routine1)
	funcName := strAfter(strBefore(stackLines[7], "("), "/")
	lineNum := strBefore(strAfter(stackLines[8], ":"), "+")

	format := fmt.Sprintf("routine:%s-%s%s", routine, funcName, lineNum)
	return format
}

func output(level int, format string, args ...interface{}) {
	if level >= DebugLevel {
		log.Printf(format, args...)
	}
}

func Print(format string, args ...interface{}) {
	format = fmt.Sprintf("[%s] ", callInfo()) + format
	output(LogPRINT, format, args...)
}

func Debug(format string, args ...interface{}) {
	format = fmt.Sprintf("[%s] ", callInfo()) + format
	output(LogDEBUG, format, args...)
}

func Info(format string, args ...interface{}) {
	format = fmt.Sprintf("[%s] ", callInfo()) + format
	output(LogINFO, format, args...)
}

func Warn(format string, args ...interface{}) {
	format = fmt.Sprintf("[WARN:%s] ", callInfo()) + format
	output(LogWARN, format, args...)
}

func Error(format string, args ...interface{}) {
	format = "\n" + errorWrapperStart + "\n" + fmt.Sprintf("[ERROR:%s] ", callInfo()) + format + "\n" + errorwrapperEnd + "\n"
	output(LogERROR, format, args...)
}

func Fatal(format string, args ...interface{}) {
	format = "\n" + errorWrapperStart + "\n" + fmt.Sprintf("[FATAL:%s] ", callInfo()) + format + "\n" + errorwrapperEnd + "\n"
	output(LogFATAL, format, args...)
}

func Critical(format string, args ...interface{}) {
	format = fmt.Sprintf("[%s] ", callInfo()) + format
	output(LogCritical, format, args...)
}
