package sum

import (
	"encoding/json"
	"time"
)

// Config defines the server configuration.
type Config struct {
	Host         string        `json:"host"`
	Port         int           `json:"port"`
	ReadTimeout  time.Duration `json:"read_timeout"`
	WriteTimeout time.Duration `json:"write_timeout"`
	IdleTimeout  time.Duration `json:"idle_timeout"`
}

// DefaultConfig returns a Config with sensible defaults.
func DefaultConfig() Config {
	return Config{
		Host:         "localhost",
		Port:         8080,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  60 * time.Second,
	}
}

// WithHost sets the host address.
func (c Config) WithHost(host string) Config {
	c.Host = host
	return c
}

// WithPort sets the port number.
func (c Config) WithPort(port int) Config {
	c.Port = port
	return c
}

// WithReadTimeout sets the read timeout.
func (c Config) WithReadTimeout(d time.Duration) Config {
	c.ReadTimeout = d
	return c
}

// WithWriteTimeout sets the write timeout.
func (c Config) WithWriteTimeout(d time.Duration) Config {
	c.WriteTimeout = d
	return c
}

// WithIdleTimeout sets the idle timeout.
func (c Config) WithIdleTimeout(d time.Duration) Config {
	c.IdleTimeout = d
	return c
}

// MarshalJSON implements json.Marshaler.
func (c Config) MarshalJSON() ([]byte, error) {
	type alias struct {
		Host         string `json:"host"`
		Port         int    `json:"port"`
		ReadTimeout  string `json:"read_timeout"`
		WriteTimeout string `json:"write_timeout"`
		IdleTimeout  string `json:"idle_timeout"`
	}
	return json.Marshal(alias{
		Host:         c.Host,
		Port:         c.Port,
		ReadTimeout:  c.ReadTimeout.String(),
		WriteTimeout: c.WriteTimeout.String(),
		IdleTimeout:  c.IdleTimeout.String(),
	})
}

// UnmarshalJSON implements json.Unmarshaler.
func (c *Config) UnmarshalJSON(data []byte) error {
	type alias struct {
		Host         string `json:"host"`
		Port         int    `json:"port"`
		ReadTimeout  string `json:"read_timeout"`
		WriteTimeout string `json:"write_timeout"`
		IdleTimeout  string `json:"idle_timeout"`
	}
	var a alias
	if err := json.Unmarshal(data, &a); err != nil {
		return err
	}
	c.Host = a.Host
	c.Port = a.Port
	if a.ReadTimeout != "" {
		d, err := time.ParseDuration(a.ReadTimeout)
		if err != nil {
			return err
		}
		c.ReadTimeout = d
	}
	if a.WriteTimeout != "" {
		d, err := time.ParseDuration(a.WriteTimeout)
		if err != nil {
			return err
		}
		c.WriteTimeout = d
	}
	if a.IdleTimeout != "" {
		d, err := time.ParseDuration(a.IdleTimeout)
		if err != nil {
			return err
		}
		c.IdleTimeout = d
	}
	return nil
}
