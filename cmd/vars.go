package cmd

var cfgDirPaths = []string{"$XDG_CONFIG_HOME/" + appName, "$HOME/.config/" + appName, "$HOME"}
var ignores = []string{mapConfigName, ".git"}
