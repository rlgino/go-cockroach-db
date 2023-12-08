package logging

import (
	"fmt"
	"log"
	"os"
	"strings"

	"go-users-service/internal/core/logger"
)

func (d *defaultLogging) log(msg string, level LevelLog, fields map[string]interface{}) {
	sb := strings.Builder{}
	for key, value := range fields {
		sb.WriteString(fmt.Sprintf("%s=%v", key, value))
	}
	fmt.Printf("level=%s msg=\"%s\" %s \n", level, msg, sb.String())
	d.logger.Printf("level=%s msg=\"%s\" %s \n", level, msg, sb.String())
}

type LevelLog string

const (
	error LevelLog = "error"
	info  LevelLog = "info"
	debug LevelLog = "debug"
	warn  LevelLog = "warn"
	fatal LevelLog = "fatal"
)

func NewDefaultLogging(appName string) logger.Logger {
	mylog := log.New(os.Stdout, appName+": ", log.LstdFlags)
	mylog.SetFlags(log.LstdFlags)
	return &defaultLogging{
		logger: mylog,
	}
}

type defaultLogging struct {
	logger *log.Logger
}

func (d *defaultLogging) Debug(msg string, fields map[string]interface{}) {
	d.log(msg, debug, fields)
}

func (d *defaultLogging) Info(msg string, fields map[string]interface{}) {
	d.log(msg, info, fields)
}

func (d *defaultLogging) Warn(msg string, fields map[string]interface{}) {
	d.log(msg, warn, fields)
}

func (d *defaultLogging) Error(msg string, fields map[string]interface{}) {
	d.log(msg, error, fields)
}

func (d *defaultLogging) Fatal(msg string, fields map[string]interface{}) {
	d.log(msg, fatal, fields)
}
