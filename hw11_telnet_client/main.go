package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"regexp"
	"strconv"
	"syscall"
	"time"

	"github.com/spf13/pflag"
)

var (
	ErrWrongHost = errors.New("wrong host")
	ErrWrongPort = errors.New("wrong port")
)

type config struct {
	timeout time.Duration
	host    string
	port    string
	socket  string
}

func main() {
	config, err := initConfig(parseArg())
	if err != nil {
		log.Fatal(err)
	}

	client := NewTelnetClient(config.socket, config.timeout, os.Stdin, os.Stdout)
	if err := client.Connect(); err != nil {
		log.Fatal(err)
	}

	defer func(client TelnetClient) {
		err := client.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(client)

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGQUIT, syscall.SIGINT)

	go func() {
		defer cancel()
		if err := client.Send(); err != nil {
			_, err := fmt.Fprint(os.Stderr, err)
			if err != nil {
				log.Fatal(err)
			}
		}
	}()

	go func() {
		defer cancel()
		if err := client.Receive(); err != nil {
			_, err := fmt.Fprint(os.Stderr, err)
			if err != nil {
				log.Fatal(err)
			}
		}
	}()

	<-ctx.Done()
}

func initConfig(host, port string, timeout time.Duration) (*config, error) {
	c := config{}

	host, err := validateHost(host)
	if err != nil {
		return &c, err
	}

	port, err = validatePort(port)
	if err != nil {
		return &c, err
	}

	c.timeout = timeout
	c.host = host
	c.port = port
	c.socket = net.JoinHostPort(c.host, c.port)

	return &c, nil
}

func validateHost(host string) (string, error) {
	if host == "" {
		return "", ErrWrongHost
	}

	pattern := `^[a-zA-Z\d]{1}[a-zA-Z\d\.\-]+$`

	matched, _ := regexp.MatchString(pattern, host)
	if !matched {
		return "", ErrWrongHost
	}

	return host, nil
}

func validatePort(port string) (string, error) {
	if port == "" {
		return "", ErrWrongPort
	}

	tcpPort, err := strconv.Atoi(port)
	if err != nil {
		return "", err
	}

	if tcpPort < 1 || tcpPort > 65535 {
		return "", ErrWrongPort
	}

	return port, nil
}

func parseArg() (string, string, time.Duration) {
	var timeout time.Duration

	pflag.DurationVarP(&timeout, "timeout", "t", 10*time.Second, "")
	pflag.Parse()

	return pflag.Arg(0), pflag.Arg(1), timeout
}
