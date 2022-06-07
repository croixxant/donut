package internal

import "github.com/spf13/viper"

type Option func(v *viper.Viper) error

func WithFile(file string) Option {
	return func(v *viper.Viper) error {
		v.SetConfigFile(file)
		if err := v.ReadInConfig(); err != nil {
			return err
		}
		return nil
	}
}

func WithNameAndPath(name string, paths ...string) Option {
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

func WithData(data map[string]interface{}) Option {
	return func(v *viper.Viper) error {
		for k, d := range data {
			v.Set(k, d)
		}
		return nil
	}
}
