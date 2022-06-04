package internal

import "github.com/spf13/viper"

type Option func(v *viper.Viper) error

func WithFile(name string, paths ...string) Option {
	return func(v *viper.Viper) error {
		v.SetConfigName(name)
		for _, path := range paths {
			v.AddConfigPath(path)
		}
		if err := v.ReadInConfig(); err != nil {
			return err
		}
		return nil
	}
}
