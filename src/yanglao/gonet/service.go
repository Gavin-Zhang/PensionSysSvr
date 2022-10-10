package gonet

import (
	"errors"
	"fmt"
	"log"
	"os"
	"reflect"
	"yanglao/gonet/utils"
)

type IService interface {
	setContext(ctx *context)
	setHandle(handle uint32)
	GetHandle() uint32
}

type ServiceCreateFunc func(args string) interface{}

type svrInfo struct {
	className string
	typeValue reflect.Type
}

var svrInfoSet map[string]*svrInfo = make(map[string]*svrInfo)

func RegService(svr IService) {
	if svr == nil {
		panic("RegService: svr is nil.")
	}

	newType := &svrInfo{}
	newType.className = utils.TypeName(svr)
	newType.typeValue = reflect.ValueOf(svr).Elem().Type()

	if svrInfoSet[newType.className] != nil {
		panic(fmt.Sprintf("RegService %s 已存在.", newType.className))
	}

	svrInfoSet[newType.className] = newType

	log.Printf("service:%s register", newType.className)
}

func NewService(svrname string, args ...interface{}) uint32 {
	defer func() {
		if err := recover(); err != nil {
			errorInfo := fmt.Sprint(err)
			log.Printf("NewService: panic:%s", errorInfo)
			os.Exit(1)
			return
		}
	}()

	if svrInfoSet[svrname] == nil {
		return 0
	}

	ctx := &context{0, "", nil, make(chan *contextMessage, contextMessageCount), make(chan int, 1), true, 20, ""}
	ctx.svr = reflect.New(svrInfoSet[svrname].typeValue).Interface().(IService)

	handle := register(ctx)
	ctx.handle = handle
	ctx.svr.setContext(ctx)
	ctx.svr.setHandle(handle)

	go context_thread(ctx)
	ctx.send(nil, "Init", args...)

	return ctx.handle
}

func KillService(handle uint32) error {
	ctx := get(handle)
	if ctx == nil {
		return errors.New(fmt.Sprintf("Can't find Context. %d", handle))
	}

	ctx.kill()
	retire(handle)
	return nil
}

// add by zhw 20150202
func KillServiceByName(name string) error {
	handle, err := getByName(name)
	if err != nil {
		return err
	}

	return KillService(handle)
}
