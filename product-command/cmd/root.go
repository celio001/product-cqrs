package cmd

import (
	"github.com/celio001/product-command/pkg/logger"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

var rootCmd = &cobra.Command{
	Use:  "product-command",
	Args: cobra.MaximumNArgs(1),
}

func Execute() {
	defer func() {
		err := recover()
		if err != nil {
			logger.Error("unexpected error while executing command",
				zap.String("error.type", "PanicError"),
				zap.String("error.message", err.(error).Error()),
				zap.String("error.code", "COMMAND_PANIC"),
			)
		}
	}()
	err := rootCmd.Execute()
	if err != nil {
		logger.Error("error executing command",
			zap.String("error.type", "CommandError"),
			zap.String("error.message", err.Error()),
			zap.String("error.code", "COMMAND_FAILED"),
		)
	}
}
