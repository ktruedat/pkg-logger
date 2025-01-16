package log

import (
	"testing"

	"github.com/pkg/errors"
)

func TestLogger(t *testing.T) {
	tests := []struct {
		name string
		env  string
	}{
		{"Production", productionEnv},
		{"Development", developmentEnv},
	}

	for _, test := range tests {
		t.Run(
			test.name, func(_ *testing.T) {
				log := NewLogger(test.env)

				err := outer()
				log.Debug("This is a debug message")
				log.Info("This is an info message")
				log.Warning("This is a warning message")
				log.Error("This is an error message", err)
				log.Error("This is an error message", nil)
			},
		)
	}
}

func inner() error {
	return errors.New("seems we have an error here")
}

func middle() error {
	err := inner()
	if err != nil {
		return errors.Wrap(err, "middle")
	}
	return nil
}

func outer() error {
	err := middle()
	if err != nil {
		return errors.Wrap(err, "outer")
	}
	return nil
}

func TestLoggerArgs(_ *testing.T) {
	log := NewLogger(developmentEnv)

	log.Debug("Debug with args", "key1", "value1", "key2", 123, 456)
	log.Info("Info with args", "key1", "value1", "key2", 123)
	log.Warning("Warning with args", "key1", "value1", "key2", 123)
	log.Error("Error with args", errors.New("error"), "key1", "value1", "key2", 123)
}

func TestLoggerWith(_ *testing.T) {
	log := NewLogger(developmentEnv)

	newLog := log.With("key1", "value1", "key2", 123, "key3", "value3")
	newLog.Debug("Debug with context")
	newLog.Info("Info with context")
	newLog.Warning("Warning with context")
	newLog.Error("Error with context", errors.New("error"))

	log.Info("Info without context")
}

func TestLoggerNewGroup(_ *testing.T) {
	log := NewLogger(developmentEnv)

	groupLog := log.NewGroup("testGroup")
	groupLog.Debug("Debug in group")
	groupLog.Info("Info in group")

	anotherGroupLog := groupLog.NewGroup("anotherGroup")
	anotherGroupLog.Warning("Warning in group")

	groupLog.Error("Error in group", errors.New("error"))
}
