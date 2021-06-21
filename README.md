# conf

[![Go Reference](https://pkg.go.dev/badge/github.com/ucarion/conf.svg)](https://pkg.go.dev/github.com/ucarion/conf)

```go
package main

import (
	"fmt"

	"github.com/ucarion/conf"
)

// If you compile this into an executable named ./example, these are equivalent:
//
//  ./example -Password=yyy
//  EXAMPLE_PASSWORD=yyy ./example
//
// And in both cases, you get the following output:
//
//  config with sensitive fields redacted {admin }
//  password yyy
func main() {
	cfg := struct {
		Username string
		Password string `redact:"true"`
	}{
		Username: "admin",
		Password: "letmein",
	}

	conf.Load(&cfg)
	fmt.Println("config with sensitive fields redacted", conf.Redact(cfg))
	fmt.Println("password", cfg.Password)
}
```

`conf` is a minimal (~200 lines of code) wrapper around the standard library's
`flag` package that provides three things:

* **A less-verbose way to set up flags.** The `flag` package in the standard
  library is popular, convenient, and simple. But when you're writing programs
  that have dozens of config variables, the `flag` API is tedious and
  error-prone.

  `conf` provides a more convenient interface for creating flags. Instead of
  doing dozens of `flag.String` calls, you pass `conf.Load` a struct, and each
  field of that struct will be converted into a flag. The result is brief but
  nonetheless type-safe code.

* **A way to read flags from environment variables.** The `flag` package only
  supports reading from a list of arguments (such as `os.Args`), but in many
  environments (such as AWS Lambda) it's far more convenient to read parameters
  from env vars.

  `conf.Load` reads flag values from both env vars and `os.Args`,
  so `./my-tool -foo-bar=3` and `MY_TOOL_FOO_BAR=3 ./my-tool` do the same thing.

* **A way to log your config, but with sensitive fields redacted.** It's common
  for programs to log their configuration on startup, but doing so requires you
  to first clear out sensitive fields (such as database passwords, secret keys,
  etc.) before passing your config to your logging system.

  `conf.Redact` does that sensitive-field-clearing for you. `conf.Redact`
  returns a deep copy of your config struct, but with any fields marked as
  `redact:"true"` reset to their zero values.

## Installation

You can start using `conf` by running:

```bash
go get github.com/ucarion/conf
```
