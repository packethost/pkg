# grpc

```go
import "github.com/packethost/pkg/grpc"
```


## Index

 - [type Option](#Option)
     - [func LoadX509KeyPair(certFile, keyFile string) Option](#LoadX509KeyPair)
     - [func Port(port int) Option](#Port)
     - [func Register(r func(*grpc.Server)) Option](#Register)
     - [func ServerOption(opt grpc.ServerOption) Option](#ServerOption)
     - [func StreamInterceptor(si grpc.StreamServerInterceptor) Option](#StreamInterceptor)
     - [func UnaryInterceptor(ui grpc.UnaryServerInterceptor) Option](#UnaryInterceptor)
     - [func X509KeyPair(certPEMBlock, keyPEMBlock string) Option](#X509KeyPair)
 - [type Server](#Server)
     - [func NewServer(l log.Logger, reg ServiceRegister, options ...Option) (*Server, error)](#NewServer)
 - [type ServiceRegister](#ServiceRegister)

## <a name='Option'></a>type [Option]()

```go
type Option func(*Server)
```

The Option type describes functions that operate on Server during NewServer.
It is a convenience type to make it easier for callers to build up the slice of options apart from the call to NewServer.

## <a name='LoadX509KeyPair'></a> func  [LoadX509KeyPair]()

```go
func LoadX509KeyPair(certFile, keyFile string) Option
```
LoadX509KeyPair will setup server as a secure server by reading the cert and key from the provided file locations.
This function overrides GRPC_CERT and GRPC_KEY environment variables.
NewServer will return an error if both X509KeyPair and LoadX509KeyPair are used.

## <a name='Port'></a> func  [Port]()

```go
func Port(port int) Option
```
Port will set the port the server will bind to, Port must be > 0

## <a name='Register'></a> func  [Register]()

```go
func Register(r func(*grpc.Server)) Option
```
Register will call the callback func after the main grpc service has been setup.
The Prometheus register is always included in the set

## <a name='ServerOption'></a> func  [ServerOption]()

```go
func ServerOption(opt grpc.ServerOption) Option
```
ServerOption will add the opt param to the underlying grpc.NewServer() call.

## <a name='StreamInterceptor'></a> func  [StreamInterceptor]()

```go
func StreamInterceptor(si grpc.StreamServerInterceptor) Option
```
StreamInterceptor adds the argument to the list of interceptors in a grpc_middleware.Chain
Logging and Prometheus interceptors are always included in the set

## <a name='UnaryInterceptor'></a> func  [UnaryInterceptor]()

```go
func UnaryInterceptor(ui grpc.UnaryServerInterceptor) Option
```
UnaryInterceptor adds the argument to the list of interceptors in a grpc_middleware.Chain
Logging and Prometheus interceptors are always included in the set

## <a name='X509KeyPair'></a> func  [X509KeyPair]()

```go
func X509KeyPair(certPEMBlock, keyPEMBlock string) Option
```
X509KeyPair will setup server as a secure server using the provided cert and key
This function overrides GRPC_CERT and GRPC_KEY environment variables.
NewServer will return an error if both X509KeyPair and LoadX509KeyPair are used.

## <a name='Server'></a>type [Server]()

```go
type Server struct {
}
```

Server is used to hold configured info and ultimately the grpc server

## <a name='NewServer'></a> func  [NewServer]()

```go
func NewServer(l log.Logger, reg ServiceRegister, options ...Option) (*Server, error)
```
NewServer creates a new grpc server.
By default the server will be an insecure server listening on port 8080 with logging and prometheus interceptors set up.

The server's port is configured via the GRPC_PORT env variable, but can be overriden by the Port helper func.
A tls server is setup if keys are provided in either the environment variables GRPC_CERT/GRPC_KEY, or using the X509KeyPair or LoadX509KeyPair helper funcs.
Logging is always setup using the provided log.Logger.
Prometheus is always setup using the default prom interceptors and Register func.

req is called after the server has been setup.
This is where your service is gets registered with grpc, equivalent to pb.RegisterMyServiceServer(s, &myServiceImpl{}).

After your service has been registered any callbacks that were setup with Register will be called to finish up registration.

## <a name='Port'></a> func (*Server) [Port]()

```go
func (s *Server) Port() int
```
Port returns the port the server is listening on

## <a name='Serve'></a> func (*Server) [Serve]()

```go
func (s *Server) Serve() error
```
Serve starts the grpc server

## <a name='Server'></a> func (*Server) [Server]()

```go
func (s *Server) Server() *grpc.Server
```
Server returns the grpc server

## <a name='ServiceRegister'></a>type [ServiceRegister]()

```go
type ServiceRegister func(*Server)
```

The ServiceRegister type is used as a callback once the underlying grpc server is setup to register the main service.
