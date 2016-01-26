package main

import (
	"encoding/json"
	log "github.com/Sirupsen/logrus"
	"github.com/allegro/marathon-consul/metrics"
	flag "github.com/ogier/pflag"
	"io/ioutil"
	"strings"
	"time"
)

type Config struct {
	Web        struct {
				   Listen string
			   }
	Metrics    metrics.Config
	Log        struct {
				   Level  string
				   Format string
			   }
	configFile string
}

var _config = &Config{}

func New() (*Config, error) {
	if !flag.Parsed() {
		_config.parseFlags()
	}
	flag.Parse()
	err := _config.loadConfigFromFile()

	if err != nil {
		return nil, err
	}

	_config.setLogFormat()
	err = _config.setLogLevel()

	if err != nil {
		return nil, err
	}

	return _config, err
}

func (config *Config) parseFlags() {

	// Web
	flag.StringVar(&config.Web.Listen, "listen", ":4000", "accept connections at this address")

	// Metrics
	flag.StringVar(&config.Metrics.Target, "metrics-target", "stdout", "Metrics destination stdout or graphite (empty string disables metrics)")
	flag.StringVar(&config.Metrics.Prefix, "metrics-prefix", "default", "Metrics prefix (default is resolved to <hostname>.<app_name>")
	flag.DurationVar(&config.Metrics.Interval, "metrics-interval", 30 * time.Second, "Metrics reporting interval")
	flag.StringVar(&config.Metrics.Addr, "metrics-location", "", "Graphite URL (used when metrics-target is set to graphite)")

	// Log
	flag.StringVar(&config.Log.Level, "log-level", "info", "Log level: panic, fatal, error, warn, info, or debug")
	flag.StringVar(&config.Log.Format, "log-format", "text", "Log format: JSON, text")

	// General
	flag.StringVar(&config.configFile, "config-file", "", "Path to a JSON file to read configuration from. Note: Will override options set earlier on the command line")
}

func (config *Config) loadConfigFromFile() error {
	if config.configFile == "" {
		return nil
	}
	jsonBlob, err := ioutil.ReadFile(config.configFile)
	if err != nil {
		return err
	}
	return json.Unmarshal(jsonBlob, config)
}

func (config *Config) setLogLevel() error {
	level, err := log.ParseLevel(config.Log.Level)
	if err != nil {
		log.WithError(err).WithField("Level", config.Log.Level).Error("Bad level")
		return err
	}
	log.SetLevel(level)
	return nil
}

func (config *Config) setLogFormat() {
	format := strings.ToUpper(config.Log.Format)
	if format == "JSON" {
		log.SetFormatter(&log.JSONFormatter{})
	} else if format == "TEXT" {
		log.SetFormatter(&log.TextFormatter{})
	} else {
		log.WithField("Format", format).Error("Unknown log format")
	}
}
