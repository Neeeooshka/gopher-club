package config

import (
	"errors"
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

func (o *Options) GetAccrualSystem() string {
	return o.AccrualAddress.String()
}

type ServerAddress struct {
	Host string
	Port int
}

func (s *ServerAddress) String() string {
	return s.Host + ":" + strconv.Itoa(s.Port)
}

func (s *ServerAddress) Set(flag string) error {

	ss := strings.Split(flag, ":")

	if len(ss) != 2 {
		return errors.New("invalid server argument")
	}
	sp, err := strconv.Atoi(ss[1])
	if err != nil {
		return err
	}

	s.Host = ss[0]
	s.Port = sp

	return nil
}

type AccrualSystem struct {
	URI *url.URL
}

func (a *AccrualSystem) String() string {
	return a.URI.String()
}

func (a *AccrualSystem) Set(flag string) error {

	as, err := url.Parse(flag)

	if err != nil {
		return err
	}

	a.URI = as

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
		AccrualAddress: AccrualSystem{&url.URL{Host: "localhost:8080", Scheme: "http"}},
		DB:             DB{},
	}
}
