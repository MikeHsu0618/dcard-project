package Log

import (
	"io"
	"log"
	"os"
)

var (
	Info    *log.Logger
	Warning *log.Logger
	Error   *log.Logger
)

func init() {
	errFile, err := os.OpenFile("storage/log/errors.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalln("open log file error：", err)
	}

	logFile, err := os.OpenFile("storage/log/logs.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalln("open log file error：", err)
	}

	Info = log.New(io.MultiWriter(os.Stdout, logFile), "Info:", log.Ldate|log.Ltime|log.Lshortfile)
	Warning = log.New(io.MultiWriter(os.Stdout, logFile), "Warning:", log.Ldate|log.Ltime|log.Lshortfile)
	Error = log.New(io.MultiWriter(os.Stderr, errFile), "Error:", log.Ldate|log.Ltime|log.Lshortfile)
}

//func main() {
//	Info.Println("使用示範:", "demo")
//	Warning.Printf("測試範例：%s\n", "我是範例")
//	Error.Println("範例")
//}
