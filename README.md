<p align='center'>
    <img src='./logo.jpg' width='130px' height='125px'/>
</p>

Wave file reader/writer for Go.

# godoc

https://godoc.org/github.com/takurooo/wavgo

# Examples

## Reader

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
    defer r.Close()

    // read and parse wave file
    err = r.ReadOnMemory()
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

## Writer

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
        for ch := 0; i < int(format.NumChannels); i++ {
            samples[i][ch] = i + ch
        }
    }

    w.WriteSamples(samples)
}
```
