package kubenet

import (
	"fmt"
	"math/rand"
	"time"

	"dslab.sjtu/kubenet/api/pkg/kubenet_config"
)

// add flow tasks to each networknode
// each networknode will get 2 flow tasks if possible: intra-plane task and inter-plane task
// intra-plane: to the next satellite in the same plane
// inter-plane: to the satellite in the next plane with the same index
// 预设这里constellations内部所有层级元素的idx都是有序且+1递增的，不然会有问题
func (g *KubenetGenerator) GenerateNeighbourFlowTasks(constellations []kubenet_config.Constellation, networknodes []kubenet_config.NetworkNode) {
	// 将networknodes转换成map
	networknodesMap := make(map[string]*kubenet_config.NetworkNode)
	for i := range networknodes {
		node := &networknodes[i]
		networknodesMap[node.Name] = node
	}

	// 遍历constellations，找到每个卫星的idx+1和plane idx+1的邻居并创建flow任务
	for _, c := range constellations {
		for _, s := range c.Spheres {
			w := s.WalkerSphere
			if w == nil || len(w.Planes) == 0 {
				continue
			}
			for i := range w.Planes {
				currPlane := &w.Planes[i]
				nextPlane := &w.Planes[(i+1)%len(w.Planes)]
				if len(currPlane.Satellites) == 0 {
					continue
				}

				for j := range currPlane.Satellites {
					currSat := &currPlane.Satellites[j]
					// find networknode
					networknode, ok := networknodesMap[currSat.Satellite]
					if !ok {
						continue
					}

					taskSequence := kubenet_config.TaskSequence{
						Uid: fmt.Sprintf("%s-gen-flow", currSat.Satellite),
					}

					// generate intra-plane flow task
					if len(currPlane.Satellites) > 1 {
						nextSat := &currPlane.Satellites[(j+1)%len(currPlane.Satellites)]
						if nextSat.Satellite != currSat.Satellite {
							taskSequence.Sequence = append(taskSequence.Sequence, makeFlowTask(nextSat.Satellite))
						}
					}

					// generate inter-plane flow task
					if j < len(nextPlane.Satellites) && nextPlane.Idx != currPlane.Idx {
						offset := g.Config.Offset
						if i == g.Config.PlaneNum {
							offset += int32(g.Config.Phase)
						}

						nextIdx := (j + int(offset))
						if nextIdx < 0 {
							nextIdx += len(nextPlane.Satellites)
						}
						nextIdx = nextIdx % len(nextPlane.Satellites)

						nextPlaneSat := &nextPlane.Satellites[nextIdx]
						taskSequence.Sequence = append(taskSequence.Sequence, makeFlowTask(nextPlaneSat.Satellite))
					}

					// add task to networknode
					if len(taskSequence.Sequence) > 0 {
						networknode.Tasks = append(networknode.Tasks, taskSequence)
					}
				}
			}
		}
	}
}

func (g *KubenetGenerator) GenerateRandomFlowTasks(
	networknodes []kubenet_config.NetworkNode,
) {
	if len(networknodes) < 2 {
		return
	}

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	n := len(networknodes)

	for i := range networknodes {
		src := &networknodes[i]

		j := rng.Intn(n - 1)
		if j >= i {
			j++
		}
		dst := networknodes[j].Name

		taskSeq := kubenet_config.TaskSequence{
			Uid: fmt.Sprintf("%s-gen-flow", src.Name),
			Sequence: []kubenet_config.Task{
				makeFlowTask(dst),
			},
		}

		src.Tasks = append(src.Tasks, taskSeq)
	}
}

func makeFlowTask(dst string) kubenet_config.Task {
	return kubenet_config.Task{
		Param: kubenet_config.TaskParam{
			FlowTaskParam: &kubenet_config.FlowTaskParam{
				Type:      "Flow",
				Dst:       dst,
				Bandwidth: FLOW_TASK_BANDWIDTH,
			},
		},
	}
}
