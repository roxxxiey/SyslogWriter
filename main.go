package syslogwriter

import (
	"net"
	"os"
	"strings"
)

// SyslogWriter - структура для записи логов в файл и отправки на сервер через UDP
type SyslogWriter struct {
	conn        *net.UDPConn
	LogFilePath *os.File
}

// NewSyslogWriter - создает новый экземпляр SyslogWriter
func NewSyslogWriter(logFile *os.File) (*SyslogWriter, error) {

	// Разрешаем адрес UDP сервера
	addr, err := net.ResolveUDPAddr("udp", "localhost:514")
	if err != nil {
		return nil, err
	}

	// Устанавливаем соединение с UDP сервером
	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		return nil, err
	}

	// Возвращаем экземпляр SyslogWriter
	return &SyslogWriter{conn: conn, LogFilePath: logFile}, nil
}

// Write - пишет лог в файл и отправляет его на сервер
func (w *SyslogWriter) Write(p string) error {
	message := p
	if !strings.HasSuffix(message, "\n") {
		message += "\n"
	}

	// Записываем лог в файл
	_, err := w.LogFilePath.WriteString(p + "\n")
	if err != nil {
		return err
	}

	// Отправляем лог на сервер
	_, err = w.conn.Write([]byte(message))
	if err != nil {
		return err
	}

	return nil
}

// Close - закрывает соединение с сервером
func (w *SyslogWriter) Close() error {
	return w.conn.Close()
}
