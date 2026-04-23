package merge

import (
	"fmt"
	"os"
	"sync"

	"dslab.sjtu/kubenet-config-generator/pkg/kubenet_config_util"
	"github.com/sirupsen/logrus"
)

func (m *ConfigMerger) Merge() error {
	logrus.Infof("🔀 Start reading yaml files...")

	wg := sync.WaitGroup{}
	wg.Add(len(m.InputPaths))

	for _, inputPath := range m.InputPaths {
		go func(path string) {
			if err := m.ReadNetworkExperiment(path); err != nil {
				logrus.Errorf("error reading config %v: %s", path, err)
			}
			wg.Done()
		}(inputPath)
	}

	wg.Wait()

	logrus.Infof("🔀 Merging %d configs...", len(m.networkExperiments))

	if len(m.networkExperiments) == 0 {
		return fmt.Errorf("no valid config needed to be merged")
	}

	networkExperiment := MergeNetworkExperiment(m.networkExperiments)

	logrus.Infof("✅ Merged %d configs\n", len(m.networkExperiments))

	if err := kubenet_config_util.WriteConfig(m.OutputPath, networkExperiment); err != nil {
		return fmt.Errorf("failed to write merged config: %w", err)
	}

	logrus.Infof("✅ Wrote merged config to %s\n", m.OutputPath)

	return nil
}

func (m *ConfigMerger) ReadNetworkExperiment(inputPath string) error {
	if _, err := os.Stat(inputPath); os.IsNotExist(err) {
		return fmt.Errorf("input path does not exist: %v", inputPath)
	}

	entryNetworkExperiment, err := kubenet_config_util.ReadConfig(inputPath)
	if err != nil {
		if m.IgnoreError {
			logrus.Warnf("skip invalid config: %v", inputPath)
			return nil
		}
		return fmt.Errorf("failed to read config %v: %w", inputPath, err)
	}

	m.mu.Lock()
	m.networkExperiments = append(m.networkExperiments, entryNetworkExperiment)
	m.mu.Unlock()
	return nil
}
