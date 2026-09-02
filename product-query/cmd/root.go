package cmd

import (
	"github.com/celio001/product-cqrs/product-query/pkg/logger"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

var rootCmd = &cobra.Command{
	Use:  "product-query",
	Args: cobra.MaximumNArgs(1),
}

func Execute() {
	defer func() {
		if err := recover(); err != nil {
			logger.Error("unexpected error while executing command",
				zap.String("error.type", "PanicError"),
				zap.String("error.message", err.(error).Error()),
				zap.String("error.code", "COMMAND_PANIC"),
			)
		}
	}()

	if err := rootCmd.Execute(); err != nil {
		logger.Error("error executing command",
			zap.String("error.type", "CommandError"),
			zap.String("error.message", err.Error()),
			zap.String("error.code", "COMMAND_FAILED"),
		)
	}
}
