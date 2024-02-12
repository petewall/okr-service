package internal

import (
	"github.com/google/uuid"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

const (
	OKRTypeBoolean    = "boolean"
	OKRTypeNumber     = "number"
	OKRTypePercentage = "percent"
)

type OKR struct {
	ID          string  `json:"id" yaml:"id"`
	Quarter     string  `json:"quarter" yaml:"quarter"`
	Category    string  `json:"category" yaml:"category"`
	ValueType   string  `json:"type" yaml:"type"`
	Description string  `json:"description" yaml:"description"`
	Goal        float64 `json:"goal" yaml:"goal"`
	Progress    float64 `json:"progress" yaml:"progress"`

	goalMetric     prometheus.Gauge
	progressMetric prometheus.Gauge
}

func CreateOKR(quarter, category, valueType, description string, goal float64) *OKR {
	okr := &OKR{
		ID:          uuid.New().String(),
		Quarter:     quarter,
		Category:    category,
		ValueType:   valueType,
		Description: description,
		Goal:        goal,
		Progress:    0,
	}
	okr.UpdateMetrics()
	return okr
}

func (okr *OKR) promLabels() prometheus.Labels {
	return prometheus.Labels{
		"quarter":     okr.Quarter,
		"category":    okr.Category,
		"type":        okr.ValueType,
		"description": okr.Description,
	}
}

func (okr *OKR) UpdateMetrics() {
	if okr.goalMetric == nil {
		okr.goalMetric = promauto.NewGauge(prometheus.GaugeOpts{
			Namespace:   "okr",
			Subsystem:   "goal",
			Name:        "info",
			Help:        "A metric representing the goal of an OKR (objective and key result). Must have a matching okr_progress_info metric.",
			ConstLabels: okr.promLabels(),
		})
	}
	okr.goalMetric.Set(okr.Goal)

	if okr.progressMetric == nil {
		okr.progressMetric = promauto.NewGauge(prometheus.GaugeOpts{
			Namespace:   "okr",
			Subsystem:   "progress",
			Name:        "info",
			Help:        "A metric representing the progress of an OKR (objective and key result). Must have a matching okr_goal_info metric.",
			ConstLabels: okr.promLabels(),
		})
	}
	okr.progressMetric.Set(okr.Progress)
}

func (okr *OKR) Set(other *OKR) {
	okr.Quarter = other.Quarter
	okr.Category = other.Category
	okr.ValueType = other.ValueType
	okr.Description = other.Description
	okr.Goal = other.Goal
	okr.Progress = other.Progress
	okr.UpdateMetrics()
}
