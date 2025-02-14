package syslogwriter

import (
	"fmt"
	"net"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// SyslogWriter - структура для записи логов в файл и отправки на сервер через UDP
type SyslogWriter struct {
	conn        *net.UDPConn
	LogFilePath *os.File
}

// NewSyslogWriter - создает новый экземпляр SyslogWriter
func NewSyslogWriter(address string, logFile *os.File) (*SyslogWriter, error) {

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

func (s *SyslogWriter) Emergency(input string) {
	s.Write(8, input)
}

func (s *SyslogWriter) Alert(input string) {
	s.Write(9, input)
}

func (s *SyslogWriter) Critical(input string) {
	s.Write(10, input)
}

func (s *SyslogWriter) Error(input string) {
	s.Write(11, input)
}

func (s *SyslogWriter) Warning(input string) {
	s.Write(12, input)
}

func (s *SyslogWriter) Notice(input string) {
	s.Write(13, input)
}

func (s *SyslogWriter) Info(input string) {
	s.Write(14, input)
}

func (s *SyslogWriter) Debug(input string) {
	s.Write(15, input)
}

// Write - пишет лог в файл и отправляет его на сервер
func (s *SyslogWriter) Write(priority int, input string) error {

	host, err := os.Hostname()
	if err != nil {
		host = "unknown"
	}
	executable, err := os.Executable()
	if err != nil {
		executable = "unknown"
	}
	processName := filepath.Base(executable)
	procID := os.Getpid()

	// Формируем сообщение в формате RFC 5424
	message := formatRFC5424(priority, host, processName, procID, input)

	// Записываем лог в файл
	_, err = s.LogFilePath.WriteString(message + "\n")
	if err != nil {
		return err
	}

	// Отправляем лог на сервер
	_, err = s.conn.Write([]byte(message))
	if err != nil {
		return err
	}

	// Выводим лог в консоль
	fmt.Println(message)

	return nil
}

// Close - закрывает соединение с сервером
func (s *SyslogWriter) Close() error {
	return s.conn.Close()
}

func formatRFC5424(priority int, hostname, appName string, procID int, message string) string {
	timestamp := time.Now().Format(time.RFC3339) // ISO 8601 формат
	version := 1
	msgID := "ID1"        // Пример идентификатора сообщения
	structuredData := "-" // Структурированные данные отсутствуют

	return strings.TrimSpace(
		fmt.Sprintf(
			"<%d>%d %s %s %s %d %s [%s] %s",
			priority, version, timestamp, hostname, appName, procID, msgID, structuredData, message,
		),
	)
}
