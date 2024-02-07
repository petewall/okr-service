package internal

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

type OKR struct {
	Quarter     string  `json:"quarter"`
	Category    string  `json:"category"`
	ValueType   string  `json:"type"`
	Description string  `json:"description"`
	Goal        float64 `json:"goal"`
	Progress    float64 `json:"progress"`

	goalMetric     prometheus.Gauge
	progressMetric prometheus.Gauge
}

func CreateOKR(quarter, category, valueType, description string, goal float64) *OKR {
	okr := &OKR{
		Quarter:     quarter,
		Category:    category,
		ValueType:   valueType,
		Description: description,
		Goal:        goal,
		Progress:    0,
		goalMetric: promauto.NewGauge(prometheus.GaugeOpts{
			Namespace:   "okr",
			Subsystem:   "goal",
			Name:        "info",
			Help:        "A metric representing the goal of an OKR (objective and key result). Must have a matching okr_progress_info metric.",
			ConstLabels: prometheus.Labels{"quarter": quarter, "category": category, "type": valueType, "description": description},
		}),
		progressMetric: promauto.NewGauge(prometheus.GaugeOpts{
			Namespace:   "okr",
			Subsystem:   "progress",
			Name:        "info",
			Help:        "A metric representing the progress of an OKR (objective and key result). Must have a matching okr_goal_info metric.",
			ConstLabels: prometheus.Labels{"quarter": quarter, "category": category, "type": valueType, "description": description},
		}),
	}
	okr.goalMetric.Set(goal)
	okr.progressMetric.Set(0)

	return okr
}
