package config

import (
	"errors"
	"flag"
	"fmt"
	"io/fs"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type (
	Config struct {
		Log      `yaml:"Log" env-prefix:"BOT_" env-description:"Logging configuration"`
		Users    map[int64]string `yaml:"Users" env-prefix:"BOT_" env-required:"true" env-description:"Users configuration"`
		BotToken string           `yaml:"BotToken" env:"BOT_TOKEN" env-required:"true" env-description:"Telegram bot token"`
	}

	Log struct {
		// Path to log file
		LogFile string `yaml:"LogFile" env:"LOGFILE" env-default:"./bot.log" env-description:"Path to log file"`
	}
)

func NewConfig() (*Config, error) {
	var cfg Config

	// create flag set using `flag` package
	fset := flag.NewFlagSet("tb", flag.ExitOnError)

	var configPath string
	fset.StringVar(&configPath, "c", "./config.yml", "Path to config.yml")
	//fset.StringVar(&help, "h", "", "Help")

	// get config usage with wrapped flag usage
	fset.Usage = cleanenv.FUsage(fset.Output(), &cfg, nil, fset.Usage)
	if err := fset.Parse(os.Args[1:]); err != nil {

		fset.Usage()
		os.Exit(0)

	}

	if err := cleanenv.ReadConfig(configPath, &cfg); errors.Is(err, fs.ErrNotExist) {
		//Reread env vars if file not found
		if err := cleanenv.ReadEnv(&cfg); err != nil {
			return nil, fmt.Errorf("env config error %w", err)
		}

	} else if err != nil {
		return nil, fmt.Errorf("read config error %w", err)
	}

	return &cfg, nil

}
