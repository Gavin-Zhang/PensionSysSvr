package main

import (
	"github.com/cihub/seelog"
)

func initSeelog(cfg string) {
	if cfg == "" {
		cfg = `
<seelog>
	<outputs formatid="main">
		<filter levels="debug,info,critical,error">
			<console />
		</filter>
		<filter levels="info">
			<file path="log/info.log"/>
		</filter>
		<filter levels="debug">
			<file path="log/debug.log"/>
		</filter>
		<filter levels="error">
			<file path="log/error.log"/>
		</filter>
		<rollingfile formatid="rool" type="date" filename="log/roll.log" datepattern="2006.01.02." maxrolls="30" />
	</outputs>

	<formats>
		<format id="main" format="%Date/%Time [%Level] %Msg%n"/>
		<format id="info" format="%Date/%Time [%Level] %Msg%n"/>
		<format id="debug" format="%Date/%Time [%Level] %Msg%n"/>
		<format id="error" format="%Date/%Time [%Level] %Msg%n"/>
		<format id="rool" format="%Date/%Time [%Level] %Msg%n    %RelFile:%Line(%Func)%n"/>
	</formats>
</seelog>`
	}

	logger, err := seelog.LoggerFromConfigAsBytes([]byte(cfg))
	if err != nil {
		panic(err)
	}
	seelog.ReplaceLogger(logger)
}
