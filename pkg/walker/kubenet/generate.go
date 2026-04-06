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

	switch g.Config.Schedule {
	case "EasyOrbit":
		networkExperiment.Extra.SchedulePolicy = &kubenet_config.ExtraSchedulePolicy{
			EasyOrbitPolicy: &kubenet_config.EasyOrbitPolicy{
				Type: "EasyOrbit",
			},
		}
	case "Random":
	case "":
	default:
		logrus.Fatalf("unknown schedule: %s", g.Config.Schedule)
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

		// intfs
		intfs, err := g.GenerateSatelliteIntfs(IntfGenerateParam{
			RoutingAreaId: g.Config.RoutingArea,
			CrossSphere:   g.Config.CrossSphere,
			GroundAccess:  g.Config.GroundAccess,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to generate satellite intfs: %v", err)
		}

		node := kubenet_config.NetworkNode{
			Name:               nodeName,
			Type:               "Container",
			VisualOptions:      visualOption,
			Intfs:              intfs,
			MobilityModel:      mobilityModel,
			OrientationModel:   orientationModel,
			Param:              networkNodeParam,
			NetworkDeviceModel: networkDeviceModelptr,
		}

		if g.Config.Root != "" {
			node.Root = &g.Config.Root
		}

		networkNodes = append(
			networkNodes,
			node,
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

	// generate flow tasks
	switch flowPattern := g.Config.GenFlow; flowPattern {
	case FLOW_GEN_MODE_WALKER_NEIGHBOUR:
		g.GenerateNeighbourFlowTasks(constellations, networkNodes)
	case FLOW_GEN_MODE_RANDOM:
		g.GenerateRandomFlowTasks(networkNodes)
	default:
	}

	return &ConstellationResult{
		Constellations: constellations,
		NetworkNodes:   networkNodes,
	}, nil
}

func (g *KubenetGenerator) GenerateNetworkNodeParam() (*kubenet_config.NetworkNodeParam, error) {
	// 创建计算单元列表
	computationUnits := []kubenet_config.ComputationUnit{}

	// 如果配置了用户镜像，添加用户计算单元
	for _, container := range g.Config.Containers {
		computationUnits = append(computationUnits, kubenet_config.ComputationUnit{
			Name:  container.Name,
			Image: &container.Image,
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
		ModelId: "1",
	}, nil
}
