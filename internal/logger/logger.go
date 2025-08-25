package logger

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"sync"
)

// logger singleton
var (
	mu     sync.Mutex
	logger *log.Logger
)

func init() {
	file, err := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("Não foi possível abrir o arquivo de log: %v", err)
	}
	logger = log.New(file, "", 0) // sem prefixo, vamos formatar manual
}

// logInternal escreve a mensagem formatada
func logInternal(level string, msg string) {
	mu.Lock()
	defer mu.Unlock()

	// pegar info de arquivo/linha
	_, file, line, ok := runtime.Caller(2)
	if !ok {
		file = "???"
		line = 0
	}
	logger.Printf("[%s] %s (%s:%d)\n", level, msg, file, line)
}

// PrintLog registra uma mensagem de nível INFO
func PrintLog(msg string) {
	logInternal("INFO", msg)
}

// Warn registra uma mensagem de nível WARN
func Warn(msg string) {
	logInternal("WARN", msg)
}

// Error registra uma mensagem de nível ERROR
func Error(msg string) {
	logInternal("ERROR", msg)
}

// Printf registra mensagem formatada
func Printf(level string, format string, args ...interface{}) {
	logInternal(level, fmt.Sprintf(format, args...))
}
