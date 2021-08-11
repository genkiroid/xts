# xts

xts is a script that generates an insert statement from a mysql xml dump file.

## Use case

xts **is not** a tool for restoring database.If your goal is to restore, you may want to refer to [this](https://rpbouman.blogspot.com/2010/04/restoring-xml-formatted-mysql-dumps.html).

xts is designed for tiny use cases such as the following.

- In languages other than Go, MySQL dump files in xml format are used as fixtures for unit tests.
- I want to reuse the above fixtures in unit tests for a project whose implementation language is Go.

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
$ xts path-to-mysql-dump.xml
```

Example is [here](https://github.com/genkiroid/xts/blob/main/xts_example_test.go).

