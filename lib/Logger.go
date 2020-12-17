package lib

import (
	"github.com/withmandala/go-log"
	"os"
)

var (
	TerminalLogger	  *log.Logger
	DocuLogger		  *log.Logger
)

func init() {
	file, err := os.OpenFile("logs.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		TerminalLogger.Fatal(err)
	}

	//One logger will write on terminal with different colors
	TerminalLogger = log.New(os.Stderr).WithDebug().WithColor().WithTimestamp()

	//The other logger will write on a txt
	DocuLogger = log.New(file).WithTimestamp().WithDebug()

}