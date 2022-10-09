package gonet

import (
	"errors"
	"fmt"
	"log"
	"time"
)

type ActorModel struct {
	ctx    *context
	handle uint32
	name   string
}

func (actor *ActorModel) setContext(ctx *context) {
	actor.ctx = ctx
}

func (actor *ActorModel) setHandle(handle uint32) {
	actor.handle = handle
}

func (actor *ActorModel) GetHandle() uint32 {
	return actor.handle
}

func (actor *ActorModel) GetName() string {
	return actor.name
}

func (actor *ActorModel) RegName(name string) error {
	if len(actor.name) > 0 {
		return errors.New("actor.name > 0")
	}

	_, err := RegName(actor.handle, name)
	if err == nil {
		actor.name = name
	}
	return err
}

func (actor *ActorModel) NewService(name string, args ...interface{}) uint32 {
	return NewService(name, args...)
}

func (actor *ActorModel) CheckError(err error) {
	if err != nil {
		panic(err)
	}
}

func (actor *ActorModel) Send(handle uint32, fname string, args ...interface{}) {
	ctx := get(handle)
	if ctx == nil {
		log.Println("not found context by handle", handle, "for", fname)
		return
	}

	ctx.send(actor.ctx, fname, args...)
}

func (actor *ActorModel) SendByName(name string, fname string, args ...interface{}) error {
	handle, err := getByName(name)
	if err != nil {
		return err
	}

	actor.Send(handle, fname, args...)

	return nil
}

func (actor *ActorModel) Call(handle uint32, fname string, args ...interface{}) ([]interface{}, error) {
	ctx := get(handle)
	if ctx == nil {
		return nil, errors.New(fmt.Sprintf("Can't find Context. %s", fname))
	}

	return ctx.call(actor.ctx, fname, args...) /*, nil*/
}

func (actor *ActorModel) CallByName(name string, fname string, args ...interface{}) ([]interface{}, error) {
	handle, err := getByName(name)
	if err != nil {
		return nil, err
	}

	return actor.Call(handle, fname, args...)
}

func (actor *ActorModel) Timeout(d time.Duration, funcName string, args ...interface{}) *Timer {
	return Timeout(actor.handle, d, funcName, args...)
}

func (actor *ActorModel) SetTimeout(handle uint32, name string, timeout int64) error {
	set := func(handle uint32, timout int64) error {
		ctx := get(handle)
		if ctx == nil {
			return fmt.Errorf("Can't find Context.")
		}
		ctx.setTimeout(timeout)
		return nil
	}

	if handle != 0 {
		return set(handle, timeout)
	} else {
		handle, err := getByName(name)
		if err != nil {
			return err
		}
		return set(handle, timeout)
	}

}
