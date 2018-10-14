package nmap

import (
	"fmt"
	"github.com/globalsign/mgo"
	"net/http"
	"strings"
	"time"
)

type mongoMapper struct{}

func NewMongoMapper() *mongoMapper { return &mongoMapper{} }

func (*mongoMapper) App() string {
	return "mongod"
}

func (m *mongoMapper) KnownPorts() []int {
	// https://docs.mongodb.com/manual/reference/default-mongodb-port/
	return []int{27017, 27018, 27019}
}

func (*mongoMapper) HasApp(port int, respHeader http.Header, body string) bool {
	body = strings.ToLower(body)
	return strings.Contains(body, "mongo")
}

func (*mongoMapper) Insecure(addr string) (bool, string) {
	conn, err := mgo.DialWithTimeout(fmt.Sprintf("mongodb://%s", addr), time.Second*1)
	if err == nil {
		_, err := conn.DatabaseNames()
		conn.Close()
		if err == nil {
			return true, "missing authorization"
		}
	}
	return false, ""
}
