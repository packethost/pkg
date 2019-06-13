# log

```go
import "github.com/packethost/pkg/log"
```

Package log sets up a shared logger that can be used by all packages run under one binary.

This package wraps zap very lightly so zap best practices apply here too, namely use `With` for KV pairs to add context to a line.
The lack of a wide gamut of logging levels is by design.
The intended use case for each of the levels are:
  Error:
    Logs a message as an error, may also have external side effects such as posting to rollbar, sentry or alerting directly.
  Info:
    Used for production.
    Context should all be in K=V pairs so they can be useful to ops and future-you-at-3am.
  Debug:
    Meant for developer use *during development*.

## Index

 - [type Logger](#Logger)
     - [func Init(service string) (Logger, error)](#Init)
     - [func Test(t zaptest.TestingT, service string) Logger](#Test)

#### Examples

 - [Logger (Debug)](#ExampleLogger_Debug)
 - [Logger (Error)](#ExampleLogger_Error)
 - [Logger (Fatal)](#ExampleLogger_Fatal)
 - [Logger (Info)](#ExampleLogger_Info)
 - [Logger (Package)](#ExampleLogger_Package)
 - [Logger (With)](#ExampleLogger_With)

## <a name='Logger'></a>type [Logger]()

```go
type Logger struct {
}
```

Logger is a wrapper around zap.SugaredLogger

<a name='ExampleLogger_Debug'></a><details><summary>Example (Debug)</summary><p>


```go
l := setupForExamples("debug")
defer l.Close()

l.Debug("debug message")
```

Output:
```
{"level":"debug","caller":"log/log_test.go:252","msg":"debug message","service":"github.com/packethost/pkg","pkg":"debug"}
```
</p></details>

<a name='ExampleLogger_Error'></a><details><summary>Example (Error)</summary><p>


```go
l := setupForExamples("error")
defer l.Close()

l.Error(fmt.Errorf("oh no an error"))
```

Output:
```
{"level":"error","caller":"log/log_test.go:275","msg":"oh no an error","service":"github.com/packethost/pkg","pkg":"error","error":"oh no an error"}
```
</p></details>

<a name='ExampleLogger_Fatal'></a><details><summary>Example (Fatal)</summary><p>


```go
l := setupForExamples("fatal")
defer l.Close()

defer func() {
	recover()
}()
l.Fatal(fmt.Errorf("oh no an error"))
```

Output:
```
{"level":"error","caller":"log/log_test.go:288","msg":"oh no an error","service":"github.com/packethost/pkg","pkg":"fatal","error":"oh no an error"}
```
</p></details>

<a name='ExampleLogger_Info'></a><details><summary>Example (Info)</summary><p>


```go
l := setupForExamples("info")
defer l.Close()

defer func() {
	recover()
}()
l.Info("info message")
```

Output:
```
{"level":"info","caller":"log/log_test.go:265","msg":"info message","service":"github.com/packethost/pkg","pkg":"info"}
```
</p></details>

<a name='ExampleLogger_Package'></a><details><summary>Example (Package)</summary><p>


```go
l := setupForExamples("info")
defer l.Close()

l.Info("info message")
l = l.Package("package")
l.Info("info message")
```

Output:
```
{"level":"info","caller":"log/log_test.go:308","msg":"info message","service":"github.com/packethost/pkg","pkg":"info"}
{"level":"info","caller":"log/log_test.go:310","msg":"info message","service":"github.com/packethost/pkg","pkg":"info","pkg":"package"}
```
</p></details>

<a name='ExampleLogger_With'></a><details><summary>Example (With)</summary><p>


```go
l := setupForExamples("with")
defer l.Close()

l.With("true", true).Info("info message")
```

Output:
```
{"level":"info","caller":"log/log_test.go:298","msg":"info message","service":"github.com/packethost/pkg","pkg":"with","true":true}
```
</p></details>

## <a name='Init'></a> func  [Init]()

```go
func Init(service string) (Logger, error)
```
Init initializes the logging system and sets the "service" key to the provided argument.
This func should only be called once and after flag.Parse() has been called otherwise leveled logging will not be configured correctly.

## <a name='Test'></a> func  [Test]()

```go
func Test(t zaptest.TestingT, service string) Logger
```
Test returns a logger that does not log to rollbar and can be used with testing.TB to only log on test failure or run with -v

## <a name='AddCallerSkip'></a> func (Logger) [AddCallerSkip]()

```go
func (l Logger) AddCallerSkip(skip int) Logger
```
AddCallerSkip increases the number of callers skipped by caller annotation.
When building wrappers around the Logger, supplying this option prevents Logger from always reporting the wrapper code as the caller.

## <a name='Close'></a> func (Logger) [Close]()

```go
func (l Logger) Close()
```
Close finishes and flushes up any in-flight logs

## <a name='Debug'></a> func (Logger) [Debug]()

```go
func (l Logger) Debug(args ...interface{})
```
Debug is used to log messages in development, not even for lab.
No one cares what you pass to Debug.
All the values of arg are stringified and concatenated without any strings.

## <a name='Error'></a> func (Logger) [Error]()

```go
func (l Logger) Error(err error, args ...interface{})
```
Error is used to log an error, the error will be forwared to rollbar and/or other external services.
All the values of arg are stringified and concatenated without any strings.
If no args are provided err.Error() is used as the log message.

## <a name='Fatal'></a> func (Logger) [Fatal]()

```go
func (l Logger) Fatal(err error, args ...interface{})
```
Fatal calls Error followed by a panic(err)

## <a name='GRPCLoggers'></a> func (Logger) [GRPCLoggers]()

```go
func (l Logger) GRPCLoggers() (grpc.StreamServerInterceptor, grpc.UnaryServerInterceptor)
```
GRPCLoggers returns server side logging middleware for gRPC servers

## <a name='Info'></a> func (Logger) [Info]()

```go
func (l Logger) Info(args ...interface{})
```
Info is used to log message in production, only simple strings should be given in the args.
Context should be added as K=V pairs using the `With` method.
All the values of arg are stringified and concatenated without any strings.

## <a name='Package'></a> func (Logger) [Package]()

```go
func (l Logger) Package(pkg string) Logger
```
Package returns a copy of the logger with the "pkg" set to the argument.
It should be called before the original Logger has had any keys set to values, otherwise confusion may ensue.

## <a name='With'></a> func (Logger) [With]()

```go
func (l Logger) With(args ...interface{}) Logger
```
With is used to add context to the logger, a new logger copy with the new K=V pairs as context is returned.
