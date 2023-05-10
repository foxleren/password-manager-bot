package telegram

import (
	"fmt"
	"github.com/siruspen/logrus"
	"time"
)

type ParserData struct {
	PythonFile string
	ExcelFile  string
}

var parsingTime string
var parsingUpdateCounter string = "0"

func setParsingTime() {
	t := time.Now()
	parsingTime = fmt.Sprintf("%d-%02d-%02d %02d:%02d",
		t.Year(), t.Month(), t.Day(),
		t.Hour(), t.Minute())
	logrus.Printf("Parsing time: %s", parsingTime)
}

func (b *Bot) compileParser() error {
	//logrus.Println("Started python compilation.")
	//cmd := exec.Command("python3", b.parserData.PythonFile)
	//err := cmd.Run()
	//if err != nil {
	//	return err
	//}
	//logrus.Println("Finished python compilation.")
	////parsingUpdateCounter = string(out)
	////logrus.Printf("Result of python compilation: %s\n", out)
	//
	//setParsingTime()

	return nil
}
