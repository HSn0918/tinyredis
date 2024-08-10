package config

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
)

var Configures *Config

func init() {

}

var (
	DefaultHost     = "0.0.0.0"
	DefaultPort     = 6379
	DefaultLogDir   = "./"
	DefaultLogLevel = "info"
	defaultShardNum = 1024
)

type Config struct {
	ConfFile string
	Host     string
	Port     int
	LogDir   string
	LogLevel string
	ShardNum int
}
type CfgError struct {
	message string
}

func (err *CfgError) Error() string {
	return err.message
}

func Setup(cmd *cobra.Command) (*Config, error) {
	cfg := &Config{
		Host:     DefaultHost,
		Port:     DefaultPort,
		LogDir:   DefaultLogDir,
		LogLevel: DefaultLogLevel,
		ShardNum: defaultShardNum,
	}
	var err error
	if err = cmd.ParseFlags(os.Args[1:]); err != nil {
		return nil, err
	}
	if cfg.Host, err = cmd.Flags().GetString("host"); err != nil {
		return nil, fmt.Errorf("failed to parse host flag: %w", err)
	}
	if cfg.Port, err = cmd.Flags().GetInt("port"); err != nil {
		return nil, fmt.Errorf("failed to parse port flag: %w", err)
	}
	if cfg.LogDir, err = cmd.Flags().GetString("logdir"); err != nil {
		return nil, fmt.Errorf("failed to parse logdir flag: %w", err)
	}
	if cfg.LogLevel, err = cmd.Flags().GetString("loglevel"); err != nil {
		return nil, fmt.Errorf("failed to parse loglevel flag: %w", err)
	}
	Configures = cfg
	return cfg, nil
}

func (cfg *Config) Parse(cfgFile string) error {
	fl, err := os.Open(cfgFile)
	if err != nil {
		return err
	}

	defer func() {
		err := fl.Close()
		if err != nil {
			fmt.Printf("Close config file error: %s \n", err.Error())
		}
	}()

	reader := bufio.NewReader(fl)
	for {
		line, ioErr := reader.ReadString('\n')
		if ioErr != nil && ioErr != io.EOF {
			return ioErr
		}

		if len(line) > 0 && line[0] == '#' {
			continue
		}

		fields := strings.Fields(line)
		if len(fields) >= 2 {
			cfgName := strings.ToLower(fields[0])
			if cfgName == "host" {
				if ip := net.ParseIP(fields[1]); ip == nil {
					ipErr := &CfgError{
						message: fmt.Sprintf("Given ip address %s is invalid", cfg.Host),
					}
					return ipErr
				}
				cfg.Host = fields[1]
			} else if cfgName == "port" {
				port, err := strconv.Atoi(fields[1])
				if err != nil {
					return err
				}
				if port <= 1024 || port >= 65535 {
					portErr := &CfgError{
						message: fmt.Sprintf("Listening port should be between 1024 and 65535, but %d is given.", port),
					}
					return portErr
				}
				cfg.Port = port
			} else if cfgName == "logdir" {
				cfg.LogDir = strings.ToLower(fields[1])
			} else if cfgName == "loglevel" {
				cfg.LogLevel = strings.ToLower(fields[1])
			} else if cfgName == "shardnum" {
				cfg.ShardNum, err = strconv.Atoi(fields[1])
				if err != nil {
					fmt.Println("ShardNum should be a number. Get: ", fields[1])
					panic(err)
				}
			}
		}
		if ioErr == io.EOF {
			break
		}
	}
	return nil
}

func NewDefaultConfig() *Config {
	return &Config{
		Host:     DefaultHost,
		Port:     DefaultPort,
		LogDir:   DefaultLogDir,
		LogLevel: DefaultLogLevel,
		ShardNum: defaultShardNum,
	}
}
