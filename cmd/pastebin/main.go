package main

import (
	"context"
	"github.com/artshirshov/gastebin/internal/app"
	"github.com/artshirshov/gastebin/pkg/logger"
	"go.uber.org/zap"
)

func init() {
	logger.SetDefaultZapLogger()
}

func main() {
	ctx := context.Background()

	a, err := app.NewApp(ctx)
	if err != nil {
		logger.Log.With(zap.Error(err)).Fatal("failed to init app")
	}

	err = a.Run()
	if err != nil {
		logger.Log.With(zap.Error(err)).Fatal("failed to run app")
	}
}
