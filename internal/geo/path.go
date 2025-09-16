package geo

import (
	"encoding/json"
	"fmt"
	"math"
	"regexp"
	"strconv"
	"strings"

	"dslab.sjtu/kubenet/api/pkg/kubenet_config"
	"gopkg.in/yaml.v2"
)

const EarthRadius = 6371000.0 // 地球半径，单位 m

// haversine 公式计算两点球面距离（单位：米）
func haversine(lat1, lon1, lat2, lon2 float64) float64 {
	dLat := (lat2 - lat1) * math.Pi / 180.0
	dLon := (lon2 - lon1) * math.Pi / 180.0

	lat1Rad := lat1 * math.Pi / 180.0
	lat2Rad := lat2 * math.Pi / 180.0

	a := math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Cos(lat1Rad)*math.Cos(lat2Rad)*
			math.Sin(dLon/2)*math.Sin(dLon/2)

	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	return EarthRadius * c
}

func ConvertPathFromSAS(pathStr string, altitude, velocity float64, outputType string) (string, error) {
	// 匹配 N/S + 数字 和 E/W + 数字
	re := regexp.MustCompile(`([NS])([\d\.]+)\s+([EW])([\d\.]+)`)
	matches := re.FindAllStringSubmatch(pathStr, -1)

	mobility_model := kubenet_config.PolylineMobilityModel{
		Type: "Polyline",
	}

	var prevLat, prevLon float64
	var cumulativeTime float64 // 纳秒累计时间

	speedMS := velocity * 1000.0 / 3600.0 // km/h -> m/s

	for i, m := range matches {
		latStr, lonStr := m[2], m[4]
		lat, _ := strconv.ParseFloat(latStr, 64)
		lon, _ := strconv.ParseFloat(lonStr, 64)

		// 判断正负
		if strings.ToUpper(m[1]) == "S" {
			lat = -lat
		}
		if strings.ToUpper(m[3]) == "W" {
			lon = -lon
		}

		// 计算与前一个点的时间
		if i > 0 {
			dist := haversine(prevLat, prevLon, lat, lon) // 米
			sec := dist / speedMS                         // 秒
			cumulativeTime += sec * 1e9                   // 纳秒
		}

		mobility_model.Sequence = append(mobility_model.Sequence, kubenet_config.PolylineMobilityModelSequenceInner{
			Latitude:  lat,
			Longitude: lon,
			Altitude:  altitude,
			Time:      int64(cumulativeTime),
		})

		prevLat, prevLon = lat, lon
	}

	switch outputType {
	case "json":
		jsonBytes, err := json.Marshal(mobility_model)
		if err != nil {
			return "", err
		}
		return string(jsonBytes), nil
	case "yaml":
		// ① 先用 json.Marshal 得到 JSON 字节
		jsonBytes, err := json.Marshal(mobility_model)
		if err != nil {
			return "", err
		}

		// ② 再把 JSON 反序列化为通用结构
		var obj interface{}
		if err := json.Unmarshal(jsonBytes, &obj); err != nil {
			return "", err
		}

		// ③ 用 yaml.Marshal 转成 YAML
		yamlBytes, err := yaml.Marshal(obj)
		if err != nil {
			return "", err
		}

		fmt.Println(string(yamlBytes))
		return string(yamlBytes), nil
	default:
		return "", fmt.Errorf("unsupported output type: %s", outputType)
	}
}
