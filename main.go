package SyslogWriter

import (
	"log"
	"net"
	"os"
	"strings"
)

type SyslogWriter struct {
	conn        *net.UDPConn
	LogFilePath *os.File
}

func NewSyslogWriter(address string) (*SyslogWriter, error) {

	logFile, err := os.OpenFile("./log.txt", os.O_CREATE|os.O_RDWR|os.O_APPEND, 0777)
	if err != nil {
		log.Fatalln("can't open log file:", err)
	}

	addr, err := net.ResolveUDPAddr("udp", address)
	if err != nil {
		return nil, err
	}

	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		return nil, err
	}

	return &SyslogWriter{conn: conn, LogFilePath: logFile}, nil
}

func (w *SyslogWriter) Write(p string) error {
	message := p
	if !strings.HasSuffix(message, "\n") {
		message += "\n"
	}
	w.LogFilePath.WriteString(p + "\n")
	w.conn.Write([]byte(message))

	return nil
}

func (w *SyslogWriter) Close() error {
	return w.conn.Close()
}
