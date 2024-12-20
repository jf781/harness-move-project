package operation

import (
	// "fmt"
	"harness-copy-project/services"

	"github.com/go-resty/resty/v2"
	"go.uber.org/zap"
)

const (
	Reset  = "\033[0m"
	Red    = "\033[31m"
	Green  = "\033[32m"
	Yellow = "\033[33m"
)

type (
	Config struct {
		Token    string
		Account  string
		BaseURL  string
		Logger   *zap.Logger
		CopyCD   bool
		CopyFF   bool
		ShowPB   bool
		LogLevel string
	}

	// NOT SURE WHICH NAME TO CHOSE TO THAT TYPE
	NoName struct {
		Org     string
		Project string
	}

	Copy struct {
		Config Config
		Source NoName
		Target NoName
	}
)

func (v *Copy) ValidateProjects() bool {

	api := services.ApiRequest{
		Client:  resty.New(),
		Token:   v.Config.Token,
		Account: v.Config.Account,
		BaseURL: v.Config.BaseURL,
	}

	var validateSourcePrj, validateTargetPrj bool

	// Validate Source project
	if api.ValidateProject(v.Source.Org, v.Source.Project, v.Config.Logger) {
		v.Config.Logger.Info("Source project exists",
			zap.String("Project", v.Source.Project),
			zap.String("Org", v.Source.Org),
		)
		validateSourcePrj = true
	} else {
		v.Config.Logger.Error("Failed to confirm source project exists",
			zap.String("Project", v.Source.Project),
			zap.String("Org", v.Source.Org),
		)
		validateSourcePrj = false
	}

	// Validate Target project
	if api.ValidateProject(v.Target.Org, v.Target.Project, v.Config.Logger) {
		v.Config.Logger.Error("Target project already exists",
			zap.String("Project", v.Target.Project),
			zap.String("Org", v.Target.Org),
		)
		validateTargetPrj = false
	} else {
		v.Config.Logger.Info("Target project does not exist",
			zap.String("Project", v.Target.Project),
			zap.String("Org", v.Target.Org),
		)
		validateTargetPrj = true
	}

	if validateSourcePrj && validateTargetPrj {
		return true
	} else {
		return false
	}

}

func (o *Copy) Exec() error {

	api := services.ApiRequest{
		Client:  resty.New(),
		Token:   o.Config.Token,
		Account: o.Config.Account,
		BaseURL: o.Config.BaseURL,
	}

	var operations []services.Operation

	if o.Config.CopyCD || o.Config.CopyFF {
		// CREATE NEW PROJECT IF IT DOES NOT EXIST IN THE TARGET ORG
		operations = append(operations, services.NewProjectOperation(&api, o.Source.Org, o.Source.Project, o.Target.Org, o.Target.Project, o.Config.Logger))
		operations = append(operations, services.RemoveCurrentUserOperation(&api, o.Target.Org, o.Target.Project, o.Config.Logger))
		operations = append(operations, services.NewConnectorOperation(&api, o.Source.Org, o.Source.Project, o.Target.Org, o.Target.Project, o.Config.Logger, o.Config.ShowPB))
		operations = append(operations, services.NewEnvironmentOperation(&api, o.Source.Org, o.Source.Project, o.Target.Org, o.Target.Project, o.Config.Logger, o.Config.ShowPB))
		operations = append(operations, services.NewEnvGroupOperation(&api, o.Source.Org, o.Source.Project, o.Target.Org, o.Target.Project, o.Config.Logger, o.Config.ShowPB))
	}

	if o.Config.CopyCD {
		operations = append(operations, services.NewVariableOperation(&api, o.Source.Org, o.Source.Project, o.Target.Org, o.Target.Project, o.Config.Logger, o.Config.ShowPB))
		operations = append(operations, services.NewFileStoreOperation(&api, o.Source.Org, o.Source.Project, o.Target.Org, o.Target.Project, o.Config.Logger, o.Config.ShowPB))
		operations = append(operations, services.NewInfrastructureOperation(&api, o.Source.Org, o.Source.Project, o.Target.Org, o.Target.Project, o.Config.Logger, o.Config.ShowPB))
		operations = append(operations, services.NewServiceOperation(&api, o.Source.Org, o.Source.Project, o.Target.Org, o.Target.Project, o.Config.Logger, o.Config.ShowPB))
		operations = append(operations, services.NewServiceOverrideOperation(&api, o.Source.Org, o.Source.Project, o.Target.Org, o.Target.Project, o.Config.Logger, o.Config.ShowPB))
		operations = append(operations, services.NewTemplateOperation(&api, o.Source.Org, o.Source.Project, o.Target.Org, o.Target.Project, o.Config.Logger, o.Config.ShowPB))
		operations = append(operations, services.NewPipelineOperation(&api, o.Source.Org, o.Source.Project, o.Target.Org, o.Target.Project, o.Config.Logger, o.Config.ShowPB))
		operations = append(operations, services.NewInputsetOperation(&api, o.Source.Org, o.Source.Project, o.Target.Org, o.Target.Project, o.Config.Logger, o.Config.ShowPB))
		// // operations = append(operations, services.NewSecretOperation(&api, o.Source.Org, o.Source.Project, o.Target.Org, o.Target.Project))
		operations = append(operations, services.NewTagOperation(&api, o.Source.Org, o.Source.Project, o.Target.Org, o.Target.Project, o.Config.Logger, o.Config.ShowPB))
		operations = append(operations, services.NewUserScopeOperation(&api, o.Source.Org, o.Source.Project, o.Target.Org, o.Target.Project, o.Config.Logger, o.Config.ShowPB))
		operations = append(operations, services.NewUserGroupOperation(&api, o.Source.Org, o.Source.Project, o.Target.Org, o.Target.Project, o.Config.Logger, o.Config.ShowPB))
		operations = append(operations, services.NewServiceAccountOperation(&api, o.Source.Org, o.Source.Project, o.Target.Org, o.Target.Project, o.Config.Logger, o.Config.ShowPB))
		operations = append(operations, services.NewRoleOperation(&api, o.Source.Org, o.Source.Project, o.Target.Org, o.Target.Project, o.Config.Logger, o.Config.ShowPB))
		operations = append(operations, services.NewResourceGroupOperation(&api, o.Source.Org, o.Source.Project, o.Target.Org, o.Target.Project, o.Config.Logger, o.Config.ShowPB))
		operations = append(operations, services.NewRoleAssignmentOperation(&api, o.Source.Org, o.Source.Project, o.Target.Org, o.Target.Project, o.Config.Logger, o.Config.ShowPB))
		operations = append(operations, services.NewTriggerOperation(&api, o.Source.Org, o.Source.Project, o.Target.Org, o.Target.Project, o.Config.Logger, o.Config.ShowPB))
	}

	if o.Config.CopyFF {
		operations = append(operations, services.NewFeatureFlagOperation(&api, o.Source.Org, o.Source.Project, o.Target.Org, o.Target.Project, o.Config.Logger, o.Config.ShowPB))
		operations = append(operations, services.NewTargets(&api, o.Source.Org, o.Source.Project, o.Target.Org, o.Target.Project, o.Config.Logger, o.Config.ShowPB))
		operations = append(operations, services.NewTargetGroups(&api, o.Source.Org, o.Source.Project, o.Target.Org, o.Target.Project, o.Config.Logger, o.Config.ShowPB))
	}

	for _, op := range operations {
		if err := op.Copy(); err != nil {
			return err
		}
	}

	return nil
}

func (f *Copy) Freeze() error {

	api := services.ApiRequest{
		Client:  resty.New(),
		Token:   f.Config.Token,
		Account: f.Config.Account,
		BaseURL: f.Config.BaseURL,
	}

	freezeOperation := services.FreezeSourceProjectOperation(&api, f.Source.Org, f.Source.Project, f.Config.Logger)
	if err := freezeOperation.Copy(); err != nil {
		return err
	}

	return nil
}
