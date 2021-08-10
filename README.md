# xts

xts is a script that generates an insert statement from a mysql xml dump file.

## Installation

### Homebrew

```console
$ brew install genkiroid/tap/xts
```

### For other platforms

Precompiled binaries for released versions are available in the [releases](https://github.com/genkiroid/xts/releases) page.

### go get

```console
$ go get github.com/genkiroid/xts/...
```

## Usage

```console
$ xts /path/to/mysql_dump.xml
```

Example is [here](https://github.com/genkiroid/xts/blob/main/xts_example_test.go).

