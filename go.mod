module github.com/packethost/pkg

go 1.12

require (
	github.com/grpc-ecosystem/go-grpc-middleware v1.0.0
	github.com/grpc-ecosystem/go-grpc-prometheus v1.2.0
	github.com/pkg/errors v0.8.1
	github.com/prometheus/client_golang v0.9.4 // indirect
	github.com/rollbar/rollbar-go v1.0.2
	github.com/stretchr/testify v1.3.0
	go.uber.org/atomic v1.3.2 // indirect
	go.uber.org/multierr v1.1.0 // indirect
	go.uber.org/zap v1.9.1
	golang.org/x/net v0.0.0-20190213061140-3a22650c66bd // indirect
	golang.org/x/sys v0.0.0-20190221075227-b4e8571b14e0 // indirect
	google.golang.org/genproto v0.0.0-20190219182410-082222b4a5c5 // indirect
	google.golang.org/grpc v1.18.0
)

replace github.com/golang/lint => golang.org/x/lint latest
