package logger

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/buffer"
	"go.uber.org/zap/zapcore"
)

// NewLogger возвращает human-friendly sugared logger с окрашенными сообщениями.
// Сохраняет вашу оригинальную логику: zap.NewDevelopmentConfig(), ISO8601 time, TimeKey="time",
// CapitalColorLevelEncoder и Development=true. Единственное добавление — обёртка-энкодер,
// которая раскрашивает Entry.Message по уровню.
func NewLogger() (*zap.SugaredLogger, error) {
	// ваш оригинальный конфиг (малые правки оставлены как есть)
	cfg := zap.NewDevelopmentConfig()
	cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	cfg.EncoderConfig.TimeKey = "time"
	cfg.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	cfg.Development = true

	// Создаём базовый console encoder согласно вашей EncoderConfig
	baseEnc := zapcore.NewConsoleEncoder(cfg.EncoderConfig)

	// Оборачиваем энкодер, чтобы раскрашивать только Entry.Message
	wrapped := &colorEncoder{
		Encoder: baseEnc,
		colorFn: colorizeMsg, // функция, задающая цвет по уровню
	}

	// Создаём core: wrapped encoder -> stdout -> уровень из cfg
	core := zapcore.NewCore(wrapped, zapcore.Lock(os.Stdout), cfg.Level)

	// Строим logger с теми же опциями, что и вы использовали ранее
	logger := zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))

	// Возвращаем sugared logger (как у вас было)
	return logger.Sugar(), nil
}

/* --------------------------- color wrapper --------------------------- */

// colorizeMsg возвращает ту же строку, обёрнутую в ANSI escape-коды по уровню.
// Лёгкая цветовая схема: ERROR -> red, WARN -> yellow, INFO -> green, DEBUG -> blue.
func colorizeMsg(level zapcore.Level, msg string) string {
	const (
		red    = "\x1b[31m"
		yellow = "\x1b[33m"
		green  = "\x1b[32m"
		blue   = "\x1b[34m"
		reset  = "\x1b[0m"
	)
	switch {
	case level >= zapcore.ErrorLevel:
		return red + msg + reset
	case level == zapcore.WarnLevel:
		return yellow + msg + reset
	case level == zapcore.InfoLevel:
		return green + msg + reset
	case level == zapcore.DebugLevel:
		return blue + msg + reset
	default:
		return msg
	}
}

// colorEncoder wraps zapcore.Encoder и заменяет Entry.Message на цветной вариант перед кодированием.
type colorEncoder struct {
	zapcore.Encoder
	colorFn func(zapcore.Level, string) string
}

// Clone required to satisfy zapcore.Encoder contract.
func (c *colorEncoder) Clone() zapcore.Encoder {
	return &colorEncoder{
		Encoder: c.Encoder.Clone(),
		colorFn: c.colorFn,
	}
}

// EncodeEntry заменяет сообщение в записи на раскрашенное, затем делегирует реальному Encoder.
func (c *colorEncoder) EncodeEntry(ent zapcore.Entry, fields []zapcore.Field) (*buffer.Buffer, error) {
	// Создаём копию Entry чтобы не мутировать оригинал, и заменяем Message
	e := ent
	if c.colorFn != nil {
		e.Message = c.colorFn(ent.Level, ent.Message)
	}
	return c.Encoder.EncodeEntry(e, fields)
}
