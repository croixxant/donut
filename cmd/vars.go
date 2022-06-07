package cmd

import "path/filepath"

const appName string = "donut"

var defaultConfigPath = filepath.Join("$HOME", ".config", appName, appName+".toml")
var cfgDirPaths = []string{filepath.Join("$XDG_CONFIG_HOME", appName), filepath.Join("$HOME", ".config", appName)}
var ignores = []string{".git"}
