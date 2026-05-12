package supervisor_log

import (
	"blockEmulator/params"
	"io"
	"log"
	"os"
)

type SupervisorLog struct {
	Slog *log.Logger
}

func NewSupervisorLog() *SupervisorLog {
	writer1 := os.Stdout

	dirpath := params.LogWrite_path
	err := os.MkdirAll(dirpath, os.ModePerm)
	if err != nil {
		log.Panic(err)
	}
	writer2, err := os.OpenFile(dirpath+"/Supervisor.log", os.O_WRONLY|os.O_CREATE, 0755)
	if err != nil {
		log.Panic(err)
	}
	pl := log.New(io.MultiWriter(writer1, writer2), "Supervisor: ", log.Lshortfile|log.Ldate|log.Ltime)
	return &SupervisorLog{
		Slog: pl,
	}
}
