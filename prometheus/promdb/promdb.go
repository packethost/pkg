package promdb

import (
	"database/sql"

	"github.com/prometheus/client_golang/prometheus"
)

// DBWithStats is an interface used by DatabaseCollector
type DBWithStats interface {
	Stats() sql.DBStats
}

// DatabaseCollector implements prometheus.Collector interface
type DatabaseCollector struct {
	db     DBWithStats
	labels prometheus.Labels

	maxOpenConnections *prometheus.Desc
	openConnections    *prometheus.Desc
	inUse              *prometheus.Desc
	idle               *prometheus.Desc
	waitCount          *prometheus.Desc
	waitDuration       *prometheus.Desc
	maxIdleClosed      *prometheus.Desc
	maxLifetimeClosed  *prometheus.Desc
}

// With returns a new instance of DatabaseCollector with the specified labels merging with any previous labels that may have been defined
func (c *DatabaseCollector) With(labels prometheus.Labels) *DatabaseCollector {
	// Initialize with at least as many labels as we already have, the provided
	// labels could overwrite current labels, so we can't be sure of the ending
	// size.
	newLabels := make(prometheus.Labels, len(c.labels))
	for k, v := range c.labels {
		newLabels[k] = v
	}
	for k, v := range labels {
		newLabels[k] = v
	}
	return NewDatabaseCollectorWithLabels(c.db, newLabels)
}

// WithLabel is a shortcut to With for providing a single label and value
func (c *DatabaseCollector) WithLabel(name, value string) *DatabaseCollector {
	labels := prometheus.Labels{name: value}
	return c.With(labels)
}

// Describe implements prometheus.Collector interface
func (c *DatabaseCollector) Describe(ch chan<- *prometheus.Desc) {
	prometheus.DescribeByCollect(c, ch)
}

// Collect implements prometheus.Collector interface
func (c *DatabaseCollector) Collect(ch chan<- prometheus.Metric) {
	stats := c.db.Stats()

	var maxOpenConnections float64 = float64(stats.MaxOpenConnections)
	var openConnections float64 = float64(stats.OpenConnections)
	var inUse float64 = float64(stats.InUse)
	var idle float64 = float64(stats.Idle)
	var waitCount float64 = float64(stats.WaitCount)
	var waitDuration float64 = float64(stats.WaitDuration)
	var maxIdleClosed float64 = float64(stats.MaxIdleClosed)
	var maxLifetimeClosed float64 = float64(stats.MaxLifetimeClosed)

	ch <- prometheus.MustNewConstMetric(c.maxOpenConnections, prometheus.GaugeValue, maxOpenConnections)
	ch <- prometheus.MustNewConstMetric(c.openConnections, prometheus.GaugeValue, openConnections)
	ch <- prometheus.MustNewConstMetric(c.inUse, prometheus.GaugeValue, inUse)
	ch <- prometheus.MustNewConstMetric(c.idle, prometheus.GaugeValue, idle)
	ch <- prometheus.MustNewConstMetric(c.waitCount, prometheus.GaugeValue, waitCount)
	ch <- prometheus.MustNewConstMetric(c.waitDuration, prometheus.GaugeValue, waitDuration)
	ch <- prometheus.MustNewConstMetric(c.maxIdleClosed, prometheus.GaugeValue, maxIdleClosed)
	ch <- prometheus.MustNewConstMetric(c.maxLifetimeClosed, prometheus.GaugeValue, maxLifetimeClosed)
}

// NewDatabaseCollector creates a new instance of DatabaseCollector with no labels
func NewDatabaseCollector(db DBWithStats) *DatabaseCollector {
	return NewDatabaseCollectorWithLabels(db, nil)
}

// NewDatabaseCollectorWithLabels creates a new instance of DatabaseCollector with the specified prometheus labels
func NewDatabaseCollectorWithLabels(db DBWithStats, labels prometheus.Labels) *DatabaseCollector {
	return &DatabaseCollector{
		db:     db,
		labels: labels,

		maxOpenConnections: prometheus.NewDesc("database_connections_max_open",
			"Maximum number of open connections to the database.",
			nil, labels,
		),
		openConnections: prometheus.NewDesc("database_connections_open",
			"The number of established connections both in use and idle.",
			nil, labels,
		),
		inUse: prometheus.NewDesc("database_connections_in_use",
			"The number of connections currently in use.",
			nil, labels,
		),
		idle: prometheus.NewDesc("database_connections_idle",
			"The number of idle connections.",
			nil, labels,
		),
		waitCount: prometheus.NewDesc("database_connections_wait_count",
			"The total number of connections waited for.",
			nil, labels,
		),
		waitDuration: prometheus.NewDesc("database_connections_wait_duration_seconds",
			"The total time blocked waiting for a new connection.",
			nil, labels,
		),
		maxIdleClosed: prometheus.NewDesc("database_connections_max_idle_closed",
			"The total number of connections closed due to SetMaxIdleConns.",
			nil, labels,
		),
		maxLifetimeClosed: prometheus.NewDesc("database_connections_max_lifetime_closed",
			"The total number of connections closed due to SetConnMaxLifetime.",
			nil, labels,
		),
	}
}
