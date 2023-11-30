package logging

import (
	"go-users-service/internal/core/logger"
	"log"
)

func NewDefaultLogging() logger.Logger {
	return &defaultLogging{}
}

type defaultLogging struct {
	logger log.Logger
}

func (d defaultLogging) Debug(msg string, fields map[string]interface{}) {
	//TODO implement me
	panic("implement me")
}

func (d defaultLogging) Info(msg string, fields map[string]interface{}) {
	//TODO implement me
	panic("implement me")
}

func (d defaultLogging) Warn(msg string, fields map[string]interface{}) {
	//TODO implement me
	panic("implement me")
}

func (d defaultLogging) Error(msg string, fields map[string]interface{}) {
	//TODO implement me
	panic("implement me")
}

func (d defaultLogging) Fatal(msg string, fields map[string]interface{}) {
	//TODO implement me
	panic("implement me")
}
