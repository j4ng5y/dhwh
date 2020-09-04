package cmd

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
	"time"

	"gopkg.in/yaml.v2"
)

type HTTPServer struct {
	IPAddress         string `yaml:"ip_address"`
	Port              int    `yaml:"port"`
	ReadHeaderTimeout string `yaml:"read_header_timeout"`
	ReadTimeout       string `yaml:"read_timeout"`
	WriteTimeout      string `yaml:"write_timeout"`
	IdleTimeout       string `yaml:"idle_timeout"`
}

func ParseDuration(s string) (time.Duration, error) {
	switch {
	case strings.HasSuffix(s, "sec"):
		s := strings.TrimSuffix(s, "sec")
		i, err := strconv.Atoi(s)
		if err != nil {
			return 0, err
		}
		return time.Duration(i) * time.Second, nil
	case strings.HasSuffix(s, "msec"):
		s := strings.TrimSuffix(s, "msec")
		i, err := strconv.Atoi(s)
		if err != nil {
			return 0, err
		}
		return time.Duration(i) * time.Millisecond, nil
	case strings.HasSuffix(s, "min"):
		s := strings.TrimSuffix(s, "min")
		i, err := strconv.Atoi(s)
		if err != nil {
			return 0, err
		}
		return time.Duration(i) * time.Minute, nil
	case strings.HasSuffix(s, "hr"):
		s := strings.TrimSuffix(s, "hr")
		i, err := strconv.Atoi(s)
		if err != nil {
			return 0, err
		}
		return time.Duration(i) * time.Hour, nil
	case strings.HasSuffix(s, "day"):
		s := strings.TrimSuffix(s, "day")
		i, err := strconv.Atoi(s)
		if err != nil {
			return 0, err
		}
		return time.Duration(i) * (time.Hour * 24), nil
	case strings.HasSuffix(s, "wk"):
		s := strings.TrimSuffix(s, "wk")
		i, err := strconv.Atoi(s)
		if err != nil {
			return 0, err
		}
		return time.Duration(i) * ((time.Hour * 24) * 7), nil
	case strings.HasSuffix(s, "mo"):
		s := strings.TrimSuffix(s, "mo")
		i, err := strconv.Atoi(s)
		if err != nil {
			return 0, err
		}
		return time.Duration(i) * ((time.Hour * 24) * 31), nil
	case strings.HasSuffix(s, "yr"):
		s := strings.TrimSuffix(s, "sec")
		i, err := strconv.Atoi(s)
		if err != nil {
			return 0, err
		}
		return time.Duration(i) * ((time.Hour * 24) * 365), nil
	default:
		return 0, fmt.Errorf("invalid time suffix: %s", s)
	}
}

type Config struct {
	HTTPServer HTTPServer `yaml:"http"`
}

func (C *Config) Unmarshal(filename string) error {
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	return yaml.Unmarshal(b, C)
}
