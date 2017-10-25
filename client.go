package influxdb

import (
	"errors"

	"github.com/influxdata/influxdb/client/v2"
)

const (
	dbName   = "mydb"
	username = ""
	password = ""
)

type Client struct {
	Conn client.Client
}

type ClientConfig struct {
	Name string
	Type string

	Address  string
	Username string
	Password string
}

var clients = map[string]Client{}

func NewClient(config ClientConfig) (Client, error) {
	if clientExist(config.Name) {
		return clients[config.Name], nil
	}

	var c client.Client
	var err error

	switch config.Type {
	case "http":
		c, err = client.NewHTTPClient(client.HTTPConfig{
			Addr:     config.Address,
			Username: config.Username,
			Password: config.Password,
		})
	case "udp":
		c, err = client.NewUDPClient(client.UDPConfig{
			Addr: config.Address,
		})
	default:
		return Client{}, errors.New("connection type is not valid")
	}

	if err != nil {
		return Client{}, err
	}

	clients[config.Name] = Client{Conn: c}

	return clients[config.Name], nil
}

func GetClient(name string) Client {
	if clientExist(name) {
		return clients[name]
	}
	return Client{}
}

func clientExist(name string) bool {
	if c, ok := clients[name]; ok && c.Conn != nil {
		return true
	}
	return false
}
