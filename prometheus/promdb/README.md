# PromDB
A quick library for adding prometheus metrics to your service for your databases

## Usage

### Without labels

```go
db, _ := sql.Open("postgres", "")
collector := promdb.NewDatabaseCollector(db)
prometheus.MustRegister(collector)
```

### With labels

```go
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
```

### Example rendered metrics

#### Without labels
```
# HELP database_connections_idle The number of idle connections.
# TYPE database_connections_idle gauge
database_connections_idle 0
# HELP database_connections_in_use The number of connections currently in use.
# TYPE database_connections_in_use gauge
database_connections_in_use 0
```

#### With labels
```
# HELP database_connections_idle The number of idle connections.
# TYPE database_connections_idle gauge
database_connections_idle{backend="authdb",driver="mysql"} 0
database_connections_idle{backend="billingdb",driver="postgres"} 0
database_connections_idle{backend="eventsdb",driver="postgres"} 0
database_connections_idle{backend="statsdb",driver="bigquery"} 0
# HELP database_connections_in_use The number of connections currently in use.
# TYPE database_connections_in_use gauge
database_connections_in_use{backend="authdb",driver="mysql"} 0
database_connections_in_use{backend="billingdb",driver="postgres"} 0
database_connections_in_use{backend="eventsdb",driver="postgres"} 0
database_connections_in_use{backend="statsdb",driver="bigquery"} 0
```