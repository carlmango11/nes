package log

import "log"

var Debug = false

func Printf(format string, v ...interface{}) {
	log.Printf(format, v...)
}

func Panicf(format string, v ...interface{}) {
	log.Panicf(format, v...)
}

func Debugf(format string, v ...interface{}) {
	if Debug {
		log.Printf(format, v...)
	}
}
