# wavgo

[![codecov](https://codecov.io/gh/takurooo/wavgo/graph/badge.svg?token=63MH3X9PC4)](https://codecov.io/gh/takurooo/wavgo)

<p align="center">
  <img src="./logo.jpg" width="130px" height="125px" />
</p>

`wavgo` is a small library for reading and writing WAV audio files in Go.

## Features

- Parse WAV headers and access format information
- Read sample data in common bit depths
- Write new WAV files with custom formats

## Install

```bash
go get github.com/takurooo/wavgo
```

## Documentation

Full API documentation is available on [pkg.go.dev](https://pkg.go.dev/github.com/takurooo/wavgo).

## Examples

### Reader

```go
package main

import (
    "fmt"

    "github.com/takurooo/wavgo"
)

func main() {
    var err error
    r := wavgo.NewReader()
    if err = r.Open("/path/to/your/file.wav"); err != nil {
        panic(err)
    }
    defer func() {
        if err := r.Close(); err != nil {
            panic(err)
        }
    }()

    // read and parse wave file
    err = r.Load()
    if err != nil {
        panic(err)
    }

    // get format info
    format := r.GetFormat()

    fmt.Println("AudioFormat    : ", format.AudioFormat)
    fmt.Println("NumChannels    : ", format.NumChannels)
    fmt.Println("SampleRate     : ", format.SampleRate)
    fmt.Println("ByteRate       : ", format.ByteRate)
    fmt.Println("BlockAlign     : ", format.BlockAlign)
    fmt.Println("BitsPerSample  : ", format.BitsPerSample)

    // get sample data
    samples, err := r.GetSamples(2)
    for _, sample := range samples {
        for ch := 0; ch < int(format.NumChannels); ch++ {
            fmt.Println(sample[ch])
        }
    }
}
```

### Writer

```go
package main

import (
    "github.com/takurooo/wavgo"
)

func main() {
    var err error

    format := &wavgo.Format{}
    format.AudioFormat = wavgo.AudioFormatPCM
    format.NumChannels = 2
    format.SampleRate = 48000
    format.ByteRate = 128000
    format.BlockAlign = 4
    format.BitsPerSample = 16

    w := wavgo.NewWriter(format)
    err = w.Open("test.wav")
    if err != nil {
        panic(err)
    }
    defer w.Close()

    samples := make([]wavgo.Sample, 4)
    for i := 0; i < len(samples); i++ {
        for ch := 0; ch < int(format.NumChannels); ch++ {
            samples[i][ch] = i + ch
        }
    }

    w.WriteSamples(samples)
}
```
