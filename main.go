package main

import (
	"bytes"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/urfave/cli/v2"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"harness-copy-project/operation"
	"harness-copy-project/services"
)

var Version = "development"
var errs []error
var globalLogger *zap.Logger
var globalLogBuffer bytes.Buffer
var loopLogger *zap.Logger
var SummaryReport []operation.ProjectSummary

func main() {

	// Initlize and configure the logger
	globalConfig := zap.NewProductionConfig()
	globalConfig.Level = zap.NewAtomicLevelAt(zapcore.InfoLevel) // Set to Info, Debug, or Error for more verbose logging

	// Create a core that writes to the buffer
	globalCore := zapcore.NewCore(
		zapcore.NewJSONEncoder(globalConfig.EncoderConfig),
		zapcore.AddSync(&globalLogBuffer),
		globalConfig.Level,
	)

	globalLogger = zap.New(globalCore)

	startTime := time.Now()

	// Defer the logger and print the final log
	defer globalLogger.Sync()
	defer func() {
		stopTime := time.Now()
		apiCalls := services.GetApiCalls()
		projects := services.GetProjects()
		duration := stopTime.Sub(startTime)

		var avgApiCallDuration time.Duration
		if apiCalls > 0 {
			avgApiCallDuration = duration / time.Duration(apiCalls)
		} else {
			avgApiCallDuration = 0
		}

		globalLogger.Info("Harness Copy Project has completed.",
			zap.Int("Number of API Calls: ", apiCalls),
			zap.String("Run Duration: ", duration.String()),
			zap.Duration("Average API Call Duration: ", avgApiCallDuration),
			zap.Int("Number of projects moved: ", projects),
			zap.String("Stop Time: ", stopTime.Format("12:00:00")),
		)

	}()

	// Start logging
	globalLogger.Info("Harness Copy Project has started.")
	globalLogger.Info("Start time: " + startTime.Format("12:00:00"))

	// Create a new CLI app
	app := &cli.App{
		Name:    "harness-copy-project",
		Version: Version,
		Usage:   "Non-official Harness CLI to copy project between organizations",
		Action:  run,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "csvPath",
				Usage:    "The path to the CSV file that contains the source and target project information.",
				Required: true,
			},
			&cli.StringFlag{
				Name:     "apiToken",
				Usage:    "The API token that will be used to authenticate with the Harness Account.",
				Required: true,
			},
			&cli.StringFlag{
				Name:     "accountId",
				Usage:    "The account ID that contains both the source and target orgnaizations.",
				Required: true,
			},
			&cli.StringFlag{
				Name:     "baseUrl",
				Usage:    "The URL of the harness instance that your projects reside in.",
				Required: true,
			},
			&cli.BoolFlag{
				Name:     "copyCDComponents",
				Usage:    "If set to 'true', then it will copy the Continuous Delivery components.",
				Required: false,
				Value:    false,
			},
			&cli.BoolFlag{
				Name:     "copyFFComponents",
				Usage:    "If set to 'true, then it will copy the Feature Flag components.",
				Required: false,
				Value:    false,
			},
			&cli.BoolFlag{
				Name:     "showProgressBar",
				Usage:    "If set to 'true, then it will show the progress bar for items as they are copied.",
				Required: false,
				Value:    false,
			},
			&cli.StringFlag{
				Name:     "logLevel",
				Usage:    "Defines the level of logs returned.  Valid responses are 'info', 'warn' and 'error'.",
				Required: false,
				Value:    "error",
			},
		},
	}

	// Run the CLI app
	app.Run(os.Args)
}

func run(c *cli.Context) error {
	importCsv := operation.ImportCSV{
		CsvPath: c.String("csvPath"),
	}

	csvData, err := importCsv.Exec()
	if err != nil {
		globalLogger.Error("Failed to pull CSV data",
			zap.String("csvPath", c.String("csvPath")),
			zap.Error(err),
		)
		return err
	}

	logLevel := strings.ToLower(c.String("logLevel"))

	for i := 0; i < len(csvData.SourceOrg); i++ {
		// Create a new log buffer for the project
		var loopLogBuffer bytes.Buffer

		// Initialize and configure the logger for the project
		loopConfig := zap.NewProductionConfig()
		loopConfig.Level = zap.NewAtomicLevelAt(zapcore.InfoLevel) // Set to Info, Debug, or Error for more verbose logging

		loopCore := zapcore.NewCore(
			zapcore.NewJSONEncoder(loopConfig.EncoderConfig),
			zapcore.AddSync(&loopLogBuffer),
			loopConfig.Level,
		)

		loopLogger = zap.New(loopCore)

		// Increment the number of projects moved
		services.IncrementProjects()

		// Create a new copy operation
		cp := operation.Copy{
			Config: operation.Config{
				Token:    c.String("apiToken"),
				Account:  c.String("accountId"),
				BaseURL:  c.String("baseUrl"),
				Logger:   loopLogger,
				CopyCD:   c.Bool("copyCDComponents"),
				CopyFF:   c.Bool("copyFFComponents"),
				ShowPB:   c.Bool("showProgressBar"),
				LogLevel: logLevel,
			},
			Source: operation.NoName{
				Org:     csvData.SourceOrg[i],
				Project: csvData.SourceProject[i],
			},
			Target: operation.NoName{
				Org:     csvData.TargetOrg[i],
				Project: csvData.TargetProject[i],
			},
		}

		// Check for missing or empty values
		if cp.Source.Org == "" || cp.Source.Project == "" || cp.Target.Org == "" {
			loopLogger.Warn("Invalid CSV data. Missing required fields.",
				zap.String("Source Org", cp.Source.Org),
				zap.String("Source Project", cp.Source.Project),
				zap.String("Target Org", cp.Target.Org),
			)
			continue // Skip this iteration if required data is missing
		}

		// Use source project name if target project name is missing
		if cp.Target.Project == "" {
			cp.Target.Project = cp.Source.Project
		}

		fmt.Println(color.GreenString("Moving project '%v' from org '%v' to org '%v'. The target project will be named '%v'", cp.Source.Project, cp.Source.Org, cp.Target.Org, cp.Target.Project))

		// Execute the copy operation from source to target operation
		if err := cp.Exec(); err != nil {
			loopLogger.Error("Failed to Copy Project",
				zap.String("Source Project", cp.Source.Project),
				zap.String("Target Project", cp.Target.Project),
				zap.Error(err),
			)
			errs = append(errs, err)
			continue
		}

		// Validate the copy operation
		if err := operation.ValidateAndLogCopy(cp, loopLogger); err != nil {
			errs = append(errs, err)
			continue
		}

		loopLogger.Info(fmt.Sprintf("Project '%v' has been copied to org: '%v' \n", cp.Source.Project, cp.Target.Org))

		// Reset the API call counter
		services.ResetAllCounters()

		// Parse and filter error messages for the project
		operation.ParseAndPrintProjectLogs(loopLogBuffer.String(), logLevel, cp.Source.Project)

		// Create a summary report for the project
		currentProjectSummary := operation.ProjectCopySummary(cp.Source.Project, cp.Target.Project, loopLogBuffer.String())
		SummaryReport = append(SummaryReport, currentProjectSummary)

	}

	// Parse and filter error messages for the global operation
	operation.ParseAndPrintGlobalLogs(globalLogBuffer.String(), logLevel)

	// Output summary for all projects
	maxSourceLen := len("Source Project")
	maxTargetLen := len("Target Project")
	for _, summary := range SummaryReport {
		if len(summary.SourceProject) > maxSourceLen {
			maxSourceLen = len(summary.SourceProject)
		}
		if len(summary.TargetProject) > maxTargetLen {
			maxTargetLen = len(summary.TargetProject)
		}
	}

	headerFmt := fmt.Sprintf("%%-%ds  %%-%ds  %%s\n", maxSourceLen, maxTargetLen)
	rowFmt := fmt.Sprintf("%%-%ds  %%-%ds  %%s\n", maxSourceLen, maxTargetLen)

	fmt.Println("\nSummary Report:")
	fmt.Printf(headerFmt, "Source Project", "Target Project", "Successful")
	fmt.Println(strings.Repeat("-", maxSourceLen+maxTargetLen+15))

	for _, summary := range SummaryReport {
		successStr := "No"
		if summary.Successful {
			successStr = "Yes"
		}
		fmt.Printf(rowFmt, summary.SourceProject, summary.TargetProject, successStr)
	}

	return nil
}
