package gonet

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	_ "net/http/pprof"
	"runtime"
	"runtime/debug"
	"time"
)

func Send(handle uint32, fname string, args ...interface{}) {
	ctx := get(handle)
	if ctx == nil {
		fmt.Println("Send function can't find context", handle, fname)
		return
	}

	ctx.send(nil, fname, args...)
}

func SendByName(name string, fname string, args ...interface{}) error {
	handle, err := getByName(name)
	if err != nil {
		return err
	}

	Send(handle, fname, args...)

	return nil
}

func Call(handle uint32, fname string, args ...interface{}) ([]interface{}, error) {
	ctx := get(handle)
	if ctx == nil {
		return nil, errors.New(fmt.Sprintf("Can't find Context. %s", fname))
	}

	return ctx.call(nil, fname, args...) /*, nil*/
}

func CallByName(name string, fname string, args ...interface{}) ([]interface{}, error) {
	handle, err := getByName(name)
	if err != nil {
		return nil, err
	}

	return Call(handle, fname, args...)
}

func QueryName(name string) (uint32, error) {
	handle, err := getByName(name)
	return handle, err
}

func RegName(handle uint32, name string) (string, error) {
	ctx := get(handle)
	if ctx == nil {
		return "", errors.New(fmt.Sprintf("Can't find Context. %s", name))
	}

	hs.lock.Lock()
	hs.name[name] = ctx
	ctx.name = name
	hs.lock.Unlock()

	return name, nil
}

func Run(mainService string, maxThread int, debugEnable bool) {
	runtime.GOMAXPROCS(maxThread)

	if debugEnable {
		go func() {
			log.Println(http.ListenAndServe("localhost:6066", nil))
		}()
	}

	log.Println("Welcome gonet.")

	NewService(mainService, "")

	for {
		select {
		case <-time.After(1 * time.Minute):
			//log.Printf("Auto GC.\n")
			debug.FreeOSMemory()
			runtime.GC()
		}
	}
}

func init() {
	log.Println("gonet init.")
}
