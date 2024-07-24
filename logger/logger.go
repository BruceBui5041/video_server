package logger

import (
	"fmt"

	"go.uber.org/zap"
)

var AppLogger *zap.Logger

func init() {
	var err error
	AppLogger, err = zap.NewProduction()
	if err != nil {
		panic(fmt.Sprintf("Failed to initialize logger: %v", err))
	}
}
