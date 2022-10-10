package gonet

import (
	"fmt"
	"log"

	"runtime"
	"runtime/debug"
	"time"
	"yanglao/gonet/utils"
)

type context struct {
	handle      uint32
	name        string
	svr         IService
	messageChan chan *contextMessage
	quitChan    chan int
	check       bool
	timeout     int64
	Handle_func string
}

func (ctx *context) sendMessage(message *contextMessage) {
	ctx.messageChan <- message
}

func (ctx *context) send(src *context, fname string, args ...interface{}) error {
	msg := &contextMessage{src, fname, args, nil}
	ctx.messageChan <- msg

	return nil
}

func (ctx *context) setTimeout(timeout int64) {
	ctx.timeout = timeout
}

func (ctx *context) call(src *context, fname string, args ...interface{}) ([]interface{}, error) {
	msg := &contextMessage{src, fname, args, make(chan []interface{}, 1)}
	ctx.messageChan <- msg

	// ret := <-msg.reply

	var err error = nil
	if ctx.timeout != 0 {
		select {
		case ret, _ := <-msg.reply:
			return ret, err
		case <-time.After(time.Second * time.Duration(ctx.timeout)):
			{
				err = fmt.Errorf("time out context current handle function:%s", ctx.Handle_func)
				caller := func(skip int) (name string, file string, line int, ok bool) {
					var pc uintptr
					if pc, file, line, ok = runtime.Caller(skip + 1); !ok {
						return
					}
					name = runtime.FuncForPC(pc).Name()
					return
				}
				info := make([]string, 0)
				info = append(info, err.Error())
				for skip := 0; ; skip++ {
					name, file, line, ok := caller(skip)
					if !ok {
						break
					}
					skipInfo := fmt.Sprintf("%d %s Skip=%d name=%s file=%s, line=%d", ctx.handle, fname, skip, name, file, line)
					info = append(info, skipInfo)
				}
				argsInfo := fmt.Sprintln("param:", args)
				info = append(info, argsInfo)
				utils.StackLogContext(ctx.handle, ctx.name, info)
			}
		}
	} else {
		return <-msg.reply, err
	}

	return make([]interface{}, 0), err
}

func (ctx *context) kill() error {
	ctx.quitChan <- 1
	return nil
}

func context_thread(ctx *context) {
	quit := false
	for !quit {
		select {
		case msg := <-ctx.messageChan:
			if ctx.svr == nil {
				log.Printf("ctx.svr is nil")
				// 是否需要中断?
				continue
			}

			ctx.Handle_func = msg.fname
			result, err := utils.CallMethod(ctx.svr, msg.fname, msg.args)
			if err != nil {
				srcType := ""
				if msg.src != nil {
					srcType = utils.TypeName(msg.src.svr) + " "
				}

				destType := utils.TypeName(ctx.svr)
				stack := debug.Stack()
				info := make([]string, 0)
				info = append(info, fmt.Sprintf("%scall %s::%s error.", srcType, destType, msg.fname))
				info = append(info, err.Error())
				info = append(info, string(stack))
				utils.StackLog(info)
				// 是否需要中断?
				continue
			}

			if msg.reply != nil {
				msg.reply <- result
			}
		case <-ctx.quitChan:
			quit = true
		}
	}

	//log.Printf("context_thread[%d]: destroy\n", ctx.handle)
}
