package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var ActiveIncidents = promauto.NewGauge(prometheus.GaugeOpts{
	Name: "warden_active_incidents_total",
	Help: "Total number of active incidents",
})
