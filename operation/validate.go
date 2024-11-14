package operation

import (
	"fmt"
	"harness-copy-project/services"
	"go.uber.org/zap"
	"github.com/fatih/color"
)

func ValidateAndLogCopy (cp Copy, logger *zap.Logger) error {
	var projectErr []bool

	fmt.Println(color.GreenString("Project '%v' has been copied to '%v'", cp.Source.Project, cp.Target.Project))
	fmt.Println(color.GreenString("Connectors Total: %v ", services.GetConnectorsTotal()))
	fmt.Println(color.GreenString("Connectors Moved: %v ", services.GetConnectorsMoved()))
	projectErr = append(projectErr, services.ConfirmSuccessfulCopy(services.GetConnectorsTotal(), services.GetConnectorsMoved()))

	fmt.Println(color.GreenString("Environments Total: %v ", services.GetEnvironmentsTotal()))
	fmt.Println(color.GreenString("Environments Moved: %v ", services.GetEnvironmentsMoved()))
	projectErr = append(projectErr, services.ConfirmSuccessfulCopy(services.GetEnvironmentsTotal(), services.GetEnvironmentsMoved()))

	fmt.Println(color.GreenString("EnvironmentGroups Total: %v ", services.GetEnvironmentGroupsTotal()))
	fmt.Println(color.GreenString("EnvironmentGroups Moved: %v ", services.GetEnvironmentGroupsMoved()))
	projectErr = append(projectErr, services.ConfirmSuccessfulCopy(services.GetEnvironmentGroupsTotal(), services.GetEnvironmentGroupsMoved()))

	fmt.Println(color.GreenString("FeatureFlags Total: %v ", services.GetFeatureFlagsTotal()))
	fmt.Println(color.GreenString("FeatureFlags Moved: %v ", services.GetFeatureFlagsMoved()))
	projectErr = append(projectErr, services.ConfirmSuccessfulCopy(services.GetFeatureFlagsTotal(), services.GetFeatureFlagsMoved()))

	fmt.Println(color.GreenString("FileStores Total: %v ", services.GetFileStoresTotal()))
	fmt.Println(color.GreenString("FileStores Moved: %v ", services.GetFileStoresMoved()))
	projectErr = append(projectErr, services.ConfirmSuccessfulCopy(services.GetFileStoresTotal(), services.GetFileStoresMoved()))

	fmt.Println(color.GreenString("Infrastructure Total: %v ", services.GetInfrastructureTotal()))
	fmt.Println(color.GreenString("Infrastructure Moved: %v ", services.GetInfrastructureMoved()))
	projectErr = append(projectErr, services.ConfirmSuccessfulCopy(services.GetInfrastructureTotal(), services.GetInfrastructureMoved()))

	fmt.Println(color.GreenString("InputSets Total: %v ", services.GetInputSetsTotal()))
	fmt.Println(color.GreenString("InputSets Moved: %v ", services.GetInputSetsMoved()))
	projectErr = append(projectErr, services.ConfirmSuccessfulCopy(services.GetInputSetsTotal(), services.GetInputSetsMoved()))

	fmt.Println(color.GreenString("Pipelines Total: %v ", services.GetPipelinesTotal()))
	fmt.Println(color.GreenString("Pipelines Moved: %v ", services.GetPipelinesMoved()))
	projectErr = append(projectErr, services.ConfirmSuccessfulCopy(services.GetPipelinesTotal(), services.GetPipelinesMoved()))

	fmt.Println(color.GreenString("ResourceGroups Total: %v ", services.GetResourceGroupsTotal()))
	fmt.Println(color.GreenString("ResourceGroups Moved: %v ", services.GetResourceGroupsMoved()))
	projectErr = append(projectErr, services.ConfirmSuccessfulCopy(services.GetResourceGroupsTotal(), services.GetResourceGroupsMoved()))

	fmt.Println(color.GreenString("RoleAssignments Total: %v ", services.GetRoleAssignmentsTotal()))
	fmt.Println(color.GreenString("RoleAssignments Moved: %v ", services.GetRoleAssignmentsMoved()))
	projectErr = append(projectErr, services.ConfirmSuccessfulCopy(services.GetRoleAssignmentsTotal(), services.GetRoleAssignmentsMoved()))

	fmt.Println(color.GreenString("Roles Total: %v ", services.GetRolesTotal()))
	fmt.Println(color.GreenString("Roles Moved: %v ", services.GetRolesMoved()))
	projectErr = append(projectErr, services.ConfirmSuccessfulCopy(services.GetRolesTotal(), services.GetRolesMoved()))

	fmt.Println(color.GreenString("Service Overrides Total: %v ", services.GetOverridesTotal()))
	fmt.Println(color.GreenString("Service Overrides Moved: %v ", services.GetOverridesMoved()))
	projectErr = append(projectErr, services.ConfirmSuccessfulCopy(services.GetOverridesTotal(), services.GetOverridesMoved()))

	fmt.Println(color.GreenString("Services Total: %v ", services.GetServicesTotal()))
	fmt.Println(color.GreenString("Services Moved: %v ", services.GetServicesMoved()))
	projectErr = append(projectErr, services.ConfirmSuccessfulCopy(services.GetServicesTotal(), services.GetServicesMoved()))

	fmt.Println(color.GreenString("Tags Total: %v ", services.GetTagsTotal()))
	fmt.Println(color.GreenString("Tags Moved: %v ", services.GetTagsMoved()))
	projectErr = append(projectErr, services.ConfirmSuccessfulCopy(services.GetTagsTotal(), services.GetTagsMoved()))

	fmt.Println(color.GreenString("TargetGroups Total: %v ", services.GetTargetGroupsTotal()))
	fmt.Println(color.GreenString("TargetGroups Moved: %v ", services.GetTargetGroupsMoved()))
	projectErr = append(projectErr, services.ConfirmSuccessfulCopy(services.GetTargetGroupsTotal(), services.GetTargetGroupsMoved()))

	fmt.Println(color.GreenString("Targets Total: %v ", services.GetTargetsTotal()))
	fmt.Println(color.GreenString("Targets Moved: %v ", services.GetTargetsMoved()))
	projectErr = append(projectErr, services.ConfirmSuccessfulCopy(services.GetTargetsTotal(), services.GetTargetsMoved()))

	fmt.Println(color.GreenString("Templates Total: %v ", services.GetTemplatesTotal()))
	fmt.Println(color.GreenString("Templates Moved: %v ", services.GetTemplatesMoved()))
	projectErr = append(projectErr, services.ConfirmSuccessfulCopy(services.GetTemplatesTotal(), services.GetTemplatesMoved()))

	fmt.Println(color.GreenString("UserGroups Total: %v ", services.GetUserGroupsTotal()))
	fmt.Println(color.GreenString("UserGroups Moved: %v ", services.GetUserGroupsMoved()))
	projectErr = append(projectErr, services.ConfirmSuccessfulCopy(services.GetUserGroupsTotal(), services.GetUserGroupsMoved()))

	fmt.Println(color.GreenString("Users Total: %v ", services.GetUsersTotal()))
	fmt.Println(color.GreenString("Users Moved: %v ", services.GetUsersMoved()))
	projectErr = append(projectErr, services.ConfirmSuccessfulCopy(services.GetUsersTotal(), services.GetUsersMoved()))

	fmt.Println(color.GreenString("Variables Total: %v ", services.GetVariablesTotal()))
	fmt.Println(color.GreenString("Variables Moved: %v ", services.GetVariablesMoved()))
	projectErr = append(projectErr, services.ConfirmSuccessfulCopy(services.GetVariablesTotal(), services.GetVariablesMoved()))

	fmt.Println(color.GreenString("Triggers Total: %v ", services.GetTriggersTotal()))
	fmt.Println(color.GreenString("Triggers Moved: %v ", services.GetTriggersMoved()))
	projectErr = append(projectErr, services.ConfirmSuccessfulCopy(services.GetTriggersTotal(), services.GetTriggersMoved()))

	if services.ValidateCopy(projectErr) {
		if err := cp.Freeze(); err != nil {
			logger.Error("Failed to Copy Project",
				zap.String("Source Project", cp.Source.Project),
				zap.String("Target Project", cp.Target.Project),
				zap.Error(err),
			)
			fmt.Println(color.RedString("Error encountered while moving project %v: %v ", cp.Target.Project, err))
			return err
		}
	} else {
		fmt.Println(color.RedString("Error encountered while moving project %v: %v ", cp.Target.Project, projectErr))
	}

	fmt.Println(color.GreenString("---------------------------------\n\n"))

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

	return nil
}