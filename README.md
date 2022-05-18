# donut

Tiny dotfiles management tool written in Go.

# Usage

1. Clone dotfiles into your `src_dir`

```
git clone git@github.com:croixxant/dotfiles.git $XDG_DATA_HOME/donut
```

2. Create donut's config file and specify your `src_dir`.

```
echo "src_dir = \"$XDG_DATA_HOME/donut\"" > $XDG_CONFIG_HOME/donut.toml
```
