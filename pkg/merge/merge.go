package merge

import (
	"fmt"
	"os"
	"path/filepath"

	"dslab.sjtu/kubenet-config-generator/pkg/kubenet_config_util"
	"dslab.sjtu/kubenet/api/pkg/kubenet_config"
	"github.com/sirupsen/logrus"
)

func (m *ConfigMerger) Merge() error {
	if m.InputPath == "" || m.OutputPath == "" {
		return fmt.Errorf("input and output paths are required")
	}

	if _, err := os.Stat(m.InputPath); os.IsNotExist(err) {
		return fmt.Errorf("input path does not exist: %v", m.InputPath)
	}

	entries, err := os.ReadDir(m.InputPath)
	if err != nil {
		return fmt.Errorf("failed to read directory: %v", err)
	}

	var networkExperiments []*kubenet_config.NetworkExperiment

	for _, entry := range entries {
		if entry.IsDir() && !m.Recursive {
			continue
		}

		filename := filepath.Join(m.InputPath, entry.Name())
		if m.FilterExt != "" && filepath.Ext(filename) != m.FilterExt {
			continue
		}

		// fmt.Println("merge file:", filename)

		entryNetworkExperiment, err := kubenet_config_util.ReadConfig(filename)
		if err != nil {
			if m.IgnoreError {
				logrus.Warnf("skip invalid config: %v", filename)
				continue
			}
			return fmt.Errorf("failed to read config %v: %w", filename, err)
		}
		networkExperiments = append(networkExperiments, entryNetworkExperiment)
	}

	if len(networkExperiments) == 0 {
		return fmt.Errorf("no valid input files found in %s", m.InputPath)
	}

	networkExperiment := MergeNetworkExperiment(networkExperiments)
	if err := kubenet_config_util.WriteConfig(m.OutputPath, networkExperiment); err != nil {
		return fmt.Errorf("failed to write merged config: %w", err)
	}

	fmt.Printf("✅ Merged %d configs into %s\n", len(networkExperiments), m.OutputPath)
	return nil
}
