package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"sync"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func defaultConfig(fns ...func(*Config)) *Config {
	cfg := &Config{
		Path:      "/dev/stdout",
		Imports:   true,
		Returns:   true,
		Local:     "github.com/demosdemon",
		Comments:  true,
		TabIndent: true,
		TabWidth:  8,
	}

	for _, fn := range fns {
		fn(cfg)
	}

	return cfg
}

func captureFile(fp **os.File) (<-chan string, func(), error) {
	prev := *fp
	cancel := func() {
		*fp = prev
	}

	r, w, err := os.Pipe()
	if err != nil {
		return nil, nil, err
	}

	ch := make(chan string)
	*fp = w
	go func() {
		var buf bytes.Buffer
		_, _ = io.Copy(&buf, r)
		ch <- string(buf.Bytes())
		close(ch)
	}()

	return ch, cancel, nil
}

func captureOutput(fn func()) (stdout string, stderr string, err error) {
	log := logrus.StandardLogger()
	prevExit := log.ExitFunc
	defer func() {
		log.ExitFunc = prevExit
	}()

	log.ExitFunc = func(i int) {
		panic(exit(i))
	}

	stdoutC, cancel, err := captureFile(&os.Stdout)
	if err != nil {
		return "", "", err
	}
	defer cancel()

	stderrC, cancel, err := captureFile(&os.Stderr)
	if err != nil {
		return "", "", err
	}
	defer cancel()

	go func() {
		defer os.Stdout.Close()
		defer os.Stderr.Close()
		defer func() {
			if r := recover(); r != nil {
				if r, ok := r.(error); ok {
					err = r
				} else {
					err = fmt.Errorf("panic: %v", r)
				}
			}
		}()

		fn()
	}()

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		stdout = <-stdoutC
		wg.Done()
	}()

	go func() {
		stderr = <-stderrC
		wg.Done()
	}()

	wg.Wait()

	return stdout, stderr, err
}

func configDevNull(cfg *Config)           { cfg.Path = "/dev/null" }
func configDevStdout(cfg *Config)         { cfg.Path = "/dev/stdout" }
func configDevStderr(cfg *Config)         { cfg.Path = "/dev/stderr" }
func configDiff(cfg *Config)              { cfg.Diff = true }
func configWrite(cfg *Config)             { cfg.Write = true }
func configExitCode(cfg *Config)          { cfg.ExitCode = true }
func configNoClobber(cfg *Config)         { cfg.NoClobber = true }
func configImports(cfg *Config)           { cfg.Imports = false }
func configReturns(cfg *Config)           { cfg.Returns = false }
func configLocal(cfg *Config)             { cfg.Local = "go.platform.sh" }
func configAllErrors(cfg *Config)         { cfg.AllErrors = true }
func configComments(cfg *Config)          { cfg.Comments = false }
func configTabIndent(cfg *Config)         { cfg.TabIndent = false }
func configTabWidth(cfg *Config)          { cfg.TabWidth = 4 }
func configFormatOnly(cfg *Config)        { cfg.FormatOnly = true }
func configPrintErrors(cfg *Config)       { cfg.PrintErrors = true }
func configRemoveBareReturns(cfg *Config) { cfg.RemoveBareReturns = true }

func TestNewConfig(t *testing.T) {
	cases := []struct {
		name         string
		args         []string
		config       *Config
		expectErr    bool
		expectStdout string
		expectStderr string
	}{
		{
			name:   "zero",
			args:   []string{},
			config: defaultConfig(),
		},
		{
			name:   "path",
			args:   []string{"-path", "/dev/null"},
			config: defaultConfig(configDevNull),
		},
		{
			name:   "diff",
			args:   []string{"-diff"},
			config: defaultConfig(configDiff),
		},
		{
			name:   "write",
			args:   []string{"-write"},
			config: defaultConfig(configWrite),
		},
		{
			name:   "exit-code",
			args:   []string{"-exit-code"},
			config: defaultConfig(configExitCode),
		},
		{
			name:   "no-clobber",
			args:   []string{"-no-clobber"},
			config: defaultConfig(configNoClobber),
		},
		{
			name:   "imports",
			args:   []string{"-imports=false"},
			config: defaultConfig(configImports),
		},
		{
			name:   "returns",
			args:   []string{"-returns=false"},
			config: defaultConfig(configReturns),
		},
		{
			name:   "local",
			args:   []string{"-local", "go.platform.sh"},
			config: defaultConfig(configLocal),
		},
		{
			name:   "all-errors",
			args:   []string{"-all-errors"},
			config: defaultConfig(configAllErrors),
		},
		{
			name:   "comments",
			args:   []string{"-comments=false"},
			config: defaultConfig(configComments),
		},
		{
			name:   "tab-indent",
			args:   []string{"-tab-indent=false"},
			config: defaultConfig(configTabIndent),
		},
		{
			name:   "tab-width",
			args:   []string{"-tab-width", "4"},
			config: defaultConfig(configTabWidth),
		},
		{
			name:   "format-only",
			args:   []string{"-format-only"},
			config: defaultConfig(configFormatOnly),
		},
		{
			name:   "print-errors",
			args:   []string{"-print-errors"},
			config: defaultConfig(configPrintErrors),
		},
		{
			name:   "remove-bare-returns",
			args:   []string{"-remove-bare-returns"},
			config: defaultConfig(configRemoveBareReturns),
		},
		{
			name:         "unknown-flag",
			args:         []string{"-unknown-flag"},
			expectErr:    true,
			expectStderr: "flag provided but not defined: -unknown-flag\nUsage of gen:\n  -all-errors\n    \treport all errors (not just the first 10 on different lines) (default false)\n  -comments\n    \tkeep comments (default true)\n  -diff\n    \tdisplay a diff instead of rewriting files (default false)\n  -exit-code\n    \texit with a failure code if no change was detected (default false)\n  -format-only\n    \tdisable the insertion and deletions of imports (default false)\n  -imports\n    \trun goimports on the file post generation (default true)\n  -local value\n    \tput imports beginning with this string after 3rd-party packages (see goimports) (default github.com/demosdemon)\n  -no-clobber\n    \tfail to write a file if it already exists (default false)\n  -path value\n    \tpath of the generated file (default /dev/stdout)\n  -print-errors\n    \tprint non-fatal typechecking errors to stderr (default false)\n  -remove-bare-returns\n    \tremove bare returns (default false)\n  -returns\n    \trun goreturns on the file post generation (default true)\n  -tab-indent\n    \tuse tabs for indent (default true)\n  -tab-width value\n    \tset tab width (default 8)\n  -write\n    \twrite the generated file (default false)\n",
		},
	}

	for _, c := range cases {
		c := c
		t.Run(c.name, func(t *testing.T) {
			var cfg *Config
			var err error
			stdout, stderr, err2 := captureOutput(func() {
				cfg, err = NewConfig(c.args)
			})
			require.NoError(t, err2)
			assert.True(t, (err != nil) == c.expectErr)
			assert.EqualValues(t, c.config, cfg)
			assert.Equal(t, c.expectStdout, stdout)
			assert.Equal(t, c.expectStderr, stderr)
		})
	}
}

func testMain(tb testing.TB) {
	prevArgs := append(([]string)(nil), os.Args...)
	defer func() {
		os.Args = prevArgs
	}()

	tmp, err := ioutil.TempFile("", "*.go")
	require.NoError(tb, err)
	require.NoError(tb, tmp.Close())

	tmpName := tmp.Name()
	require.NoError(tb, os.Remove(tmpName))
	defer func() {
		_ = os.Remove(tmpName)
	}()

	os.Args = []string{
		os.Args[0],
		"-path",
		tmpName,
		"-write",
	}

	stdout, stderr, err := captureOutput(main)
	assert.Equal(tb, "", stdout)
	assert.Equal(tb, "", stderr)

	got, err := ioutil.ReadFile(tmpName)
	require.NoError(tb, err)

	expected, err := ioutil.ReadFile("../../generated.go")
	require.NoError(tb, err)

	require.Equal(tb, expected, got)
}

func Test_main(t *testing.T) {
	testMain(t)
}

func Benchmark_main(b *testing.B) {
	for idx := 0; idx < b.N; idx++ {
		testMain(b)
	}
}

func Test_must(t *testing.T) {
	cases := []struct {
		name   string
		err    error
		panics bool
	}{
		{
			name: "nil",
		},
		{
			name:   "err",
			err:    assert.AnError,
			panics: true,
		},
	}

	for _, c := range cases {
		c := c
		t.Run(c.name, func(t *testing.T) {
			fn := func() {
				must(c.err)
			}

			if c.panics {
				assert.PanicsWithValue(t, c.err, fn)
			} else {
				assert.NotPanics(t, fn)
			}
		})
	}
}

func Benchmark_must(b *testing.B) {
	cases := []struct {
		name   string
		err    error
		panics bool
	}{
		{
			name: "nil",
		},
		{
			name:   "err",
			err:    assert.AnError,
			panics: true,
		},
	}

	for idx := 0; idx < b.N; idx++ {
		for _, c := range cases {
			c := c
			b.Run(c.name, func(b *testing.B) {
				fn := func() {
					must(c.err)
				}

				for idx := 0; idx < b.N; idx++ {
					if c.panics {
						assert.PanicsWithValue(b, c.err, fn)
					} else {
						assert.NotPanics(b, fn)
					}
				}
			})
		}
	}
}

func TestExecute(t *testing.T) {
	cases := []struct {
		name   string
		args   []string
		data   Render
		exit   int
		stdout string
		stderr string
	}{
		{
			name:   "args error",
			args:   []string{"-help"},
			data:   testRender{},
			exit:   1,
			stderr: "Usage of gen:\n  -all-errors\n    \treport all errors (not just the first 10 on different lines) (default false)\n  -comments\n    \tkeep comments (default true)\n  -diff\n    \tdisplay a diff instead of rewriting files (default false)\n  -exit-code\n    \texit with a failure code if no change was detected (default false)\n  -format-only\n    \tdisable the insertion and deletions of imports (default false)\n  -imports\n    \trun goimports on the file post generation (default true)\n  -local value\n    \tput imports beginning with this string after 3rd-party packages (see goimports) (default github.com/demosdemon)\n  -no-clobber\n    \tfail to write a file if it already exists (default false)\n  -path value\n    \tpath of the generated file (default /dev/stdout)\n  -print-errors\n    \tprint non-fatal typechecking errors to stderr (default false)\n  -remove-bare-returns\n    \tremove bare returns (default false)\n  -returns\n    \trun goreturns on the file post generation (default true)\n  -tab-indent\n    \tuse tabs for indent (default true)\n  -tab-width value\n    \tset tab width (default 8)\n  -write\n    \twrite the generated file (default false)\n  -all-errors\n    \treport all errors (not just the first 10 on different lines) (default false)\n  -comments\n    \tkeep comments (default true)\n  -diff\n    \tdisplay a diff instead of rewriting files (default false)\n  -exit-code\n    \texit with a failure code if no change was detected (default false)\n  -format-only\n    \tdisable the insertion and deletions of imports (default false)\n  -imports\n    \trun goimports on the file post generation (default true)\n  -local value\n    \tput imports beginning with this string after 3rd-party packages (see goimports) (default github.com/demosdemon)\n  -no-clobber\n    \tfail to write a file if it already exists (default false)\n  -path value\n    \tpath of the generated file (default /dev/stdout)\n  -print-errors\n    \tprint non-fatal typechecking errors to stderr (default false)\n  -remove-bare-returns\n    \tremove bare returns (default false)\n  -returns\n    \trun goreturns on the file post generation (default true)\n  -tab-indent\n    \tuse tabs for indent (default true)\n  -tab-width value\n    \tset tab width (default 8)\n  -write\n    \twrite the generated file (default false)\n",
		},
		{
			name: "render error",
			args: []string{},
			data: testRender{err: assert.AnError},
			exit: 1,
		},
		{
			name: "no error",
			args: []string{"-path", "/dev/null"},
			data: testRender{output: []byte("package foobar\n")},
			exit: 0,
		},
	}

	for _, c := range cases {
		c := c
		t.Run(c.name, func(t *testing.T) {
			stdout, stderr, err := captureOutput(func() {
				Execute(c.args, c.data)
			})

			if err == nil {
				err = exit(0)
			}

			assert.Equal(t, exit(c.exit), err)
			assert.Equal(t, c.stdout, stdout)
			assert.Equal(t, c.stderr, stderr)
		})
	}
}

func TestConfig_Execute(t *testing.T) {
	tmpData := []byte(`// this file is temporary and safe to delete

package main

import "log"

func main() {
	log.Println("Hello, World")
}
`)

	invalid := []byte("package go\n")

	populated := func() (string, error) {
		tmp, err := ioutil.TempFile("", "*.go")
		if err != nil {
			return "", err
		}

		_, err = tmp.Write(tmpData)
		if err != nil {
			_ = os.Remove(tmp.Name())
			return "", err
		}

		return tmp.Name(), tmp.Close()
	}

	empty := func() (string, error) {
		tmp, err := ioutil.TempFile("", "*.go")
		if err != nil {
			return "", err
		}

		return tmp.Name(), tmp.Close()
	}

	notExist := func() (string, error) {
		tmp, err := ioutil.TempFile("", "*.go")
		if err != nil {
			return "", err
		}

		err = tmp.Close()
		if err != nil {
			_ = os.Remove(tmp.Name())
			return "", err
		}

		return tmp.Name(), os.Remove(tmp.Name())
	}

	nothing := func() (string, error) {
		return "", nil
	}

	accessDenied := func() (string, error) {
		tmp, err := ioutil.TempFile("", "*.go")
		if err != nil {
			return "", err
		}

		err = tmp.Close()
		if err != nil {
			_ = os.Remove(tmp.Name())
			return "", err
		}

		return tmp.Name(), os.Chmod(tmp.Name(), 0200)
	}

	cases := []struct {
		name   string
		config *Config
		mktemp func() (string, error)
		render Render
		expect []byte
		stdout string
		stderr string
		err    bool
	}{
		{
			name:   "Defaults",
			config: defaultConfig(),
			mktemp: empty,
			render: testRender{output: tmpData},
			expect: []byte{},
		},
		{
			name:   "ReadPrevious/Stdout",
			config: defaultConfig(configDevStdout),
			mktemp: nothing,
			render: testRender{output: tmpData},
		},
		{
			name:   "ReadPrevious/Stderr",
			config: defaultConfig(configDevStderr),
			mktemp: nothing,
			render: testRender{output: tmpData},
		},
		{
			name:   "ReadPrevious/NotExist",
			config: defaultConfig(),
			mktemp: notExist,
			render: testRender{output: tmpData},
		},
		{
			name:   "ReadPrevious/AccessDenied",
			config: defaultConfig(),
			mktemp: accessDenied,
			render: testRender{},
			expect: []byte{},
			err:    true,
		},
		{
			name:   "Render/Error",
			config: defaultConfig(),
			mktemp: populated,
			render: testRender{err: assert.AnError},
			expect: tmpData,
			err:    true,
		},
		{
			name:   "GoImports/Disabled",
			config: defaultConfig(configImports),
			mktemp: empty,
			render: testRender{output: tmpData},
			expect: []byte{},
		},
		{
			name:   "GoImports/Error",
			config: defaultConfig(),
			mktemp: empty,
			render: testRender{output: invalid},
			expect: []byte{},
			err:    true,
		},
		{
			name:   "GoReturns/Disabled",
			config: defaultConfig(configReturns),
			mktemp: empty,
			render: testRender{output: tmpData},
			expect: []byte{},
		},
		{
			name:   "GoReturns/Error",
			config: defaultConfig(configImports),
			mktemp: empty,
			render: testRender{output: invalid},
			expect: []byte{},
			err:    true,
		},
		{
			name:   "PrintDiff/Enabled",
			config: defaultConfig(configDiff),
			mktemp: empty,
			render: testRender{output: tmpData},
			expect: []byte{},
			stdout: "--- a\n+++ b\n@@ -1 +1,10 @@\n+// this file is temporary and safe to delete\n \n+package main\n+\n+import \"log\"\n+\n+func main() {\n+\tlog.Println(\"Hello, World\")\n+}\n+\n",
		},
		{
			name:   "WriteFile/NoChange/NoExit",
			config: defaultConfig(configWrite),
			mktemp: populated,
			render: testRender{output: tmpData},
			expect: tmpData,
			err:    false,
		},
		{
			name:   "WriteFile/NoChange/Exit",
			config: defaultConfig(configWrite, configExitCode),
			mktemp: populated,
			render: testRender{output: tmpData},
			expect: tmpData,
			err:    true,
		},
		{
			name:   "WriteFile/NoClobber/Exists",
			config: defaultConfig(configWrite, configNoClobber),
			mktemp: empty,
			render: testRender{output: tmpData},
			expect: []byte{},
			err:    true,
		},
		{
			name:   "WriteFile/Enabled",
			config: defaultConfig(configWrite),
			mktemp: notExist,
			render: testRender{output: tmpData},
			expect: tmpData,
		},
	}

	for _, c := range cases {
		c := c
		t.Run(c.name, func(t *testing.T) {
			tmp, err := c.mktemp()
			require.NoError(t, err)
			if tmp != "" {
				defer os.Remove(tmp)
				c.config.Path = tmp
			}

			stdout, stderr, err2 := captureOutput(func() {
				err = c.config.Execute(c.render)
			})
			assert.NoError(t, err2)

			if tmp != "" {
				tmpData, err := ioutil.ReadFile(tmp)
				if c.expect == nil {
					assert.True(t, os.IsNotExist(err))
				} else {
					assert.Equal(t, string(c.expect), string(tmpData))
				}
			}

			assert.True(t, (err != nil) == c.err)
			assert.Equal(t, c.stdout, stdout)
			assert.Equal(t, c.stderr, stderr)
		})
	}
}

type testRender struct {
	output []byte
	err    error
}

func (t testRender) Render(w io.Writer) error {
	if t.err != nil {
		return t.err
	}
	_, err := w.Write(t.output)
	return err
}

type exit int

func (e exit) Error() string {
	return fmt.Sprintf("exit status %d", int(e))
}
