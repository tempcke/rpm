package configs

import (
	"flag"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"time"
)

type (
	Config struct {
		Err           error
		name          string
		errorHandling flag.ErrorHandling // default flag.ContinueOnError
		flags         map[string]flag.Flag
		args          []string
		envFunc       func(string) string
		vMap          map[string]flag.Value
		envPrefix     string
		output        io.Writer // default os.Stderr
	}
	Option func(*Config)

	FlagSet interface {
		VisitAll(func(*flag.Flag))
		Name() string
		ErrorHandling() flag.ErrorHandling
	}
	Value = interface {
		flag.Value
		IsSet() bool
	}

	stringValue struct {
		fv  flag.Value
		v   *string
		def string
	}
	boolValue struct {
		*stringValue
	}
)

var _ flag.Value = (*stringValue)(nil)

func New(options ...Option) Config {
	var c = &Config{
		vMap: make(map[string]flag.Value),
	}
	for _, option := range options {
		option(c)
	}
	return parse(c)
}
func (c Config) Name() string { return c.name }
func (c Config) GetString(key string) string {
	key = flagCase(key)
	if c.vMap == nil || c.vMap[key] == nil {
		return getFromEnv(c.envFunc, envName(c.envPrefix, key), "")
	}
	return c.vMap[key].String()
}
func (c Config) GetBool(key string) bool { return c.GetString(key) == "true" }
func (c Config) GetDuration(key string) time.Duration {
	// TODO: consider logging error if user provided a logger
	d, _ := time.ParseDuration(c.GetString(key))
	return d
}
func (c Config) GetInt(key string) int {
	s := c.GetString(key)
	if s == "" {
		return 0
	}
	i, _ := strconv.Atoi(s) // TODO: consider logging error if user provided a logger
	return i
}
func (c Config) GetInt64(key string) int64 { return int64(c.GetInt(key)) }
func (c Config) GetFloat64(key string) float64 {
	s := c.GetString(key)
	if s == "" {
		return 0
	}
	f, _ := strconv.ParseFloat(s, 64) // TODO: consider logging error if user provided a logger
	return f
}
func (c Config) newValue(key string, fv flag.Value) flag.Value {
	sv := &stringValue{
		fv:  fv,
		def: getFromEnv(c.envFunc, envName(c.envPrefix, key), fv.String()),
	}
	if isBoolValue(fv) {
		return &boolValue{stringValue: sv}
	}
	return sv
}
func (sv *stringValue) String() string {
	if sv == nil {
		return ""
	}
	if sv.v != nil {
		return *sv.v
	}
	return sv.def
}
func (sv *stringValue) Set(v string) error {
	if err := sv.fv.Set(v); err != nil {
		return err
	}
	v = sv.fv.String()
	sv.v = &v
	return nil
}
func (sv *stringValue) IsSet() bool {
	return sv.v != nil
}
func (bv *boolValue) IsBoolFlag() bool { return true }

func WithFlagSet(fs FlagSet) Option {
	return func(c *Config) {
		if fs == nil {
			c.Err = fmt.Errorf("flagset is nil")
			return
		}
		var flags = make(map[string]flag.Flag)
		fs.VisitAll(func(f *flag.Flag) {
			f2 := *f
			f2.Name = flagCase(f.Name)
			flags[f2.Name] = f2
		})
		c.name = fs.Name()
		c.errorHandling = fs.ErrorHandling()
		c.flags = flags
	}
}
func WithArgs(args []string) Option {
	return func(c *Config) {
		c.args = args
	}
}
func WithEnvFunc(envFunc func(string) string) Option {
	// this should only ever be used for testing when you don't want to use os.Setenv or t.Setenv
	return func(c *Config) { c.envFunc = envFunc }
}
func WithEnvFromMap(envMap map[string]string) Option {
	// hint: read .env to map with godotenv.Read(envFile)
	m := make(map[string]string, len(envMap))
	for k, v := range envMap {
		m[envCase(k)] = v
	}
	return WithEnvFunc(func(k string) string { return m[k] })
}
func WithEnvPrefix(prefix string) Option { return func(c *Config) { c.envPrefix = prefix } }
func WithOutput(output io.Writer) Option { return func(c *Config) { c.output = output } }

func parse(c *Config) Config {
	fs := flag.NewFlagSet(c.Name(), c.errorHandling)
	if c.output != nil {
		fs.SetOutput(c.output)
	}
	c.vMap = make(map[string]flag.Value)
	for _, f := range c.flags {
		val := c.newValue(f.Name, f.Value)
		c.vMap[f.Name] = val
		fs.Var(val, f.Name, f.Usage)
	}
	c.Err = fs.Parse(c.args)

	// fs.Parse does no call Set on any values that there are no args for
	// so the default values will come from the FlagSet and not the env if it exists
	// as a result the methods on this struct will return correct values
	// however the vars defined by the user while creating the flagset will not
	// to fix that we need to call Set on all values not set by an argument
	for _, v := range c.vMap {
		if sv, ok := v.(Value); ok && !sv.IsSet() {
			if err := sv.Set(sv.String()); err != nil {
				c.Err = err
				return *c
			}
		}
	}

	return *c
}
func isBoolValue(v flag.Value) bool {
	if bf, ok := v.(interface{ IsBoolFlag() bool }); ok {
		return bf.IsBoolFlag()
	}
	return false
}
func getFromEnv(fn func(string) string, key, def string) string {
	key = envCase(key)
	if fn == nil {
		if v, ok := os.LookupEnv(key); ok {
			return v
		}
		return def
	}
	if v := fn(key); v != "" {
		return v
	}
	return def
}
func envName(envPrefix, key string) string {
	key = envCase(key)
	if len(envPrefix) == 0 {
		return key
	}
	p := envCase(envPrefix) + "_"
	if strings.HasPrefix(key, p) {
		return key
	}
	return p + key
}
func flagCase(s string) string {
	key := strings.ToLower(strings.ReplaceAll(s, "_", "-"))
	return strings.Trim(key, "-")
}
func envCase(s string) string {
	key := strings.ToUpper(strings.ReplaceAll(s, "-", "_"))
	return strings.Trim(key, "_")
}
