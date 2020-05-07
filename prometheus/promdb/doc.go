/*
Package promdb is a quick library for adding prometheus metrics to your service for your databases, utilizing the sql.DBStats struct results.

Exable without labels

	db, _ := sql.Open("postgres", "")
	collector := promdb.NewDatabaseCollector(db)
	prometheus.MustRegister(collector)

Example with labels

	var labels prometheus.Labels
	var collector *promdb.DatabaseCollector

	authdb, _ := sql.Open("mysql", "")
	labels = prometheus.Labels{"driver": "mysql", "backend": "authdb"}
	collector = promdb.NewDatabaseCollectorWithLabels(authdb, labels)
	prometheus.MustRegister(collector)

	billingdb, _ := sql.Open("postgres", "")
	labels = prometheus.Labels{"driver": "postgres"}
	collector = promdb.NewDatabaseCollectorWithLabels(billingdb, labels)
	prometheus.MustRegister(collector.WithLabel("backend", "billingdb"))

	eventsdb, _ := sql.Open("postgres", "")
	labels = prometheus.Labels{"driver": "postgres", "backend": "eventsdb"}
	collector = promdb.NewDatabaseCollector(eventsdb)
	prometheus.MustRegister(collector.With(labels))

	statsdb, _ := sql.Open("bigquery", "")
	collector = promdb.NewDatabaseCollector(statsdb).WithLabel("driver", "bigquery")
	prometheus.MustRegister(collector.WithLabel("backend", "statsdb"))
*/
package promdb
