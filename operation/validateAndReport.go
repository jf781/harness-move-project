package operation

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	"harness-copy-project/model"
	"harness-copy-project/services"

	"go.uber.org/zap"
)

// Function to identify the log level to report on
func logReportingLevel(logLevel string) []string {
	var reportingLevel []string

	switch logLevel {
	case "error":
		reportingLevel = []string{"error"}
	case "warn":
		reportingLevel = []string{"error", "warn"}
	case "info":
		reportingLevel = []string{"error", "warn", "info"}
	}

	return reportingLevel
}

// Function to report on logs from copying an individual project
func ParseAndPrintProjectLogs(logs, logLevel, sourceProject string) {
	logEntries := strings.Split(logs, "\n")

	entryCount := len(logEntries)
	entryIndex := int(1)

	reportingLevel := logReportingLevel(logLevel)

	fmt.Printf("Logs from migrating project: '%v'. Log Level: '%v'. \n ---- \n", sourceProject, logLevel)
	for _, logEntry := range logEntries {
		entryIndex++
		if logEntry == "" {
			continue
		}
		var entry model.LogEntry
		err := json.Unmarshal([]byte(logEntry), &entry)
		if err != nil {
			fmt.Println("Failed to parse log entry:", err)
			continue
		}

		for _, level := range reportingLevel {
			if level == entry.Level {

				printLogs(entry)

				if entryIndex == entryCount {
					fmt.Println(" --")
					fmt.Printf("--- End of log entries for project: '%v' --- \n", sourceProject)
					fmt.Println(" --")
				} else {
					fmt.Println(" ---- ")
				}
			}
		}
	}
}

// Function to report on logs from global operations
func ParseAndPrintGlobalLogs(logs, logLevel string) {
	logEntries := strings.Split(logs, "\n")

	entryCount := len(logEntries)
	entryIndex := int(1)

	reportingLevel := logReportingLevel(logLevel)

	fmt.Printf("Logs encountered when during global operations. \n ---- \n")
	for _, logEntry := range logEntries {
		entryIndex++
		if logEntry == "" {
			continue
		}
		var entry model.LogEntry
		err := json.Unmarshal([]byte(logEntry), &entry)
		if err != nil {
			fmt.Println("Failed to parse log entry:", err)
			continue
		}

		for _, level := range reportingLevel {
			if level == entry.Level {

				printLogs(entry)

				if entryIndex == entryCount {
					fmt.Println(" --")
					fmt.Printf("--- End of log entries for global opreration: --- \n")
					fmt.Println(" --")
				} else {
					fmt.Println(" ---- ")
				}
			}
		}
	}
}

// Function to print log entries
func printLogs(entry model.LogEntry) {
	v := reflect.ValueOf(entry)
	t := v.Type()

	level := entry.Level

	for i := 0; i < v.NumField(); i++ {
		fieldValue := v.Field(i)
		fieldType := t.Field(i)

		// Skip empty fields
		if fieldValue.Kind() == reflect.String && fieldValue.String() == "" {
			continue
		}

		// Print field name and value
		switch level {
		case "error":
			fmt.Printf(Red+"%s: %v\n"+Reset, fieldType.Name, fieldValue.Interface())
		case "warn":
			fmt.Printf(Yellow+"%s: %v\n"+Reset, fieldType.Name, fieldValue.Interface())
		case "info":
			fmt.Printf("%s: %v\n", fieldType.Name, fieldValue.Interface())
		}
	}
}

// Function to report on summary of all projects copied
func OperationSummary(summaryReport []model.ProjectSummary) {
	// Output summary for all projects
	maxSourceLen := len("Source Project")
	maxTargetLen := len("Target Project")
	summaryColor := ""
	for _, summary := range summaryReport {
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

	for _, summary := range summaryReport {
		successStr := "No"
		summaryColor = Red
		if summary.Successful {
			successStr = "Yes"
			summaryColor = Green
		}
		fmt.Printf(summaryColor+rowFmt, summary.SourceProject, summary.TargetProject, successStr+Reset)
	}
}

// Function to create a summary report for each project
func ProjectCopySummary(sourceProject, targetProject string, copyStatus bool) model.ProjectSummary {
	var projectSummary model.ProjectSummary
	projectSummary.SourceProject = sourceProject
	projectSummary.TargetProject = targetProject
	projectSummary.Successful = copyStatus
	return projectSummary
}

// Function to validate and log the copy operation
func ValidateAndLogCopy(cp Copy, logger *zap.Logger) bool {
	var projectErr []bool

	// Output project copy complete message
	fmt.Printf("Project '%v' has been copied to '%v' \n", cp.Source.Project, cp.Target.Project)

	// Output project entity counts
	projectErr = append(projectErr, ConfirmSuccessfulCopy("Connectors", services.GetConnectorsTotal(), services.GetConnectorsMoved()))
	projectErr = append(projectErr, ConfirmSuccessfulCopy("Environments", services.GetEnvironmentsTotal(), services.GetEnvironmentsMoved()))
	projectErr = append(projectErr, ConfirmSuccessfulCopy("Environment Groups", services.GetEnvironmentGroupsTotal(), services.GetEnvironmentGroupsMoved()))
	projectErr = append(projectErr, ConfirmSuccessfulCopy("Feature Flags", services.GetFeatureFlagsTotal(), services.GetFeatureFlagsMoved()))
	projectErr = append(projectErr, ConfirmSuccessfulCopy("File Stores", services.GetFileStoresTotal(), services.GetFileStoresMoved()))
	projectErr = append(projectErr, ConfirmSuccessfulCopy("Infrastructure", services.GetInfrastructureTotal(), services.GetInfrastructureMoved()))
	projectErr = append(projectErr, ConfirmSuccessfulCopy("Input sets", services.GetInputSetsTotal(), services.GetInputSetsMoved()))
	projectErr = append(projectErr, ConfirmSuccessfulCopy("Pipelines", services.GetPipelinesTotal(), services.GetPipelinesMoved()))
	projectErr = append(projectErr, ConfirmSuccessfulCopy("Resource Groups", services.GetResourceGroupsTotal(), services.GetResourceGroupsMoved()))
	projectErr = append(projectErr, ConfirmSuccessfulCopy("Role Assignments", services.GetRoleAssignmentsTotal(), services.GetRoleAssignmentsMoved()))
	projectErr = append(projectErr, ConfirmSuccessfulCopy("Roles", services.GetRolesTotal(), services.GetRolesMoved()))
	projectErr = append(projectErr, ConfirmSuccessfulCopy("Service Overrides", services.GetOverridesTotal(), services.GetOverridesMoved()))
	projectErr = append(projectErr, ConfirmSuccessfulCopy("Services", services.GetServicesTotal(), services.GetServicesMoved()))
	projectErr = append(projectErr, ConfirmSuccessfulCopy("Tags", services.GetTagsTotal(), services.GetTagsMoved()))
	projectErr = append(projectErr, ConfirmSuccessfulCopy("Target Groups", services.GetTargetGroupsTotal(), services.GetTargetGroupsMoved()))
	projectErr = append(projectErr, ConfirmSuccessfulCopy("Targets", services.GetTargetsTotal(), services.GetTargetsMoved()))
	projectErr = append(projectErr, ConfirmSuccessfulCopy("Templates", services.GetTemplatesTotal(), services.GetTemplatesMoved()))
	projectErr = append(projectErr, ConfirmSuccessfulCopy("User Groups", services.GetUserGroupsTotal(), services.GetUserGroupsMoved()))
	projectErr = append(projectErr, ConfirmSuccessfulCopy("Users", services.GetUsersTotal(), services.GetUsersMoved()))
	projectErr = append(projectErr, ConfirmSuccessfulCopy("Variables", services.GetVariablesTotal(), services.GetVariablesMoved()))
	projectErr = append(projectErr, ConfirmSuccessfulCopy("Triggers", services.GetTriggersTotal(), services.GetTriggersMoved()))

	if services.ValidateCopy(projectErr) {
		if err := cp.Freeze(); err != nil {
			logger.Error("Failed to Freeze Project",
				zap.String("Source Project", cp.Source.Project),
				zap.Error(err),
			)
			fmt.Printf(Red+"Error encountered while freezing project: '%v'.  Err: %v \n"+Reset, cp.Source.Project, err)
			return false
		}
	} else {
		fmt.Printf(Red+"Error encountered while copying project: '%v'. \n"+Reset, cp.Target.Project)
		fmt.Printf(Red+"Source project: %v has not be froozen. \n"+Reset, cp.Source.Project)
		return false
	}

	// Output project entity counts to logger
	logger.Info("Project Migration Status:",
		zap.Int("ConnectorsTotal", services.GetConnectorsTotal()),
		zap.Int("ConnectorsMoved", services.GetConnectorsMoved()),
		zap.Int("EnvironmentsTotal", services.GetEnvironmentsTotal()),
		zap.Int("EnvironmentsMoved", services.GetEnvironmentsMoved()),
		zap.Int("EnvironmentGroupsTotal", services.GetEnvironmentGroupsTotal()),
		zap.Int("EnvironmentGroupsMoved", services.GetEnvironmentGroupsMoved()),
		zap.Int("FeatureFlagsTotal", services.GetFeatureFlagsTotal()),
		zap.Int("FeatureFlagsMoved", services.GetFeatureFlagsMoved()),
		zap.Int("FileStoresTotal", services.GetFileStoresTotal()),
		zap.Int("FileStoresMoved", services.GetFileStoresMoved()),
		zap.Int("InfrastructureTotal", services.GetInfrastructureTotal()),
		zap.Int("InfrastructureMoved", services.GetInfrastructureMoved()),
		zap.Int("InputSetsTotal", services.GetInputSetsTotal()),
		zap.Int("InputSetsMoved", services.GetInputSetsMoved()),
		zap.Int("PipelinesTotal", services.GetPipelinesTotal()),
		zap.Int("PipelinesMoved", services.GetPipelinesMoved()),
		zap.Int("ResourceGroupsTotal", services.GetResourceGroupsTotal()),
		zap.Int("ResourceGroupsMoved", services.GetResourceGroupsMoved()),
		zap.Int("RoleAssignmentsTotal", services.GetRoleAssignmentsTotal()),
		zap.Int("RoleAssignmentsMoved", services.GetRoleAssignmentsMoved()),
		zap.Int("RolesTotal", services.GetRolesTotal()),
		zap.Int("RolesMoved", services.GetRolesMoved()),
		zap.Int("OverridesTotal", services.GetOverridesTotal()),
		zap.Int("OverridesMoved", services.GetOverridesMoved()),
		zap.Int("ServicesTotal", services.GetServicesTotal()),
		zap.Int("ServicesMoved", services.GetServicesMoved()),
		zap.Int("TagsTotal", services.GetTagsTotal()),
		zap.Int("TagsMoved", services.GetTagsMoved()),
		zap.Int("TargetGroupsTotal", services.GetTargetGroupsTotal()),
		zap.Int("TargetGroupsMoved", services.GetTargetGroupsMoved()),
		zap.Int("TargetsTotal", services.GetTargetsTotal()),
		zap.Int("TargetsMoved", services.GetTargetsMoved()),
		zap.Int("TemplatesTotal", services.GetTemplatesTotal()),
		zap.Int("TemplatesMoved", services.GetTemplatesMoved()),
		zap.Int("TriggersTotal", services.GetTriggersTotal()),
		zap.Int("TriggersMoved", services.GetTriggersMoved()),
		zap.Int("UserGroupsTotal", services.GetUserGroupsTotal()),
		zap.Int("UserGroupsMoved", services.GetUserGroupsMoved()),
		zap.Int("UsersTotal", services.GetUsersTotal()),
		zap.Int("UsersMoved", services.GetUsersMoved()),
		zap.Int("VariablesTotal", services.GetVariablesTotal()),
		zap.Int("VariablesMoved", services.GetVariablesMoved()),
	)

	return true
}

func ConfirmSuccessfulCopy(entityType string, total, copied int) bool {
	var entityColor string
	var success bool

	if total == copied {
		// All items copied successfully
		entityColor = Green
		success = true
	} else {
		// One or more items failed to copy
		entityColor = Red
		success = false
	}

	fmt.Printf(entityColor+"%v Total: %v \n"+Reset, entityType, total)
	fmt.Printf(entityColor+"%v Moved: %v \n"+Reset, entityType, copied)

	return success
}
