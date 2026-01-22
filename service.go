package goutils

import (
	"fmt"
	"github.com/kardianos/service"
	"sync"
)

type program struct {
	// Use a wait group to track active goroutines
	wg      sync.WaitGroup
	countWg int
}

var programInstance *program
var ServiceInstance service.Service

func (p *program) Start(s service.Service) error {
	fmt.Println(s.String() + " Rodando!")
	return nil
}

func (p *program) Stop(s service.Service) error {
	fmt.Println(fmt.Sprintf("%s Parando com %d threads ainda rodando...", s.String(), p.countWg))
	p.DoneWaitGroup()

	p.wg.Wait()
	fmt.Println(s.String() + " Parado!")
	return nil
}

func (p *program) AddWaitGroup() {
	p.wg.Add(1)
	p.countWg++
}

func (p *program) DoneWaitGroup() {
	fmt.Println("Finishing Thread")
	p.wg.Done()
	p.countWg--
}

func InitService() {
	programInstance = &program{
		wg: sync.WaitGroup{},
	}
	programInstance.AddWaitGroup()
}

func ExecuteService(serviceConfig *service.Config) {
	if programInstance == nil {
		InitService()
	}

	var err error
	ServiceInstance, err = service.New(programInstance, serviceConfig)
	if err != nil {
		fmt.Println("Não pode criar o serviço: " + err.Error())
	}

	err = ServiceInstance.Run()
	if err != nil {
		fmt.Println("Não é possível iniciar o serviço: " + err.Error())
	}
}

func AddThreadForShutdown() {
	if programInstance == nil {
		InitService()
	}
	programInstance.AddWaitGroup()
}

func DoneThreadForShutdown() {
	programInstance.DoneWaitGroup()
}
