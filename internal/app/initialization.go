package app

import (
	"fmt"
	"log/slog"
)

type InitFunc func() error

func AnnounceRun(step string, init InitFunc) error {
	slog.Info(fmt.Sprintf("%s starting", step))
	err := init()
	if err != nil {
		slog.Error(fmt.Sprintf("%s failed", step), "error", err)
		return err
	}
	slog.Info(fmt.Sprintf("%s successful", step))
	return nil
}

func Run(step string, init InitFunc) error {
	slog.Debug(fmt.Sprintf("%s starting", step))
	err := init()
	if err != nil {
		slog.Error(fmt.Sprintf("%s failed", step), "error", err)
		return err
	}
	slog.Info(fmt.Sprintf("%s successful", step))
	return nil
}

func Sequence(step string, steps ...InitFunc) error {
	return AnnounceRun(step, func() error {
		for _, init := range steps {
			if err := init(); err != nil {
				return err
			}
		}
		return nil
	})
}

func Step(step string, init InitFunc) InitFunc {
	return func() error {
		return Run(step, init)
	}
}
