package util

type Named interface {
	GetName() string
}

func NotIn[Type Named](name string, datas []Type) bool {
	for _, data := range datas {
		if name == data.GetName() {
			return false
		}
	}
	return true
}
