package main

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

const timeoutDefault = 10 * time.Second

func TestTelnetClient(t *testing.T) {
	t.Run("basic", func(t *testing.T) {
		l, err := net.Listen("tcp", "127.0.0.1:")
		require.NoError(t, err)
		defer func() { require.NoError(t, l.Close()) }()

		var wg sync.WaitGroup
		wg.Add(2)

		go func() {
			defer wg.Done()

			in := &bytes.Buffer{}
			out := &bytes.Buffer{}

			timeout, err := time.ParseDuration("10s")
			require.NoError(t, err)

			client := NewTelnetClient(l.Addr().String(), timeout, io.NopCloser(in), out)
			require.NoError(t, client.Connect())
			defer func() { require.NoError(t, client.Close()) }()

			in.WriteString("hello\n")
			err = client.Send()
			require.NoError(t, err)

			err = client.Receive()
			require.NoError(t, err)
			require.Equal(t, "world\n", out.String())
		}()

		go func() {
			defer wg.Done()

			conn, err := l.Accept()
			require.NoError(t, err)
			require.NotNil(t, conn)
			defer func() { require.NoError(t, conn.Close()) }()

			request := make([]byte, 1024)
			n, err := conn.Read(request)
			require.NoError(t, err)
			require.Equal(t, "hello\n", string(request)[:n])

			n, err = conn.Write([]byte("world\n"))
			require.NoError(t, err)
			require.NotEqual(t, 0, n)
		}()

		wg.Wait()
	})

	t.Run("not connected", func(t *testing.T) {
		in := &bytes.Buffer{}
		out := &bytes.Buffer{}

		timeout, err := time.ParseDuration("3s")
		require.NoError(t, err)

		client := NewTelnetClient("", timeout, io.NopCloser(in), out)
		require.Error(t, client.Connect())
		require.Error(t, client.Send(), ErrTelnetClientNotConnected)
		require.Error(t, client.Receive(), ErrTelnetClientNotConnected)
		require.NoError(t, client.Close())
	})

	t.Run("timeout", func(t *testing.T) {
		in := &bytes.Buffer{}
		out := &bytes.Buffer{}

		client := NewTelnetClient("localhost:4242", time.Second, io.NopCloser(in), out)
		require.Error(t, client.Connect())
	})
}

func TestInitConfig(t *testing.T) {
	tests := []struct {
		host           string
		port           string
		expectedSocket string
		expectedError  error
	}{
		{
			host:           "localhost",
			port:           "8080",
			expectedSocket: "localhost:8080",
			expectedError:  nil,
		},
		{
			host:           "test.ru",
			port:           "9090",
			expectedSocket: "test.ru:9090",
			expectedError:  nil,
		},
		{
			host:           "127.0.0.1",
			port:           "80",
			expectedSocket: "127.0.0.1:80",
			expectedError:  nil,
		},
		{
			host:           "~localhost",
			port:           "80",
			expectedSocket: "",
			expectedError:  ErrWrongHost,
		},
		{
			host:           "lo`calhos/t",
			port:           "80",
			expectedSocket: "",
			expectedError:  ErrWrongHost,
		},
		{
			host:           "",
			port:           "80",
			expectedSocket: "",
			expectedError:  ErrWrongHost,
		},
		{
			host:           "localhost",
			port:           "",
			expectedSocket: "",
			expectedError:  ErrWrongPort,
		},
		{
			host:           "localhost",
			port:           "-8080",
			expectedSocket: "",
			expectedError:  ErrWrongPort,
		},
		{
			host:           "localhost",
			port:           "10000000",
			expectedSocket: "",
			expectedError:  ErrWrongPort,
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case %d: ", i), func(t *testing.T) {
			tt := tt

			config, err := initConfig(tt.host, tt.port, timeoutDefault)
			require.Equal(t, tt.expectedSocket, config.socket)

			if tt.expectedError != nil {
				require.True(t, errors.Is(err, tt.expectedError))
			} else {
				require.NoError(t, err)
			}
		})
	}
}
