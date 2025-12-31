package log

import stdlog "log"

type StdLogger struct {
	l      *stdlog.Logger
	fields []Field
}

func NewStdLogger() *StdLogger {
	return &StdLogger{l: stdlog.Default()}
}

func (s *StdLogger) Info(msg string, fields ...Field) {
	s.logWithLevel("INFO", msg, fields...)
}

func (s *StdLogger) Error(msg string, fields ...Field) {
	s.logWithLevel("ERROR", msg, fields...)
}

func (s *StdLogger) Debug(msg string, fields ...Field) {
	s.logWithLevel("DEBUG", msg, fields...)
}

func (s *StdLogger) With(fields ...Field) Logger {
	if len(fields) == 0 {
		return s
	}

	combined := make([]Field, 0, len(s.fields)+len(fields))
	combined = append(combined, s.fields...)
	combined = append(combined, fields...)

	return &StdLogger{
		l:      s.l,
		fields: combined,
	}
}

func (s *StdLogger) logWithLevel(level, msg string, fields ...Field) {
	allFields := make([]Field, 0, len(s.fields)+len(fields))
	allFields = append(allFields, s.fields...)
	allFields = append(allFields, fields...)

	s.l.Println(append([]any{level, msg}, serialize(allFields)...)...)
}

func serialize(fields []Field) []any {
	out := make([]any, 0, len(fields)*2)
	for _, f := range fields {
		out = append(out, f.Key, f.Value)
	}
	return out
}
