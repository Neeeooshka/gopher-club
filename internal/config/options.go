package config

import (
	"errors"
	"fmt"
	"net/url"
	"strconv"
	"strings"
)

type Options struct {
	ServerAddress  ServerAddress
	AccrualAddress AccrualSystem
	DB             DB
}

func (o *Options) GetServer() string {
	return o.ServerAddress.String()
}

func (o *Options) AccrualSystem() string {
	return o.AccrualAddress.String()
}

type ServerAddress struct {
	Host string
	Port int
}

func (s *ServerAddress) String() string {
	return fmt.Sprintf("%s:%d", s.Host, s.Port)
}

func (s *ServerAddress) Set(flag string) error {

	builder, err := NewServerAddressBuilder().FromString(flag)
	if err != nil {
		return err
	}

	*s = builder.Build()

	return nil
}

type ServerAddressBuilder struct {
	serverAddress ServerAddress
}

func NewServerAddressBuilder() *ServerAddressBuilder {
	return &ServerAddressBuilder{}
}

func (b *ServerAddressBuilder) WithHost(host string) *ServerAddressBuilder {
	b.serverAddress.Host = host
	return b
}

func (b *ServerAddressBuilder) WithPort(port int) *ServerAddressBuilder {
	b.serverAddress.Port = port
	return b
}

func (b *ServerAddressBuilder) FromString(addr string) (*ServerAddressBuilder, error) {

	addressParts := strings.Split(addr, ":")

	if len(addressParts) != 2 {
		return nil, errors.New("invalid server argument, expected host:port")
	}
	port, err := strconv.Atoi(addressParts[1])
	if err != nil {
		return nil, fmt.Errorf("invalid port: %v", err)
	}

	b.serverAddress.Host = addressParts[0]
	b.serverAddress.Port = port

	return b, nil
}

func (b *ServerAddressBuilder) Build() ServerAddress {
	return b.serverAddress
}

type AccrualSystem struct {
	URI url.URL
}

func (a *AccrualSystem) String() string {
	return a.URI.String()
}

func (a *AccrualSystem) Set(flag string) error {

	as, err := url.Parse(flag)

	if err != nil {
		return err
	}

	a.URI = *as

	return nil
}

type DB struct {
	db string
}

func (d *DB) String() string {
	return d.db
}

func (d *DB) Set(flag string) error {
	d.db = flag
	return nil
}

func NewOptions() Options {
	return Options{
		ServerAddress:  ServerAddress{Host: "localhost", Port: 8080},
		AccrualAddress: AccrualSystem{url.URL{Host: "localhost:8080", Scheme: "http"}},
		DB:             DB{"user=postgres password=postgres dbname=praktikum host=PostgreSQL-17 port=5432 sslmode=disable"},
		//DB: DB{},
	}
}
