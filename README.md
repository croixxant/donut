# donut

Tiny dotfiles management tool written in Go.

# Install


```
go install github.com/croixxant/donut@latest
```

# Usage

1. Clone dotfiles into your `src_dir`

```
git clone git@github.com:croixxant/dotfiles.git $XDG_DATA_HOME/donut
```

2. Create donut's config file and specify your `src_dir`.

```
donut init $XDG_DATA_HOME/donut
```

3. List files to be applied

```
donut list
```

4. Apply files

```
donut apply
```
