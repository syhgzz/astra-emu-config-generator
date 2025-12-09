package kubenet

import (
	"fmt"

	"dslab.sjtu/kubenet-config-generator/pkg/kubenet_config_util"
	"dslab.sjtu/kubenet-config-generator/pkg/walker/tle"
	"dslab.sjtu/kubenet/api/pkg/kubenet_config"
	"github.com/sirupsen/logrus"
)

type ConstellationResult struct {
	Constellations []kubenet_config.Constellation
	NetworkNodes   []kubenet_config.NetworkNode
}

func (g *KubenetGenerator) GenerateKubenetConfig() error {
	// read template
	networkExperiment, err := kubenet_config_util.ReadConfig(g.TemplatePath)
	if err != nil {
		return fmt.Errorf("failed to read networkExperiment: %v", err)
	}

	// generate networknodes and constellations
	constellationsResult, err := g.GenerateConstellation()
	if err != nil {
		return fmt.Errorf("failed to generate constellation: %v", err)
	}

	networkExperiment.Constellations = append(networkExperiment.Constellations, constellationsResult.Constellations...)
	networkExperiment.NetworkNodes = append(networkExperiment.NetworkNodes, constellationsResult.NetworkNodes...)

	// set isFakeTest
	if g.Config.Fake {
		isFake := true
		networkExperiment.Extra.IsFakeTest = &isFake
	}

	// write to output
	err = kubenet_config_util.WriteConfig(g.OutputPath, networkExperiment)
	if err != nil {
		logrus.Fatalf("failed to write networkExperiment: %v", err)
	}
	return nil
}

func (g *KubenetGenerator) GenerateConstellation() (*ConstellationResult, error) {
	// 读取TLE数据
	phyConstellation, err := tle.NewConstellation(g.InputPath)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize constellation: %v", err)
	}

	// 计算卫星分布
	satNumPerTrack := g.Config.SatelliteNum / g.Config.PlaneNum
	actualNum := satNumPerTrack * g.Config.NeededPlane
	isFull := g.Config.PlaneNum == g.Config.NeededPlane

	// 初始化数据结构
	networkNodes := []kubenet_config.NetworkNode{}
	var constellations []kubenet_config.Constellation

	// network device model
	var networkDeviceModelptr *string
	if g.Config.NetworkDeviceModel != "" {
		networkDeviceModelptr = &g.Config.NetworkDeviceModel
	} else {
		networkDeviceModelptr = nil
	}

	// networknode param
	networkNodeParam, err := g.GenerateNetworkNodeParam()
	if err != nil {
		return nil, fmt.Errorf("failed to generate network node param: %v", err)
	}

	// visual option
	visualOption, err := g.GenerateVisualOption()
	if err != nil {
		return nil, fmt.Errorf("failed to generate visual option: %v", err)
	}

	// 4. 创建卫星和网络节点
	satellites := []kubenet_config.OrbitSatellite{}
	planes := []kubenet_config.OrbitPlane{}
	planeIdx := 0
	innerIdx := 0
	for idx := range actualNum {
		nodeName := g.Config.NodeNamePrefix + phyConstellation[idx].Name
		satellites = append(
			satellites,
			kubenet_config.OrbitSatellite{
				Idx:       int32(innerIdx),
				Satellite: nodeName,
			},
		)
		innerIdx += 1
		if (innerIdx%satNumPerTrack == 0) && (innerIdx >= satNumPerTrack) {
			innerIdx -= satNumPerTrack
			planes = append(
				planes,
				kubenet_config.OrbitPlane{
					Idx:        int32(planeIdx),
					Satellites: satellites,
				},
			)
			planeIdx += 1
			satellites = []kubenet_config.OrbitSatellite{}
		}

		// mobility model
		mobilityModel := kubenet_config.NetworkNodeMobilityModel{
			SatelliteMobilityModel: &kubenet_config.SatelliteMobilityModel{
				Type:     "Satellite",
				TleLine1: phyConstellation[idx].Satellite.Line1,
				TleLine2: phyConstellation[idx].Satellite.Line2,
			},
		}

		// orientation model
		orientationModel := kubenet_config.NetworkNodeOrientationModel{
			WalkerSatelliteOrientationModel: &kubenet_config.WalkerSatelliteOrientationModel{
				Type: "WalkerSatellite",
			},
		}

		// declare flowTasks
		var flowTasks []kubenet_config.TaskSequence
		if g.Config.GenFlow {
			// generate flow tasks
			flowTasks = g.GenerateFlowTasks(nodeName, idx, satNumPerTrack, g.Config.NeededPlane, phyConstellation)
		}

		// intfs
		intfs, err := g.GenerateSatelliteIntfs(IntfGenerateParam{
			RoutingAreaId: g.Config.RoutingArea,
			CrossSphere:   g.Config.CrossSphere,
			GroundAccess:  g.Config.GroundAccess,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to generate satellite intfs: %v", err)
		}

		//networknode
		networkNodes = append(
			networkNodes,
			kubenet_config.NetworkNode{
				Name:               nodeName,
				Type:               "Container",
				VisualOptions:      visualOption,
				Intfs:              intfs,
				MobilityModel:      mobilityModel,
				OrientationModel:   orientationModel,
				Param:              networkNodeParam,
				NetworkDeviceModel: networkDeviceModelptr,
				Tasks:              flowTasks,
			},
		)
	}
	constellations = []kubenet_config.Constellation{
		{
			Name: fmt.Sprintf("s%d", g.Config.ShellIdx),
			Spheres: []kubenet_config.ConstellationSpheresInner{
				{
					WalkerSphere: &kubenet_config.WalkerSphere{
						Type:   "Walker",
						Planes: planes,
						Phase:  int32(g.Config.Phase),
						Offset: g.Config.Offset,
						IsFull: isFull,
					},
				},
			},
		},
	}

	return &ConstellationResult{
		Constellations: constellations,
		NetworkNodes:   networkNodes,
	}, nil
}

func (g *KubenetGenerator) GenerateNetworkNodeParam() (*kubenet_config.NetworkNodeParam, error) {
	// 创建计算单元列表
	testSuiteImage := TEST_SUITE_IMAGE
	computationUnits := []kubenet_config.ComputationUnit{
		{
			Name:  "test-suite",
			Image: &testSuiteImage,
			Ports: TEST_SUITE_PORTS,
		},
	}

	// 如果配置了用户镜像，添加用户计算单元
	if g.Config.Image != "" {
		computationUnits = append(computationUnits, kubenet_config.ComputationUnit{
			Name:  "user-image",
			Image: &g.Config.Image,
		})
	}

	// 创建并返回网络节点参数
	return &kubenet_config.NetworkNodeParam{
		ContainerParam: &kubenet_config.ContainerParam{
			Type:             "Container",
			ComputationUnits: computationUnits,
		},
	}, nil
}

func (g *KubenetGenerator) GenerateVisualOption() (*kubenet_config.VisualOption, error) {
	return &kubenet_config.VisualOption{
		ModelId:  "1",
		NodeType: 1,
	}, nil
}

func (g *KubenetGenerator) GenerateFlowTasks(currentNodeName string, currentPhyIdx, satNumPerTrack, totalPlanes int, phyConstellation []tle.Satellite) []kubenet_config.TaskSequence {
	tasks := []kubenet_config.TaskSequence{}

	// 计算目标卫星名称：通过 phyConstellation 中的真实名称获取，减少拼接出错风险
	// 将物理索引转换为轨道面索引和同轨道内卫星索引
	currentPlaneIdx := 0
	currentSatIdx := 0
	if satNumPerTrack > 0 {
		currentPlaneIdx = currentPhyIdx / satNumPerTrack
		currentSatIdx = currentPhyIdx % satNumPerTrack
	}

	// 1. 轨道面内目标：同一轨道面，卫星索引+1
	intraPhyIdx := currentPlaneIdx*satNumPerTrack + ((currentSatIdx + 1) % satNumPerTrack)
	var intraPlaneTargetName string
	if intraPhyIdx >= 0 && intraPhyIdx < len(phyConstellation) {
		intraPlaneTargetName = g.Config.NodeNamePrefix + phyConstellation[intraPhyIdx].Name
	}

	// 2. 轨道面间目标：下一轨道面，相同卫星索引
	interPlaneTargetPlaneIdx := 0
	if totalPlanes > 0 {
		interPlaneTargetPlaneIdx = (currentPlaneIdx + 1) % totalPlanes
	}
	interPhyIdx := interPlaneTargetPlaneIdx*satNumPerTrack + currentSatIdx
	var interPlaneTargetName string
	if interPhyIdx >= 0 && interPhyIdx < len(phyConstellation) {
		interPlaneTargetName = g.Config.NodeNamePrefix + phyConstellation[interPhyIdx].Name
	}

	// 创建轨道面内流量任务
	backend := "iperf"
	intraPlaneTask := kubenet_config.TaskSequence{
		Uid: fmt.Sprintf("%s-intra", currentNodeName),
		Sequence: []kubenet_config.Task{
			{
				Param: kubenet_config.TaskParam{
					FlowTaskParam: &kubenet_config.FlowTaskParam{
						Type:      "Flow",
						Dst:       intraPlaneTargetName,
						Bandwidth: FLOW_TASK_BANDWIDTH,
						Backend:   &backend,
					},
				},
				// StartTime 和 StopTime 留空
			},
		},
	}

	// 创建轨道面间流量任务
	interPlaneTask := kubenet_config.TaskSequence{
		Uid: fmt.Sprintf("%s-inter", currentNodeName),
		Sequence: []kubenet_config.Task{
			{
				Param: kubenet_config.TaskParam{
					FlowTaskParam: &kubenet_config.FlowTaskParam{
						Type:      "Flow",
						Dst:       interPlaneTargetName,
						Bandwidth: FLOW_TASK_BANDWIDTH,
						Backend:   &backend,
					},
				},
				// StartTime 和 StopTime 留空
			},
		},
	}

	tasks = append(tasks, intraPlaneTask, interPlaneTask)
	return tasks
}
