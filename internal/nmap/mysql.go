package nmap

import (
	"net/http"
	"strings"
)

type mysqlMapper struct{}

func NewMysqlMapper() *mysqlMapper { return &mysqlMapper{} }

func (*mysqlMapper) App() string {
	return "mysql"
}

func (*mysqlMapper) HasApp(port int, header http.Header, body string) bool {
	body = strings.ToLower(body)
	return strings.Contains(body, "packets") && strings.Contains(body, "order")
}

func (m *mysqlMapper) KnownPorts() []int {
	// https://docs.mongodb.com/manual/reference/default-mongodb-port/
	return []int{3306}
}

func (*mysqlMapper) Insecure(addr string) (bool, string) {
	return false, ""
}
