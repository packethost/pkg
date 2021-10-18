package promdb

import (
	"database/sql"
	"strings"
	"testing"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	promprotos "github.com/prometheus/client_model/go"
)

type FakeDB struct {
	MaxOpenConnections int
	OpenConnections    int
	InUse              int
	Idle               int
	WaitCount          int64
	WaitDuration       time.Duration
	MaxIdleClosed      int64
	MaxLifetimeClosed  int64
}

func (d FakeDB) Stats() sql.DBStats {
	return sql.DBStats{
		MaxOpenConnections: d.MaxOpenConnections,
		OpenConnections:    d.OpenConnections,
		InUse:              d.InUse,
		Idle:               d.Idle,
		WaitCount:          d.WaitCount,
		WaitDuration:       d.WaitDuration,
		MaxIdleClosed:      d.MaxIdleClosed,
		MaxLifetimeClosed:  d.MaxLifetimeClosed,
	}
}

func labelsMatch(a, b prometheus.Labels) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil && b != nil || a != nil && b == nil {
		return false
	}
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func TestNewDatabaseCollectorWithNoLabels(t *testing.T) {
	db := FakeDB{}

	collector := NewDatabaseCollector(db)

	if collector == nil {
		t.Fatal("got nil collector want DatabaseCollector")
	}

	if collector.db != db {
		t.Fatalf("got %+v database want %+v", collector.db, db)
	}

	if collector.labels != nil {
		t.Fatalf("got %q labels want nil", collector.labels)
	}
}

func TestNewDatabaseCollectorWithLabels(t *testing.T) {
	db := FakeDB{}

	labels := prometheus.Labels{"blarg": "ing"}
	collector := NewDatabaseCollectorWithLabels(db, labels)

	if collector == nil {
		t.Fatal("got nil collector want DatabaseCollector")
	}

	if collector.db != db {
		t.Fatalf("got %+v database want %+v", collector.db, db)
	}

	if !labelsMatch(collector.labels, labels) {
		t.Fatalf("got %q labels want %q", collector.labels, labels)
	}
}

func TestWith(t *testing.T) {
	db := FakeDB{}

	collector := NewDatabaseCollector(db).With(nil)

	if !labelsMatch(collector.labels, prometheus.Labels{}) {
		t.Errorf("got %q labels want empty map", collector.labels)
	}

	baseLabels := prometheus.Labels{
		"k1": "v1",
	}
	collector = NewDatabaseCollector(db).With(baseLabels)

	if !labelsMatch(collector.labels, baseLabels) {
		t.Errorf("got %q labels want %q", collector.labels, baseLabels)
	}

	newLabels := prometheus.Labels{
		"k2": "v2",
	}
	collector = collector.With(newLabels)

	wantLabels := prometheus.Labels{
		"k1": "v1",
		"k2": "v2",
	}
	if !labelsMatch(collector.labels, wantLabels) {
		t.Errorf("got %q labels want %q", collector.labels, wantLabels)
	}

	collector = NewDatabaseCollectorWithLabels(db, baseLabels).With(newLabels)

	if !labelsMatch(collector.labels, wantLabels) {
		t.Errorf("got %q labels want %q", collector.labels, wantLabels)
	}
}

func TestWithLabel(t *testing.T) {
	db := FakeDB{}

	baseLabels := prometheus.Labels{
		"k1": "v1",
	}
	collector := NewDatabaseCollector(db).WithLabel("k1", "v1")

	if !labelsMatch(collector.labels, baseLabels) {
		t.Errorf("got %q labels want %q", collector.labels, baseLabels)
	}

	collector = collector.WithLabel("k2", "v2")

	wantLabels := prometheus.Labels{
		"k1": "v1",
		"k2": "v2",
	}
	if !labelsMatch(collector.labels, wantLabels) {
		t.Errorf("got %q labels want %q", collector.labels, wantLabels)
	}

	collector = NewDatabaseCollectorWithLabels(db, baseLabels).WithLabel("k2", "v2")

	if !labelsMatch(collector.labels, wantLabels) {
		t.Errorf("got %q labels want %q", collector.labels, wantLabels)
	}
}

func TestDescribe(t *testing.T) {
	db := FakeDB{}
	collector := NewDatabaseCollector(db)

	finch := make(chan bool, 1)

	ch := make(chan *prometheus.Desc, 1)

	go func() {
		var desc *prometheus.Desc
		var pos int
		var metric string
		var metrics = []string{
			"database_connections_max_open",
			"database_connections_open",
			"database_connections_in_use",
			"database_connections_idle",
			"database_connections_wait_count",
			"database_connections_wait_duration_seconds",
			"database_connections_max_idle_closed",
			"database_connections_max_lifetime_closed",
		}
		for desc = range ch {
			if pos < len(metrics) {
				metric = metrics[pos] // Can't pos++ here? Weird
				pos++
				if !strings.Contains(desc.String(), "\""+metric+"\"") {
					t.Errorf("description definition does not include '\"%s\"' in '%s'", metric, desc)
				}
			} else {
				t.Errorf("An unexpected item was left on the channel '%s'", desc)
			}
		}
		finch <- true
	}()

	collector.Describe(ch)
	close(ch)

	// Wait for func to complete
	<-finch
}

func TestCollect(t *testing.T) {
	db := FakeDB{1, 2, 3, 4, 5, 6, 7, 8}
	collector := NewDatabaseCollector(db)

	finch := make(chan bool, 1)

	ch := make(chan prometheus.Metric, 1)

	go func() {
		var metric prometheus.Metric
		var protometric = &promprotos.Metric{}
		var pos int
		var value float64
		var metrics = []float64{1, 2, 3, 4, 5, 6, 7, 8}
		for metric = range ch {
			if pos < len(metrics) {
				value = metrics[pos] // Can't pos++ here? Weird
				pos++
				protometric.Reset()
				metric.Write(protometric)
				if protometric.GetGauge() == nil {
					t.Errorf("collect expected gauge metric got %v", protometric)
					continue
				}
				if protometric.GetGauge().GetValue() != value {
					t.Errorf("collect metric expected %f, got %f", value, protometric.GetGauge().GetValue())
				}
			} else {
				t.Errorf("An unexpected item was left on the channel '%s'", metric)
			}
		}
		finch <- true
	}()

	collector.Collect(ch)
	close(ch)

	// Wait for func to complete
	<-finch
}
