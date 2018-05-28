package logs

import (
	// "errors"
	"fmt"
	// "io"

	seelog "github.com/cihub/seelog"
)

var Logger seelog.LoggerInterface

func loadAppConfig() {
	appConfig := `
	<seelog>
		<outputs formatid="main">
			<filter levels="info,debug,critical,error">
				<console /> 
			</filter>
			<filter levels="debug">
				<file path="debug.txt" /> 
			</filter>
		</outputs>
		<formats>
			<format id="main" format="[%LEV] %Date/%Time :%n%Msg%n"/> 
		</formats>
	</seelog>
`
	logger, err := seelog.LoggerFromConfigAsBytes([]byte(appConfig))
	if err != nil {
		fmt.Println(err)
		return
	}
	UseLogger(logger)
}

func init() {
	DisableLog()
	loadAppConfig()
}

// DisableLog disables all library log output
func DisableLog() {
	Logger = seelog.Disabled
}

// UseLogger uses a specified seelog.LoggerInterface to output library log.
// Use this func if you are using Seelog logging system in your app.
func UseLogger(newLogger seelog.LoggerInterface) {
	Logger = newLogger
}
