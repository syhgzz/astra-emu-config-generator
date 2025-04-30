package tle

import (
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
	"time"
)

// TLEItem represents a satellite in TLE file
type TLEItem struct {
	Name string

	// params for line 1
	SatId         int     // 卫星ID
	Class         string  // 卫星类型
	LaunchYear    int     // 国际代号年份
	LaunchSeq     int     // 国际代号年份发射编号
	LaunchNum     string  // 国际代号其余部份
	EpochYear     int     // TLE 历年年份
	Epoch         float64 // TLE 历年（一年中的第 N 天，含小数）
	MeanMotion1   float64 // 平均运动一阶导数
	MeanMotion2   float64 // 平均运动二阶导数
	DragTerm      float64 // 摄动项
	EphemerisType int     // 星历类型，永远为零
	ElementSetNum int     // 元素集编号

	// params for line 2
	// SatId can be reused for line 2
	Inclination   float64 // 轨道倾角
	Raan          float64 // 升交点赤经
	Eccentricity  float64 // 偏心率
	ArgOfPerigee  float64 // 近地点幅角
	MeanAnomaly   float64 // 平近点角
	MeanMotion    float64 // 平均运动 - 每天的转数
	RevNumAtEpoch int     // 轨道周期内转数
}

type TLEItems []TLEItem

func NewTLEItem() *TLEItem {
	now := time.Now()
	epochYear, epoch := ConvertTimeToTLEEpoch(now)

	return &TLEItem{
		// line 1
		Class:         "U",
		LaunchYear:    24,
		LaunchSeq:     1,
		LaunchNum:     "A",
		EpochYear:     epochYear,
		Epoch:         epoch,
		MeanMotion1:   0.00001478,
		MeanMotion2:   0,
		DragTerm:      0.00011810,
		EphemerisType: 0,
		ElementSetNum: 1,
		// line 2
		Eccentricity:  0.0001424,
		ArgOfPerigee:  94.3232,
		MeanAnomaly:   265.7920,
		MeanMotion:    15.06405745,
		RevNumAtEpoch: 0,
	}
}

func (t *TLEItem) String() string {
	meanMotion1Str := ""
	if t.MeanMotion1 < 0 {
		meanMotion1Str = "-"
		meanMotion1Str += "." + ConvertToDotfree(t.MeanMotion1, 8)
	} else {
		meanMotion1Str += "." + ConvertToDotfree(t.MeanMotion1, 9)
	}
	line1 := fmt.Sprintf("1 %05d%1s %02d%03d%-3s %0d%012.8f %10s %7s %8s %01d %04d",
		t.SatId,
		t.Class,
		t.LaunchYear,
		t.LaunchSeq,
		t.LaunchNum,
		t.EpochYear,
		t.Epoch,
		meanMotion1Str,
		ConvertToSpecialExponential(t.MeanMotion2, 7),
		ConvertToSpecialExponential(t.DragTerm, 8),
		t.EphemerisType,
		t.ElementSetNum)

	line2 := fmt.Sprintf("2 %05d %08.4f %08.4f %07s %08.4f %08.4f %011.7f%05d",
		t.SatId,
		t.Inclination,
		t.Raan,
		ConvertToDotfree(t.Eccentricity, 7),
		t.ArgOfPerigee,
		t.MeanAnomaly,
		t.MeanMotion,
		t.RevNumAtEpoch)

	// Calculate and append checksums
	line1 = line1 + fmt.Sprintf("%d", calculateChecksum(line1))
	line2 = line2 + fmt.Sprintf("%d", calculateChecksum(line2))

	content := ""
	if t.Name != "" {
		content = t.Name + "\n"
	}
	content += line1 + "\n" + line2 + "\n"
	return content
}

func (t *TLEItems) WriteFile(filename string) error {
	content := ""
	for _, item := range *t {
		content += item.String()
	}
	err := os.WriteFile(filename, []byte(content), 0644)
	if err != nil {
		return err
	}
	return nil
}

// CalculateOrbitalPeriod 根据高度计算卫星运动周期（假设轨道是圆的）
func CalculateOrbitalPeriod(h float64) float64 {
	const (
		earthRadius = 6378000.0   // 地球平均半径（米）
		gravConst   = 6.67430e-11 // 万有引力常数（N·m²/kg²）
		earthMass   = 5.972e24    // 地球质量（kg）
	)

	// 计算轨道半径（地球半径 + 卫星高度）
	r := earthRadius + h

	// 使用开普勒第三定律公式：T = 2π√(r³/GM)
	period := 2 * math.Pi * math.Sqrt(math.Pow(r, 3)/(gravConst*earthMass))

	return period
}

type TLEParams struct {
	Height      float64
	Name        string
	Num         int
	Plane       int
	Phase       float64
	Inclination float64

	// Plane Delta Start
	DeltaInit float64

	// Plane num needed, 0 means all planes
	NeededPlane int
}

type TLEItemOverrideFunc func(item *TLEItem, planeId int, satId int)

func (t *TLEParams) NewTLEWithOverrides(overrideFunc TLEItemOverrideFunc) TLEItems {
	delta := 360.0 / float64(t.Plane)
	satNumPerTrack := t.Num / t.Plane
	phaseGap := 360.0 / float64(satNumPerTrack)
	items := make(TLEItems, 0)
	var planeNum int
	if t.NeededPlane > 0 {
		planeNum = t.NeededPlane
	} else {
		planeNum = t.Plane
	}
	for planeId := range planeNum {
		raan := math.Mod(t.DeltaInit+delta*float64(planeId)+360.0, 360.0)
		for satId := range satNumPerTrack {
			item := NewTLEItem()
			item.SatId = 10000 + 100*planeId + satId
			item.Name = strings.ToUpper(t.Name) + "-" + strconv.Itoa(item.SatId)
			item.Raan = raan
			item.Inclination = t.Inclination
			// Reference: https://ww2.mathworks.cn/help/aerotbx/ug/satellitescenario.walkerdelta.html#mw_badf1acd-4f40-4ad6-987e-c7d35ae2368c
			// item.MeanAnomaly = phaseGap*float64(satId) + t.Phase*360.0/float64(t.Num)*float64(satId)
			item.MeanAnomaly = phaseGap*float64(satId) + 360.0/float64(t.Plane)*t.Phase*float64((planeId)%t.Plane)
			item.MeanMotion = 24 * 3600 / CalculateOrbitalPeriod(t.Height*1000)
			if overrideFunc != nil {
				overrideFunc(item, planeId, satId)
			}
			items = append(items, *item)
		}
	}

	return items
}

func (t *TLEParams) NewTLE() TLEItems {
	return t.NewTLEWithOverrides(nil)
}

// ConvertTimeToTLEEpoch converts current time to TLE epoch format
func ConvertTimeToTLEEpoch(t time.Time) (epochYear int, epoch float64) {
	epochYear = t.Year() % 100

	dayOfYear := float64(t.YearDay() - 1)

	hours := float64(t.Hour())
	minutes := float64(t.Minute())
	seconds := float64(t.Second())

	fractionalDay := (hours + minutes/60.0 + seconds/3600.0) / 24.0

	epoch = dayOfYear + fractionalDay

	return epochYear, epoch
}
