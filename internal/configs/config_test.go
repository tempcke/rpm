package configs_test

import (
	"bytes"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/joho/godotenv"
	"github.com/tempcke/rpm/internal/configs"
)

const (
	FlagVerbose  = "verbose"
	FlagLogLevel = "log-level"
	FlagInterval = "interval"
	FlagName     = "name"

	DefaultVerbose  = false
	DefaultLogLevel = "info"
	DefaultInterval = 8 * time.Second
	DefaultName     = "footastic"

	inVerbose  = true
	inLogLvl   = "debug"
	inInterval = 4 * time.Second
)

func TestConfig_flags(t *testing.T) {
	var (
		appName = randLetters(5)
		args    = []string{
			"--verbose",
			"--log-level", inLogLvl,
			"--log-level2", inLogLvl,
			"--interval", inInterval.String(),
		}

		// test to ensure we can define flags in ENV_CASE and flag-case
		// this way we do not need 2 sets of constants
		flagLogLevel2 = envCase(FlagLogLevel + "2")

		fs          = flag.NewFlagSet(appName, flag.ExitOnError)
		outVerbose  = fs.Bool(FlagVerbose, DefaultVerbose, "")
		outLogLvl   = fs.String(FlagLogLevel, DefaultLogLevel, "")
		outLogLvl2  = fs.String(flagLogLevel2, DefaultLogLevel, "")
		outInterval = fs.Duration(FlagInterval, DefaultInterval, "")
		outName     = fs.String(FlagName, DefaultName, "") // no input, just make sure default works
	)

	c := configs.New(configs.WithFlagSet(fs), configs.WithArgs(args))
	if err := c.Err; err != nil {
		t.Fatal(err)
	}

	// lookup with flag-case and ENV_CASE both
	caseTests := map[string]func(string) string{
		"envCase":  envCase,
		"flagCase": flagCase,
	}
	for name, fn := range caseTests {
		t.Run(name, func(t *testing.T) {
			assertEqual(t, strconv.FormatBool(inVerbose), c.GetString(fn(FlagVerbose)))
			assertEqual(t, inVerbose, c.GetBool(fn(FlagVerbose)))
			assertEqual(t, inInterval.String(), c.GetString(fn(FlagInterval)))
			assertEqual(t, inInterval, c.GetDuration(fn(FlagInterval)))
			assertEqual(t, inLogLvl, c.GetString(fn(FlagLogLevel)))
			assertEqual(t, inLogLvl, c.GetString(fn(flagLogLevel2)))
			assertEqual(t, DefaultName, c.GetString(fn(FlagName)))
		})
	}

	// original flagSet can still be parsed, this isn't actually important but proves
	// that we set up the flagSet and args correctly for the test
	assertEqual(t, inVerbose, *outVerbose)
	assertEqual(t, inLogLvl, *outLogLvl)
	assertEqual(t, inLogLvl, *outLogLvl2)
	assertEqual(t, inInterval, *outInterval)
	assertEqual(t, DefaultName, *outName)
}
func TestConfig_env(t *testing.T) {
	var (
		appName = randLetters(5)
		args    []string // no args
		envMap  = map[string]string{
			"VERBOSE":   strconv.FormatBool(inVerbose),
			"LOG_LEVEL": inLogLvl,
			"INTERVAL":  inInterval.String(),
		}

		fs          = flag.NewFlagSet(appName, flag.ExitOnError)
		outVerbose  = fs.Bool(FlagVerbose, DefaultVerbose, "")
		outLogLvl   = fs.String(FlagLogLevel, DefaultLogLevel, "")
		outInterval = fs.Duration(FlagInterval, DefaultInterval, "")
		outName     = fs.String(FlagName, DefaultName, "") // no input, just make sure default works
	)
	t.Setenv("NAME", "something else") // it should NOT read from os.Getenv

	c := configs.New(
		configs.WithFlagSet(fs),
		configs.WithArgs(args),
		configs.WithEnvFunc(func(key string) string { return envMap[key] }),
	)
	if err := c.Err; err != nil {
		t.Fatal(err)
	}

	// lookup with flag-case and ENV_CASE both
	caseTests := map[string]func(string) string{
		"envCase":  envCase,
		"flagCase": flagCase,
	}
	for name, fn := range caseTests {
		t.Run(name, func(t *testing.T) {
			assertEqual(t, strconv.FormatBool(inVerbose), c.GetString(fn(FlagVerbose)))
			assertEqual(t, inVerbose, c.GetBool(fn(FlagVerbose)))
			assertEqual(t, inInterval.String(), c.GetString(fn(FlagInterval)))
			assertEqual(t, inInterval, c.GetDuration(fn(FlagInterval)))
			assertEqual(t, inLogLvl, c.GetString(fn(FlagLogLevel)))
			assertEqual(t, DefaultName, c.GetString(fn(FlagName)))
		})
	}

	// original flagSet can still be parsed, this isn't actually important but proves
	// that we set up the flagSet and args correctly for the test
	assertEqual(t, inVerbose, *outVerbose)
	assertEqual(t, inLogLvl, *outLogLvl)
	assertEqual(t, inInterval, *outInterval)
	assertEqual(t, DefaultName, *outName)
}
func TestConfig_envWithPrefix(t *testing.T) {
	var (
		appName = randLetters(5)
		prefix  = strings.ToUpper(appName)
		args    []string // no args
		envMap  = map[string]string{
			prefix + "_VERBOSE":   strconv.FormatBool(inVerbose),
			prefix + "_LOG_LEVEL": inLogLvl,
			prefix + "_INTERVAL":  inInterval.String(),
		}

		fs          = flag.NewFlagSet(appName, flag.ExitOnError)
		outVerbose  = fs.Bool(FlagVerbose, DefaultVerbose, "")
		outLogLvl   = fs.String(FlagLogLevel, DefaultLogLevel, "")
		outInterval = fs.Duration(FlagInterval, DefaultInterval, "")
		outName     = fs.String(FlagName, DefaultName, "") // no input, just make sure default works
	)
	t.Setenv(prefix+"_NAME", "something else") // it should NOT read from os.Getenv

	c := configs.New(
		configs.WithFlagSet(fs),
		configs.WithArgs(args),
		configs.WithEnvPrefix(prefix),
		configs.WithEnvFunc(func(key string) string { return envMap[key] }),
	)
	if err := c.Err; err != nil {
		t.Fatal(err)
	}

	// lookup with flag-case and ENV_CASE both
	caseTests := map[string]func(string) string{
		"envCase":  envCase,
		"flagCase": flagCase,
	}
	for name, fn := range caseTests {
		t.Run(name, func(t *testing.T) {
			assertEqual(t, strconv.FormatBool(inVerbose), c.GetString(fn(FlagVerbose)))
			assertEqual(t, inVerbose, c.GetBool(fn(FlagVerbose)))
			assertEqual(t, inInterval.String(), c.GetString(fn(FlagInterval)))
			assertEqual(t, inInterval, c.GetDuration(fn(FlagInterval)))
			assertEqual(t, inLogLvl, c.GetString(fn(FlagLogLevel)))
			assertEqual(t, DefaultName, c.GetString(fn(FlagName)))
		})
	}

	// original flagSet can still be parsed, this isn't actually important but proves
	// that we set up the flagSet and args correctly for the test
	assertEqual(t, inVerbose, *outVerbose)
	assertEqual(t, inLogLvl, *outLogLvl)
	assertEqual(t, inInterval, *outInterval)
	assertEqual(t, DefaultName, *outName)
}
func TestConfig_flagAndEnv(t *testing.T) {
	// flags over envs
	var (
		appName = randLetters(5)
		args    = []string{"--foo-a", "a", "--foo-b", "b"}      // a b -
		envMap  = map[string]string{"FOO_B": "B", "FOO_C": "C"} // - B C

		// expect
		inA = "a"
		inB = "b"
		inC = "C"

		fs   = flag.NewFlagSet(appName, flag.ExitOnError)
		outA = fs.String("foo-a", "", "")
		outB = fs.String("foo-b", "", "")
		outC = fs.String("foo-c", "", "")
	)

	c := configs.New(
		configs.WithFlagSet(fs),
		configs.WithArgs(args),
		configs.WithEnvFunc(func(key string) string { return envMap[key] }),
	)
	if err := c.Err; err != nil {
		t.Fatal(err)
	}

	// lookup with flag-case and ENV_CASE both
	caseTests := map[string]func(string) string{
		"envCase":  envCase,
		"flagCase": flagCase,
	}
	for name, fn := range caseTests {
		t.Run(name, func(t *testing.T) {
			assertEqual(t, inA, c.GetString(fn("foo-a")))
			assertEqual(t, inB, c.GetString(fn("foo-b")))
			assertEqual(t, inC, c.GetString(fn("foo-c")))
		})
	}

	// original flagSet can still be parsed, this isn't actually important but proves
	// that we set up the flagSet and args correctly for the test
	assertEqual(t, inA, *outA)
	assertEqual(t, inB, *outB)
	assertEqual(t, inC, *outC)
}
func TestConfig_envFile(t *testing.T) {
	var args []string

	envFile, err := writeStringToTempFile(`
VERBOSE=true
LOG_LEVEL=debug
INTERVAL=4s
`)
	assertNoError(t, err)
	envs, err := godotenv.Read(envFile)
	assertNoError(t, err)

	fs := flag.NewFlagSet("AppName", flag.ExitOnError)
	fs.Bool(FlagVerbose, DefaultVerbose, "verbose output")
	fs.String(FlagLogLevel, DefaultLogLevel, "debug|info|warn|error")
	fs.Duration(FlagInterval, DefaultInterval, "how long between sets")
	fs.String(FlagName, DefaultName, "process name")

	c := configs.New(
		configs.WithFlagSet(fs),
		configs.WithArgs(args),
		configs.WithEnvFromMap(envs),
	)
	assertNoError(t, c.Err)
	assertEqual(t, inVerbose, c.GetBool(FlagVerbose))
	assertEqual(t, inInterval, c.GetDuration(FlagInterval))
	assertEqual(t, inLogLvl, c.GetString(FlagLogLevel))
	assertEqual(t, DefaultName, c.GetString(FlagName))
}
func TestConfig_simpleFlagInit(t *testing.T) {
	args := []string{
		"--verbose",
		"--log-level", inLogLvl,
		"--interval", inInterval.String(),
	}

	// imagine in main.go you just had a function to return a flagSet
	getFlagSet := func() *flag.FlagSet {
		fs := flag.NewFlagSet("AppName", flag.ExitOnError)
		fs.Bool(FlagVerbose, DefaultVerbose, "verbose output")
		fs.String(FlagLogLevel, DefaultLogLevel, "debug|info|warn|error")
		fs.Duration(FlagInterval, DefaultInterval, "how long between sets")
		fs.String(FlagName, DefaultName, "process name")
		return fs
	}
	c := configs.New(
		configs.WithFlagSet(getFlagSet()),
		configs.WithArgs(args), // os.Args[1:] from main()
	)
	if err := c.Err; err != nil {
		t.Fatal(err)
	}

	assertEqual(t, inVerbose, c.GetBool(FlagVerbose))
	assertEqual(t, inInterval, c.GetDuration(FlagInterval))
	assertEqual(t, inLogLvl, c.GetString(FlagLogLevel))
	assertEqual(t, DefaultName, c.GetString(FlagName))
}
func TestConfig_help(t *testing.T) {
	var (
		appName = randLetters(5)
		args    = []string{"--help"}
		buf     bytes.Buffer
	)

	// imagine in main.go you just had a function to return a flagSet
	getFlagSet := func(name string) *flag.FlagSet {
		fs := flag.NewFlagSet(name, flag.ContinueOnError)
		fs.Bool(FlagVerbose, DefaultVerbose, "verbose output")
		fs.String(FlagLogLevel, DefaultLogLevel, "debug|info|warn|error")
		fs.Duration(FlagInterval, DefaultInterval, "how long between sets")
		fs.String(FlagName, DefaultName, "process name")
		return fs
	}
	c := configs.New(
		configs.WithOutput(&buf),
		configs.WithFlagSet(getFlagSet(appName)),
		configs.WithArgs(args))
	if err := c.Err; err != nil && err != flag.ErrHelp {
		t.Fatalf("unexpected error: %v", err)
	}
	assertStringContains(t, buf.String(), FlagVerbose, FlagLogLevel, DefaultLogLevel, FlagInterval, DefaultInterval.String())
}
func TestConfig_getters(t *testing.T) {
	args := []string{
		"--bool-t",
		"--bool-f=false",
		"--string", "Xyzzy",
		"--duration", "42m",
		"--int", "42",
		"--int64", "420",
		"--float", "4.2",
	}

	// imagine in main.go you just had a function to return a flagSet
	getFlagSet := func() *flag.FlagSet {
		fs := flag.NewFlagSet("AppName", flag.ExitOnError)
		fs.Bool("bool-t", false, "")
		fs.Bool("bool-f", true, "")
		fs.String("string", "", "")
		fs.Duration("duration", 0, "")
		fs.Int("int", 0, "")
		fs.Int64("int64", 0, "")
		fs.Float64("float", 0, "")
		return fs
	}
	c := configs.New(
		configs.WithFlagSet(getFlagSet()),
		configs.WithArgs(args),
	)
	if err := c.Err; err != nil {
		t.Fatal(err)
	}

	assertEqual(t, true, c.GetBool("bool-t"))
	assertEqual(t, false, c.GetBool("bool-f"))
	assertEqual(t, "Xyzzy", c.GetString("string"))
	assertEqual(t, 42*time.Minute, c.GetDuration("duration"))
	assertEqual(t, 42, c.GetInt("int"))
	assertEqual(t, int64(420), c.GetInt64("int64"))
	assertEqual(t, 4.2, c.GetFloat64("float"))
}
func TestConfig_map(t *testing.T) {
	// no flags, no env, just build a config from a map
	m := map[string]string{
		FlagVerbose:  strconv.FormatBool(inVerbose),
		FlagLogLevel: inLogLvl,
		FlagInterval: inInterval.String(),
	}
	c := configs.New(configs.WithEnvFromMap(m))
	// lookup with flag-case and ENV_CASE both
	caseTests := map[string]func(string) string{
		"envCase":  envCase,
		"flagCase": flagCase,
	}
	for name, fn := range caseTests {
		t.Run(name, func(t *testing.T) {
			assertEqual(t, inVerbose, c.GetBool(fn(FlagVerbose)))
			assertEqual(t, inInterval, c.GetDuration(fn(FlagInterval)))
			assertEqual(t, inLogLvl, c.GetString(fn(FlagLogLevel)))
		})
	}
}
func TestConfig_globalEnv(t *testing.T) {
	// no flags, should just use os.Getenv()
	var (
		// namespace env names to ensure no conflict with other tests
		prefix = strings.ToUpper(randLetters(3))

		flagVerbose  = prefix + "_" + FlagVerbose
		flagLogLevel = prefix + "_" + FlagLogLevel
		flagInterval = prefix + "_" + FlagInterval

		// flag-case or ENV_CASE keys both work
		m = map[string]string{
			flagVerbose:            strconv.FormatBool(inVerbose),
			envCase(flagLogLevel):  inLogLvl,
			flagCase(flagInterval): inInterval.String(),
		}

		// build config before envs are set, won't matter cause using os.Getenv
		c = configs.New(configs.WithEnvPrefix(prefix))
	)
	for k, v := range m {
		t.Setenv(envCase(k), v)
	}

	// lookup with flag-case and ENV_CASE both
	caseTests := map[string]func(string) string{
		"envCase":  envCase,
		"flagCase": flagCase,
	}
	for name, fn := range caseTests {
		t.Run(name, func(t *testing.T) {
			assertEqual(t, inVerbose, c.GetBool(fn(FlagVerbose)))
			assertEqual(t, inInterval, c.GetDuration(fn(FlagInterval)))
			assertEqual(t, inLogLvl, c.GetString(fn(FlagLogLevel)))
		})
	}
}

func assertEqual(t testing.TB, expected, actual interface{}) {
	t.Helper()
	if expected != actual {
		t.Errorf("Values not equal \nWant: %v \t%T\nGot:  %v \t%T\n", expected, expected, actual, actual)
	}
}
func assertStringContains(t testing.TB, haystack string, needles ...string) {
	t.Helper()
	for _, needle := range needles {
		if !strings.Contains(haystack, needle) {
			t.Errorf("expected string %q to contain %q", haystack, needle)
		}
	}
}
func assertNoError(t testing.TB, err error) {
	t.Helper()
	if err != nil {
		t.Fatal(err)
	}
}

func randLetters(n int) string {
	s := make([]rune, n)
	for i := range s {
		s[i] = 'a' + rand.Int31n(26)
	}
	return string(s)
}
func envCase(s string) string {
	return strings.ToUpper(strings.ReplaceAll(s, "-", "_"))
}
func flagCase(s string) string {
	return strings.ToLower(strings.ReplaceAll(s, "_", "-"))
}

// writeStringToTempFile writes a string to a new temporary file and returns the path to the file.
func writeStringToTempFile(data string) (string, error) {
	// Create a new temporary file in the default directory.
	tmpFile, err := os.CreateTemp("", "temp-*.txt")
	if err != nil {
		return "", fmt.Errorf("could not create temporary file: %v", err)
	}

	// Write the string data to the file.
	if _, err := tmpFile.WriteString(data); err != nil {
		tmpFile.Close()           // Ensure the file is closed before removing it.
		os.Remove(tmpFile.Name()) // Remove the temporary file.
		return "", fmt.Errorf("could not write to temporary file: %v", err)
	}

	// Close the file.
	if err := tmpFile.Close(); err != nil {
		os.Remove(tmpFile.Name()) // Remove the temporary file.
		return "", fmt.Errorf("could not close temporary file: %v", err)
	}

	return tmpFile.Name(), nil
}
