package zcli

import (
	"context"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/sohaha/zlsgo/zarray"
	"github.com/sohaha/zlsgo/zfile"
	"github.com/sohaha/zlsgo/zshell"
	"github.com/sohaha/zlsgo/zutil/daemon"
)

type (
	app struct {
		run    func()
		status bool
	}
	serviceStop struct {
	}
	serviceStart struct {
	}
	serviceRestart struct {
	}
	serviceInstall struct {
	}
	serviceUnInstall struct {
	}
	serviceStatus struct {
	}
)

var (
	service    daemon.ServiceIfe
	serviceErr error
	once       sync.Once
)

var s = make(chan struct{})

func (a *app) Start(daemon.ServiceIfe) error {
	a.status = true
	go func() {
		a.run()
		s <- struct{}{}
	}()
	return nil
}

func (a *app) Stop(daemon.ServiceIfe) error {
	if !a.status {
		return nil
	}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	select {
	case <-s:
	case <-ctx.Done():
		// return errors.New("forced timeout")
	}
	return nil
}

func (*serviceStatus) Flags(_ *Subcommand) {
	CheckErr(serviceErr, true)
}

func (*serviceStatus) Run(_ []string) {
	log.Printf("%s: %s\n", service.String(), service.Status())
}

func (*serviceInstall) Flags(_ *Subcommand) {
	CheckErr(serviceErr, true)
}

func (*serviceInstall) Run(_ []string) {
	CheckErr(service.Install(), true)
	CheckErr(service.Start(), true)
}

func (*serviceUnInstall) Flags(_ *Subcommand) {
	CheckErr(serviceErr, true)
}

func (*serviceUnInstall) Run(_ []string) {
	CheckErr(service.Uninstall(), true)
}

func (*serviceStart) Flags(_ *Subcommand) {
	CheckErr(serviceErr, true)
}

func (*serviceStart) Run(_ []string) {
	CheckErr(service.Start(), true)
}

func (*serviceStop) Flags(_ *Subcommand) {
	CheckErr(serviceErr, true)
}

func (*serviceStop) Run(_ []string) {
	CheckErr(service.Stop(), true)
}

func (*serviceRestart) Flags(_ *Subcommand) {
	CheckErr(serviceErr, true)
}

func (*serviceRestart) Run(_ []string) {
	CheckErr(service.Restart(), true)
}

func LaunchServiceRun(name string, description string, fn func(), config ...*daemon.Config) error {
	_, _ = LaunchService(name, description, fn, config...)
	Parse()
	if serviceErr != nil && (serviceErr != daemon.ErrNoServiceSystemDetected && !daemon.IsPermissionError(serviceErr)) {
		return serviceErr
	}
	if service == nil {
		fn()
		return nil
	}
	return service.Run()
}

func LaunchService(name string, description string, fn func(), config ...*daemon.Config) (daemon.ServiceIfe, error) {

	once.Do(func() {
		var daemonConfig *daemon.Config
		if len(config) > 0 {
			daemonConfig = config[0]
			daemonConfig.Name = name
			daemonConfig.Description = description
		} else {
			daemonConfig = &daemon.Config{
				Name:        name,
				Description: description,
				Option:      zarray.DefData{
					// "UserService": true,
				},
			}
		}

		// The file path is redirected to the current execution file path
		_, gogccflags, _, _ := zshell.Run("go env GOGCCFLAGS")
		if !strings.Contains(
			gogccflags, zfile.RealPath(zfile.ProgramPath()+"../../../..")) {
			zfile.ProjectPath = zfile.ProgramPath()
		}
		service, serviceErr = daemon.New(&app{
			run: fn,
		}, daemonConfig)

		langs[defaultLang]["install"] = "Install service"
		langs["zh"]["install"] = "安装服务"

		langs[defaultLang]["install"] = "Uninstall service"
		langs["zh"]["uninstall"] = "卸载服务"

		langs[defaultLang]["status"] = "ServiceIfe status"
		langs["zh"]["status"] = "服务状态"

		langs[defaultLang]["start"] = "Start service"
		langs["zh"]["start"] = "启动服务"

		langs[defaultLang]["stop"] = "Stop service"
		langs["zh"]["stop"] = "停止服务"

		langs[defaultLang]["restart"] = "Restart service"
		langs["zh"]["restart"] = "重启服务"

		Add("install", getLangs("install"), &serviceInstall{})
		Add("uninstall", getLangs("uninstall"), &serviceUnInstall{})
		Add("status", getLangs("status"), &serviceStatus{})
		Add("start", getLangs("start"), &serviceStart{})
		Add("stop", getLangs("stop"), &serviceStop{})
		Add("restart", getLangs("restart"), &serviceRestart{})
	})

	return service, serviceErr
}
