package core

import (
	"fmt"
	"log/slog"
	"reflect"
	"time"

	"github.com/spf13/cobra"
)

// RunModuleWithLogging executes a module with structured logging
func RunModuleWithLogging(entry ModuleEntry, cmd *cobra.Command) error {
	logger := slog.Default()
	startTime := time.Now()
	
	logger.Info("Starting module", slog.String("module", entry.Metadata.UniqueName))

	mod := entry.Factory()
	modVal := reflect.ValueOf(mod).Elem()

	inputField := modVal.FieldByName("Input")
	if !inputField.IsValid() {
		elapsed := time.Since(startTime)
		logger.Error("Module has no Input field",
			slog.String("module", entry.Metadata.UniqueName),
			slog.Duration("duration", elapsed),
		)
		return fmt.Errorf("module '%s' has no Input field", entry.Metadata.UniqueName)
	}

	inputType := inputField.Type()
	inputInstance := reflect.New(inputType).Elem()

	// Populate struct from flags
	setStructFromFlags(inputInstance, inputType, cmd, "")
	inputField.Set(inputInstance)

	// Run module
	err, msg := mod.Run()
	elapsed := time.Since(startTime)
	
	if err != nil {
		logger.Error("Module execution failed",
			slog.String("module", entry.Metadata.UniqueName),
			slog.String("error", err.Error()),
			slog.Duration("duration", elapsed),
		)
		return err
	}
	
	logger.Info("Module executed successfully",
		slog.String("module", entry.Metadata.UniqueName),
		slog.Duration("duration", elapsed),
		slog.String("result", msg),
	)
	
	return nil
}
