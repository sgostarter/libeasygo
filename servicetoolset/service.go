package servicetoolset

import (
	"flag"
	"os"
	"path/filepath"
	"runtime"

	"github.com/kardianos/service"
	"github.com/sgostarter/i/commerr"
	"github.com/sgostarter/i/l"
)

type program struct {
	wrapper *OnceWrapper
	logger  service.Logger
}

func (p *program) Start(_ service.Service) error {
	if service.Interactive() {
		_ = p.logger.Info("Running in terminal.")
	} else {
		_ = p.logger.Info("Running under service manager.")
	}

	p.wrapper.Start(p.logger)

	return nil
}

func (p *program) Stop(_ service.Service) error {
	p.wrapper.ExitRunnerAndWait()

	return nil
}

func MainEntry(serviceName, serviceDisplayName, serviceDescription string, runner Runner, logger l.Wrapper) error {
	return MainEntryEx(serviceName, serviceDisplayName, serviceDescription, "", "",
		true, true, runner, logger)
}

// MainEntryEx .
// nolint: funlen
func MainEntryEx(serviceName, serviceDisplayName, serviceDescription, workDirectory, serviceOp string,
	changeWorkDirectory, useServiceOpFlag bool, runner Runner, logger l.Wrapper) (err error) {
	if logger == nil {
		logger = l.NewNopLoggerWrapper()
	}

	if serviceName == "" || runner == nil {
		logger.Error("invalid args")

		err = commerr.ErrInvalidArgument

		return
	}

	if serviceOp == "" && useServiceOpFlag {
		if flag.Parsed() {
			logger.Fatal("flag parsed")
		}

		svcFlag := flag.String("service", "", "Control the system service.")
		flag.Parse()

		serviceOp = *svcFlag
	}

	if changeWorkDirectory {
		if workDirectory == "" {
			workDirectory, err = filepath.Abs(filepath.Dir(os.Args[0]))
			if err != nil {
				logger.WithFields(l.ErrorField(err)).Error("filepath abs")

				return
			}
		}

		err = os.Chdir(workDirectory)
		if err != nil {
			logger.WithFields(l.ErrorField(err)).Error("ch dir failed")

			return
		}
	}

	options := make(service.KeyValue)

	switch runtime.GOOS {
	case "windows":
		options["OnFailureResetPeriod"] = 0
		options["OnFailure"] = "restart"
		options["OnFailureDelayDuration"] = 0
	}

	svcConfig := &service.Config{
		Name:             serviceName,
		DisplayName:      serviceDisplayName,
		Description:      serviceDescription,
		WorkingDirectory: workDirectory,
		Option:           options,
	}

	wrapper := NewOnceWrapper(runner)

	prg := &program{
		wrapper: wrapper,
	}

	s, err := service.New(prg, svcConfig)
	if err != nil {
		logger.WithFields(l.ErrorField(err)).Error("service new")

		return
	}

	errs := make(chan error, 5)

	serviceLogger, err := s.Logger(errs)
	if err != nil {
		logger.WithFields(l.ErrorField(err)).Error("service logger")

		return
	}

	prg.logger = serviceLogger

	go func() {
		for {
			e := <-errs
			if e != nil {
				logger.WithFields(l.ErrorField(e)).Warn("service error log")
			}
		}
	}()

	if serviceOp != "" {
		err = service.Control(s, serviceOp)
		if err != nil {
			logger.WithFields(l.ErrorField(err), l.AnyField("valid actions", service.ControlAction)).
				Error("service control failed")
		}

		return
	}

	go func() {
		wrapper.Wait()

		os.Exit(99)
	}()

	_ = s.Run()

	return
}
