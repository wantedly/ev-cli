package config

import (
	"github.com/spf13/viper"
)

var (
	Bucket = "wantedly-evaluate" // Default value
)

func InitWithViper(v *viper.Viper) error {
	s := v.GetString("bucket")
	if s != "" {
		Bucket = s
	}
	return nil
}
