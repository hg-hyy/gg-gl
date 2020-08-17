package logging

import (
	"fmt"
	"hello/setting"
	"io"
	"log"
	"os"
	"path/filepath"
	"runtime"

	"github.com/gin-gonic/gin"
)

// Level ...
type Level int

var (
	// F ...
	F *os.File
	// DefaultPrefix ...
	DefaultPrefix = ""
	// DefaultCallerDepth ...
	DefaultCallerDepth = 2

	logger     *log.Logger
	logPrefix  = ""
	levelFlags = []string{"DEBUG", "INFO", "WARN", "ERROR", "FATAL"}
)

const (
	// DEBUG ...
	DEBUG Level = iota
	// INFO ...
	INFO
	// WARNING ...
	WARNING
	// ERROR ...
	ERROR
	// FATAL ...
	FATAL
)

// Setup initialize the log instance
func Setup() {
	var err error
	filePath := setting.AppSetting.LogPath
	fileName := setting.AppSetting.LogName
	F, err = MustOpen(fileName, filePath)
	if err != nil {
		log.Fatalf("logging.Setup err: %v", err)
	}

	logger = log.New(F, DefaultPrefix, log.LstdFlags)

	// f, _ := os.Create("./logs/gin.log")
	// log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.SetOutput(F)

	gin.DefaultWriter = io.MultiWriter(F)
}

// Debug output logs at debug level
func Debug(v ...interface{}) {
	setPrefix(DEBUG)
	logger.Println(v...)
}

// Info output logs at info level
func Info(v ...interface{}) {
	setPrefix(INFO)
	logger.Println(v...)
}

// Warn output logs at warn level
func Warn(v ...interface{}) {
	setPrefix(WARNING)
	logger.Println(v...)
}

// Error output logs at error level
func Error(v ...interface{}) {
	setPrefix(ERROR)
	logger.Println(v...)
}

// Fatal output logs at fatal level
func Fatal(v ...interface{}) {
	setPrefix(FATAL)
	logger.Fatalln(v...)
}

// setPrefix set the prefix of the log output
func setPrefix(level Level) {
	_, file, line, ok := runtime.Caller(DefaultCallerDepth)
	if ok {
		logPrefix = fmt.Sprintf("[%s][%s:%d]", levelFlags[level], filepath.Base(file), line)
	} else {
		logPrefix = fmt.Sprintf("[%s]", levelFlags[level])
	}

	logger.SetPrefix(logPrefix)
}

// IsNotExistMkDir create a directory if it does not exist
func IsNotExistMkDir(src string) error {
	_, err := os.Stat(src)
	if ok := os.IsNotExist(err); ok == true {
		if err := os.MkdirAll(src, os.ModePerm); err != nil {
			return err
		}
	}

	return nil
}

// MustOpen maximize trying to open the file
func MustOpen(fileName, filePath string) (*os.File, error) {
	dir, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("os.Getwd err: %v", err)
	}

	src := dir + "\\" + filePath
	_, err1 := os.Stat(src)
	perm := os.IsPermission(err1)
	if perm == true {
		return nil, fmt.Errorf("file.CheckPermission Permission denied src: %s", src)
	}

	err = IsNotExistMkDir(src)
	if err != nil {
		return nil, fmt.Errorf("file.IsNotExistMkDir src: %s, err: %v", src, err)
	}

	f, err := os.OpenFile(src+fileName, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return nil, fmt.Errorf("Fail to OpenFile :%v", err)
	}

	return f, nil
}
