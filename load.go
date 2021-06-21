package conf

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"time"
)

const (
	tagConf  = "conf"
	tagUsage = "usage"
)

// Load sets up command-line flags using an instance of a config struct, and
// populates flag values from environment variables.
//
// v should be a pointer to an instance of a struct. Each of the fields of that
// struct will be inspected in order to generate command-line flags. By default,
// generated flags will have the same name as the field's name. This can be
// overridden with the "conf" tag.
//
// When populating (i.e. parsing) flag values from environment variables, the
// base part of flag.CommandLine's name (default value: os.Args[0]) and the
// flag's name are combined, and then converted to SCREAMING_SNAKE_CASE, in
// order to generate a corresponding environment variable for the flag.
//
// Information about the corresponding environment variable for each flag is
// appended to each flag's usage.
//
// If the value of the "conf" tag on a field is "-", no flag is generated for
// that field.
//
// The value of the "usage" tag on a field sets the usage message for its
// generated flag.
//
// The value of a field in v sets the default value for its generated flag.
//
// Load will recursively generate flags for struct-valued fields in v.
// Sub-fields will have flags named after the "path" to that field, delimited by
// "-".
//
// Load will add flags to flag.CommandLine (the default, global flag.FlagSet)
// and will call flag.Parse before returning.
//
// Load will panic if a value in os.Args or os.Environ is invalid for its
// associated flag.
func Load(v interface{}) {
	addFlags(flag.CommandLine, "", reflect.ValueOf(v).Elem())
	addEnvInfo(flag.CommandLine)
	if err := setFromEnv(flag.CommandLine); err != nil {
		// emulate Parse behavior, which calls Usage() on parse failure
		flag.CommandLine.Usage()
		panic(err)
	}

	flag.Parse()
}

func addFlags(fs *flag.FlagSet, prefix string, v reflect.Value) {
	if v.Kind() == reflect.Struct {
		for i := 0; i < v.NumField(); i++ {
			field := v.Type().Field(i)
			usage := field.Tag.Get(tagUsage)
			conf := field.Tag.Get(tagConf)

			if conf == "-" {
				continue
			}

			name := field.Name
			if conf != "" {
				name = conf
			}

			if prefix != "" {
				name = prefix + "-" + name
			}

			switch v.Field(i).Kind() {
			case reflect.Struct:
				addFlags(fs, name, v.Field(i))
			case reflect.Bool, reflect.Int64, reflect.Float64, reflect.Int, reflect.Uint, reflect.Uint64, reflect.String:
				if !v.Field(i).CanSet() {
					continue
				}

				switch f := v.Field(i).Addr().Interface().(type) {
				case *bool:
					fs.BoolVar(f, name, *f, usage)
				case *time.Duration:
					fs.DurationVar(f, name, *f, usage)
				case *float64:
					fs.Float64Var(f, name, *f, usage)
				case *int:
					fs.IntVar(f, name, *f, usage)
				case *int64:
					fs.Int64Var(f, name, *f, usage)
				case *string:
					fs.StringVar(f, name, *f, usage)
				case *uint:
					fs.UintVar(f, name, *f, usage)
				case *uint64:
					fs.Uint64Var(f, name, *f, usage)
				}
			}
		}
	}
}

func addEnvInfo(fs *flag.FlagSet) {
	fs.VisitAll(func(f *flag.Flag) {
		if f.Usage != "" {
			f.Usage += " "
		}

		f.Usage += fmt.Sprintf("(env var %s)", envVarName(fs, f))
	})
}

func setFromEnv(fs *flag.FlagSet) error {
	var err error
	fs.VisitAll(func(f *flag.Flag) {
		if err != nil {
			return
		}

		env := envVarName(fs, f)
		if s, ok := os.LookupEnv(env); ok {
			if setErr := f.Value.Set(s); setErr != nil {
				err = fmt.Errorf("invalid value %q for env var %s: %v", s, env, setErr)
			}
		}
	})

	return err
}

func envVarName(fs *flag.FlagSet, f *flag.Flag) string {
	return envify(filepath.Base(fs.Name()) + "_" + f.Name)
}
