# env

```go
import "github.com/packethost/pkg/env"
```


## Index

 - [func Get(name string, def ...string) string](#Get)

#### Examples

 - [Get](#ExampleGet)

## <a name='Get'></a> func  [Get]()

```go
func Get(name string, def ...string) string
```
Get retrieves the value of the environment variable named by the key.
If the value is empty or unset it will return the first value of def or "" if none is given

<a name='ExampleGet'></a><details><summary>Example</summary><p>


```go
name := "some_environment_variable_that_is_not_set"
os.Unsetenv(name)
fmt.Println(Get(name))
fmt.Println(Get(name, "this is the default"))
fmt.Println(Get(name, "this is the default", "this one is ignored"))
fmt.Println(Get(name, "", "this one is ignored"))
os.Setenv(name, "this is the value set")
fmt.Println(Get(name))
fmt.Println(Get(name, "this is the default"))
```

Output:
```

this is the default
this is the default

this is the value set
this is the value set
```
</p></details>
