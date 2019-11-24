package signalprocessing

// ChangeDetector types calculate change points.
type ChangeDetector interface {
	DetectChanges([]float64) ([]ChangePoint, error)
}

type ChangePoint struct {
	Index int           `bson:"index" json:"index" yaml:"index"`
	Info  AlgorithmInfo `bson:"info" json:"info" yaml:"info"`
}

type AlgorithmInfo struct {
	Name    string            `bson:"name" json:"name" yaml:"name"`
	Version int               `bson:"version" json:"version" yaml:"version"`
	Options []AlgorithmOption `bson:"options" json:"options" yaml:"options"`
}

type AlgorithmOption struct {
	Name  string      `bson:"name" json:"name" yaml:"name"`
	Value interface{} `bson:"value" json:"value" yaml:"value"`
}
