package syslogwriter

import (
	"log"
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
func NewSyslogWriter(address string) (*SyslogWriter, error) {

	// Открываем файл для записи логов
	logFile, err := os.OpenFile("./log.txt", os.O_CREATE|os.O_RDWR|os.O_APPEND, 0777)
	if err != nil {
		log.Fatalln("can't open log file:", err)
	}

	// Разрешаем адрес UDP сервера
	addr, err := net.ResolveUDPAddr("udp", address)
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
