package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/qonto/standards-insights/checks/aggregates"
)

type Project struct {
	registry    *prometheus.Registry
	checksGauge *prometheus.GaugeVec
	extraLabels []string
}

func New(registry *prometheus.Registry, extraLabels []string) (*Project, error) {
	labels := append([]string{"name", "project"}, extraLabels...)
	checksGauge := prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "check_result_success",
			Help: "Projects checks results",
		},
		labels)
	err := registry.Register(checksGauge)
	if err != nil {
		return nil, err
	}
	return &Project{
		registry:    registry,
		checksGauge: checksGauge,
		extraLabels: extraLabels,
	}, nil
}

func (p *Project) LoadProjectsMetrics(results []aggregates.ProjectResult) {
	for _, project := range results {
		for _, result := range project.CheckResults {
			labels := prometheus.Labels{"name": result.Name, "project": project.Name}
			for _, label := range p.extraLabels {
				value, ok := result.Labels[label]
				if ok {
					labels[label] = value
				} else {
					labels[label] = ""
				}
			}
			gaugeValue := 0
			if result.Success {
				gaugeValue = 1
			}
			p.checksGauge.With(labels).Set(float64(gaugeValue))
		}
	}
}
