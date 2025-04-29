package kubenet

import "dslab.sjtu/kubenet-config-generator/api/pkg/kubenet_config"

func NewFadingModels() ([]kubenet_config.FadingModel, error) {
	fadindModels := []kubenet_config.FadingModel{}

	fadindModel := kubenet_config.FadingModel{}

	fadindModels = append(fadindModels, fadindModel)
	return fadindModels, nil
}
