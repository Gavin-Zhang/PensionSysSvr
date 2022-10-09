package pbpacket

import (
	"bytes"
	proto "code.google.com/p/goprotobuf/proto"
	"encoding/binary"
	"fmt"
	"gonet/utils"
	//"log"
	"reflect"
	"runtime/debug"
)

type packetNewFunc func() interface{}

type packetInfo struct {
	id        uint16
	funcName  string
	className string
	typeValue reflect.Type
}

var pbPacketSet [0xFFFF]*packetInfo

func PBRegister(id uint16, message proto.Message) {
	if pbPacketSet[id] != nil {
		panic(fmt.Sprintf("PBRegister %s:%d 已存在.", id, pbPacketSet[id].className))
	}

	if message == nil {
		panic("PBRegister: message is nil.")
	}

	newType := &packetInfo{}
	newType.typeValue = reflect.ValueOf(message).Elem().Type()
	newType.id = id
	newType.className = utils.TypeName(message)
	newType.funcName = fmt.Sprint("PB", newType.className)

	pbPacketSet[id] = newType
}

func Process(agent interface{}, buffer []byte) {
	if len(buffer) == 0 {
		return
	}

	buf := bytes.NewBuffer(buffer)
	var id uint16
	binary.Read(buf, binary.BigEndian, &id)
	pkg := buf.Next(len(buffer) - 2)

	packetInfo := pbPacketSet[id]
	if packetInfo == nil {
		return
	}

	message := reflect.New(packetInfo.typeValue).Interface()
	if message == nil {
		return
	}

	proto.Unmarshal(pkg, message.(proto.Message))

	args := []interface{}{message}
	_, err := utils.CallMethod(agent, packetInfo.funcName, args)
	if err != nil {
		agentType := utils.TypeName(agent)
		//log.Printf("ProtoBuf: [协议:%s] 调用%s::%s出错.", packetInfo.className, agentType, packetInfo.funcName)
		//log.Print(err.Error())
		//debug.PrintStack()
		stack := debug.Stack()
		info := make([]string, 0)
		info = append(info, fmt.Sprintf("ProtoBuf: [协议:%s] 调用%s::%s出错.", packetInfo.className, agentType, packetInfo.funcName))
		info = append(info, err.Error())
		info = append(info, string(stack))
		utils.StackLog(info)
		return
	}
}

func Encode(id uint16, pb proto.Message) []byte {
	pbBuffer, _ := proto.Marshal(pb)

	buf := bytes.NewBuffer(nil)
	binary.Write(buf, binary.BigEndian, id)
	buf.Write(pbBuffer)

	return buf.Bytes()
}
