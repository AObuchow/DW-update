package update

import (
	dwv1alpha2 "github.com/devfile/api/v2/pkg/apis/workspaces/v1alpha2"
)

func UpdateDevWorkspace(dw dwv1alpha2.DevWorkspace, devfile *dwv1alpha2.Devfile) (updatedDevWorkspace *dwv1alpha2.DevWorkspace) {
	updatedDevWorkspace = dw.DeepCopy()
	// Preserve original devworkspace spec.template.projects
	originalProjects := updatedDevWorkspace.Spec.Template.Projects

	// Find component with the controller.devfile.io/merge-contribution: true attribute
	mergeContributionComponent := ""
	for _, component := range updatedDevWorkspace.Spec.Template.Components {
		if component.Attributes != nil {
			if component.Attributes.Exists("controller.devfile.io/merge-contribution") {
				if component.Attributes.GetBoolean("controller.devfile.io/merge-contribution", nil) {
					mergeContributionComponent = component.Name
					break // There is only supposed to be one merge contribution component so we stop once we find it
				}
			}
		}
	}

	// Replace devworkspace spec.template with devfile content
	updatedDevWorkspace.Spec.Template = devfile.DevWorkspaceTemplateSpec

	// Retain original devworkspace projects
	// TODO: Append here so that the user can add more projects when updating devworkspace?
	updatedDevWorkspace.Spec.Template.Projects = originalProjects

	// Retain merge contribution attribute
	for _, component := range updatedDevWorkspace.Spec.Template.Components {
		if component.Name == mergeContributionComponent {
			if !component.Attributes.Exists("controller.devfile.io/merge-contribution") {
				component.Attributes.PutBoolean("controller.devfile.io/merge-contribution", true)
				break // There is only supposed to be one merge contribution component so we stop once we find it
			}
		}
	}
	return updatedDevWorkspace
}
