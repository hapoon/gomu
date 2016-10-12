gomu
=========
[![GoDoc](https://godoc.org/github.com/hapoon/gomu?status.png)](https://godoc.org/github.com/hapoon/gomu)
[![Build Status](https://travis-ci.org/hapoon/gomu.svg?branch=master)](https://travis-ci.org/hapoon/gomu)
[![Coverage Status](https://coveralls.io/repos/github/hapoon/gomu/badge.svg?branch=master)](https://coveralls.io/github/hapoon/gomu?branch=master)
[![MIT License](http://img.shields.io/badge/license-MIT-blue.svg?style=flat)](https://raw.githubusercontent.com/hapoon/gomu/master/LICENSE)

gomu is a library for handling of nullable values and unassigned parameter.

## Installation

Make sure that Go is installed on your computer. Type the following command in your terminal:

`go get gopkg.in/hapoon/gomu.v1`

After it the package is ready to use.

Add following line in your `*.go` file:

```go
import "gopkg.in/hapoon/gomu.v1"
```

## Usage

### String

Nullable string.

If JSON value is null, String.Null is true.
If JSON key is not assigned, String.Valid is false.

### Int

Nullable int64.

If JSON value is null, Int.Null is true.
If JSON key is not assigned, Int.Valid is false.

### Bool

Nullable bool.

If JSON value is null, Bool.Null is true.
If JSON key is not assigned, Bool.Valid is false.

### Time

Nullable Time.

If JSON value is null, Time.Null is true.
If JSON key is not assigned, Time.Valid is false.

### Validate

```go
type exampleStruct struct {
    Name String `valid:"stringlength(1|10)"`
    URL String `valid:"requrl"`
}
example := exampleStruct{
    Name: "sample name",
    URL:  "https://sample.com/top/",
}
result, err := gomu.Validate(example)
if !result {
    if err != nil {
        println("error: " + err.Error())
    }
}
```

## License

[MIT License](LICENSE)

