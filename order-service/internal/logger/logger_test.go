package logger

import (
	"bytes"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestInitAndGetLogger(t *testing.T) {
	tests := []struct {
		name          string
		mode          LogMode
		wantMode      LogMode
		wantInitLog   bool
		wantInitToken []string
	}{
		{
			name:          "debug mode emits init log",
			mode:          DEBUG,
			wantMode:      DEBUG,
			wantInitLog:   true,
			wantInitToken: []string{"msg=\"Logger initialized\"", "mode=DEBUG"},
		},
		{
			name:        "error mode suppresses init info",
			mode:        ERROR,
			wantMode:    ERROR,
			wantInitLog: false,
		},
		{
			name:        "prod mode suppresses init info",
			mode:        PROD,
			wantMode:    PROD,
			wantInitLog: false,
		},
		{
			name:          "custom mode keeps value and logs at info level",
			mode:          LogMode("TRACE"),
			wantMode:      LogMode("TRACE"),
			wantInitLog:   true,
			wantInitToken: []string{"msg=\"Logger initialized\"", "mode=TRACE"},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			resetForTests()
			t.Cleanup(resetForTests)

			var first *Logger
			out := captureOutput(t, func() {
				first = Init(tt.mode)
			})

			if first == nil {
				t.Fatal("Init returned nil logger")
			}
			if got := first.GetMode(); got != tt.wantMode {
				t.Fatalf("GetMode() = %s, want %s", got, tt.wantMode)
			}

			second := GetLogger()
			if first != second {
				t.Fatal("GetLogger did not return singleton instance")
			}
			if got := second.GetMode(); got != tt.wantMode {
				t.Fatalf("GetLogger() mode = %s, want %s", got, tt.wantMode)
			}

			if tt.wantInitLog {
				for _, token := range tt.wantInitToken {
					if !strings.Contains(out, token) {
						t.Fatalf("init output %q does not contain %q", out, token)
					}
				}
			} else if strings.TrimSpace(out) != "" {
				t.Fatalf("expected no init output, got %q", out)
			}
		})
	}
}

func TestGetLoggerBootstrapsDefaultSingleton(t *testing.T) {
	resetForTests()
	t.Cleanup(resetForTests)

	var first, second *Logger
	out := captureOutput(t, func() {
		first = GetLogger()
		second = GetLogger()
	})

	if first == nil || second == nil {
		t.Fatal("GetLogger returned nil")
	}
	if first != second {
		t.Fatal("GetLogger did not return the same instance twice")
	}
	if got := first.GetMode(); got != PROD {
		t.Fatalf("default mode = %s, want %s", got, PROD)
	}
	if strings.TrimSpace(out) != "" {
		t.Fatalf("expected no stdout from default bootstrap, got %q", out)
	}
}

func TestSetMode(t *testing.T) {
	tests := []struct {
		name   string
		start  LogMode
		change LogMode
		want   LogMode
	}{
		{name: "debug to error", start: DEBUG, change: ERROR, want: ERROR},
		{name: "error to prod", start: ERROR, change: PROD, want: PROD},
		{name: "prod to debug", start: PROD, change: DEBUG, want: DEBUG},
		{name: "custom value preserved", start: LogMode("TRACE"), change: LogMode("CUSTOM"), want: LogMode("CUSTOM")},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			resetForTests()
			t.Cleanup(resetForTests)

			var lg *Logger
			_ = captureOutput(t, func() {
				lg = Init(tt.start)
			})

			lg.SetMode(tt.change)
			if got := lg.GetMode(); got != tt.want {
				t.Fatalf("GetMode() = %s, want %s", got, tt.want)
			}
		})
	}
}

func TestLoggingMethodsRespectMode(t *testing.T) {
	tests := []struct {
		name       string
		mode       LogMode
		call       func(*Logger)
		wantShown  bool
		wantTokens []string
	}{
		{
			name: "debug info is shown",
			mode: DEBUG,
			call: func(l *Logger) {
				l.Info("info message", "request_id", "req-1", "attempt", 3)
			},
			wantShown:  true,
			wantTokens: []string{"msg=\"info message\"", "request_id=req-1", "attempt=3", "level=INFO"},
		},
		{
			name: "debug debug is shown",
			mode: DEBUG,
			call: func(l *Logger) {
				l.Debug("debug message", "step", "parse")
			},
			wantShown:  true,
			wantTokens: []string{"msg=\"debug message\"", "step=parse", "level=DEBUG"},
		},
		{
			name: "debug warn is shown",
			mode: DEBUG,
			call: func(l *Logger) {
				l.Warn("warn message", "code", 7)
			},
			wantShown:  true,
			wantTokens: []string{"msg=\"warn message\"", "code=7", "level=WARN"},
		},
		{
			name: "debug error is shown",
			mode: DEBUG,
			call: func(l *Logger) {
				l.Error("error message", "code", "E1")
			},
			wantShown:  true,
			wantTokens: []string{"msg=\"error message\"", "code=E1", "level=ERROR"},
		},
		{
			name: "error mode hides info",
			mode: ERROR,
			call: func(l *Logger) {
				l.Info("hidden info", "key", "value")
			},
			wantShown: false,
		},
		{
			name: "error mode hides debug",
			mode: ERROR,
			call: func(l *Logger) {
				l.Debug("hidden debug", "key", "value")
			},
			wantShown: false,
		},
		{
			name: "error mode hides warn",
			mode: ERROR,
			call: func(l *Logger) {
				l.Warn("hidden warn", "key", "value")
			},
			wantShown: false,
		},
		{
			name: "error mode shows error",
			mode: ERROR,
			call: func(l *Logger) {
				l.Error("visible error", "key", "value")
			},
			wantShown:  true,
			wantTokens: []string{"msg=\"visible error\"", "key=value", "level=ERROR"},
		},
		{
			name: "prod mode hides info",
			mode: PROD,
			call: func(l *Logger) {
				l.Info("prod hidden info")
			},
			wantShown: false,
		},
		{
			name: "prod mode hides debug",
			mode: PROD,
			call: func(l *Logger) {
				l.Debug("prod hidden debug")
			},
			wantShown: false,
		},
		{
			name: "prod mode hides warn",
			mode: PROD,
			call: func(l *Logger) {
				l.Warn("prod hidden warn")
			},
			wantShown: false,
		},
		{
			name: "prod mode shows error",
			mode: PROD,
			call: func(l *Logger) {
				l.Error("prod visible error", "component", "repo")
			},
			wantShown:  true,
			wantTokens: []string{"msg=\"prod visible error\"", "component=repo", "level=ERROR"},
		},
		{
			name: "custom mode hides debug",
			mode: LogMode("TRACE"),
			call: func(l *Logger) {
				l.Debug("custom hidden debug")
			},
			wantShown: false,
		},
		{
			name: "custom mode shows info",
			mode: LogMode("TRACE"),
			call: func(l *Logger) {
				l.Info("custom info", "feature", "slog")
			},
			wantShown:  true,
			wantTokens: []string{"msg=\"custom info\"", "feature=slog", "level=INFO"},
		},
		{
			name: "custom mode shows error",
			mode: LogMode("TRACE"),
			call: func(l *Logger) {
				l.Error("custom error", "feature", "slog")
			},
			wantShown:  true,
			wantTokens: []string{"msg=\"custom error\"", "feature=slog", "level=ERROR"},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			resetForTests()
			t.Cleanup(resetForTests)

			var lg *Logger
			out := captureOutput(t, func() {
				lg = Init(tt.mode)
				tt.call(lg)
			})

			if tt.wantShown {
				for _, token := range tt.wantTokens {
					if !strings.Contains(out, token) {
						t.Fatalf("output %q does not contain %q", out, token)
					}
				}
			}
		})
	}
}

func TestClose(t *testing.T) {
	tests := []struct {
		name       string
		mode       LogMode
		wantOutput bool
	}{
		{name: "debug emits close line", mode: DEBUG, wantOutput: true},
		{name: "error suppresses close line", mode: ERROR, wantOutput: false},
		{name: "prod suppresses close line", mode: PROD, wantOutput: false},
		{name: "custom emits close line", mode: LogMode("TRACE"), wantOutput: true},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			resetForTests()
			t.Cleanup(resetForTests)

			var lg *Logger
			out := captureOutput(t, func() {
				lg = Init(tt.mode)
				lg.Close()
			})

			if tt.wantOutput {
				if !strings.Contains(out, "Logger closed") {
					t.Fatalf("close output %q does not contain close message", out)
				}
				if !strings.Contains(out, "Logger initialized") {
					t.Fatalf("close output %q does not contain init message", out)
				}
			} else if strings.TrimSpace(out) != "" {
				t.Fatalf("expected no close output, got %q", out)
			}
		})
	}
}

func TestArgsToAttrs(t *testing.T) {
	type sample struct {
		Name string
	}

	tests := []struct {
		name    string
		args    []interface{}
		wantLen int
		wantKV  map[string]string
	}{
		{
			name:    "empty args",
			args:    nil,
			wantLen: 0,
		},
		{
			name:    "single key value pair",
			args:    []interface{}{"order_id", "123"},
			wantLen: 1,
			wantKV:  map[string]string{"order_id": "123"},
		},
		{
			name:    "multiple pairs",
			args:    []interface{}{"order_id", "123", "status", "created"},
			wantLen: 2,
			wantKV:  map[string]string{"order_id": "123", "status": "created"},
		},
		{
			name:    "odd args drop tail",
			args:    []interface{}{"order_id", "123", "orphan"},
			wantLen: 1,
			wantKV:  map[string]string{"order_id": "123"},
		},
		{
			name:    "non string key is stringified",
			args:    []interface{}{42, "answer"},
			wantLen: 1,
			wantKV:  map[string]string{"42": "answer"},
		},
		{
			name:    "nil value preserved",
			args:    []interface{}{"payload", nil},
			wantLen: 1,
			wantKV:  map[string]string{"payload": "<nil>"},
		},
		{
			name:    "struct value is kept as any",
			args:    []interface{}{"sample", sample{Name: "alpha"}},
			wantLen: 1,
			wantKV:  map[string]string{"sample": "{alpha}"},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			lg := &Logger{}
			attrs := lg.argsToAttrs(tt.args)

			if len(attrs) != tt.wantLen {
				t.Fatalf("len(attrs) = %d, want %d", len(attrs), tt.wantLen)
			}

			for key, want := range tt.wantKV {
				got, ok := attrValueAsString(attrs, key)
				if !ok {
					t.Fatalf("missing attr %q in %#v", key, attrs)
				}
				if got != want {
					t.Fatalf("attr %q = %q, want %q", key, got, want)
				}
			}
		})
	}
}

func TestPanicPathsWriteFailureLog(t *testing.T) {
	tests := []struct {
		name            string
		invoke          func(*Logger)
		wantPanicSubstr string
		wantFileTokens  []string
	}{
		{
			name: "must with error writes failure log",
			invoke: func(l *Logger) {
				l.Must(fmt.Errorf("boom"), "load config")
			},
			wantPanicSubstr: "load config: boom",
			wantFileTokens:  []string{"Context:", "load config", "Reason:", "error=boom"},
		},
		{
			name: "must not nil writes failure log",
			invoke: func(l *Logger) {
				l.MustNotNil(nil, "db connection")
			},
			wantPanicSubstr: "Value is nil: db connection",
			wantFileTokens:  []string{"Context:", "db connection", "Reason:", "context=db connection"},
		},
		{
			name: "fatal writes failure log",
			invoke: func(l *Logger) {
				l.Fatal("catastrophic", "component", "api")
			},
			wantPanicSubstr: "catastrophic",
			wantFileTokens:  []string{"Context:", "catastrophic", "Reason:", "component=api"},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			resetForTests()
			t.Cleanup(resetForTests)

			withTempWorkdir(t)

			out, panicVal := captureOutputAndRecover(t, func() {
				lg := Init(DEBUG)
				tt.invoke(lg)
			})

			if panicVal == nil {
				t.Fatal("expected panic, got nil")
			}
			if !strings.Contains(fmt.Sprint(panicVal), tt.wantPanicSubstr) {
				t.Fatalf("panic %q does not contain %q", panicVal, tt.wantPanicSubstr)
			}
			if !strings.Contains(out, "Logger initialized") {
				t.Fatalf("expected init output before panic, got %q", out)
			}
			if strings.TrimSpace(out) == "" {
				t.Fatal("expected stdout output before panic")
			}

			content, err := os.ReadFile(filepath.Join(".", "prod_failed_logs.md"))
			if err != nil {
				t.Fatalf("failed to read panic log: %v", err)
			}

			logContent := string(content)
			for _, token := range tt.wantFileTokens {
				if !strings.Contains(logContent, token) {
					t.Fatalf("panic log %q does not contain %q", logContent, token)
				}
			}
		})
	}
}

func TestPanicPathsDoNotWriteFailureLogOnSuccess(t *testing.T) {
	tests := []struct {
		name   string
		invoke func(*Logger)
	}{
		{
			name: "must with nil error",
			invoke: func(l *Logger) {
				l.Must(nil, "noop")
			},
		},
		{
			name: "must not nil with value",
			invoke: func(l *Logger) {
				l.MustNotNil("value", "noop")
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			resetForTests()
			t.Cleanup(resetForTests)

			withTempWorkdir(t)

			out := captureOutput(t, func() {
				lg := Init(DEBUG)
				tt.invoke(lg)
			})

			if strings.Contains(out, "prod_failed_logs.md") {
				t.Fatalf("unexpected failure log reference in stdout: %q", out)
			}

			if _, err := os.Stat("prod_failed_logs.md"); !os.IsNotExist(err) {
				t.Fatalf("prod_failed_logs.md should not exist on success, stat err=%v", err)
			}
		})
	}
}

func captureOutput(t *testing.T, fn func()) string {
	t.Helper()

	out, panicVal := captureOutputAndRecover(t, fn)
	if panicVal != nil {
		t.Fatalf("unexpected panic: %v", panicVal)
	}

	return out
}

func captureOutputAndRecover(t *testing.T, fn func()) (out string, panicVal any) {
	t.Helper()

	oldStdout := os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("failed to create pipe: %v", err)
	}

	os.Stdout = w
	defer func() {
		_ = w.Close()
		os.Stdout = oldStdout

		var buf bytes.Buffer
		_, _ = io.Copy(&buf, r)
		_ = r.Close()
		out = buf.String()

		if rec := recover(); rec != nil {
			panicVal = rec
		}
	}()

	fn()
	return
}

func withTempWorkdir(t *testing.T) string {
	t.Helper()

	oldWD, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get working dir: %v", err)
	}

	dir := t.TempDir()
	if err := os.Chdir(dir); err != nil {
		t.Fatalf("failed to chdir to temp dir: %v", err)
	}

	t.Cleanup(func() {
		_ = os.Chdir(oldWD)
	})

	return dir
}

func attrValueAsString(attrs []slog.Attr, key string) (string, bool) {
	for _, attr := range attrs {
		if attr.Key == key {
			return fmt.Sprint(attr.Value.Any()), true
		}
	}
	return "", false
}
