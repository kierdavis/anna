package main

import (
    "fmt"
    "github.com/kierdavis/anna/analyser"
    "github.com/kierdavis/anna/display"
    "github.com/kierdavis/anna/source/pulsesource"
    "runtime"
    "time"
)

const SampleRate = 44100
const InputBufferSize = 1536
const AnalyseBufferSize = 1 << 13

const ExpectedCycleSecs = float64(InputBufferSize) / float64(SampleRate)

func main() {
    runtime.LockOSThread()
    
    d, err := display.New()
    if err != nil {
        panic(err)
    }
    defer d.Close()
    
    s, err := pulsesource.New(SampleRate)
    if err != nil {
        panic(err)
    }
    
    la := &analyser.Analyser{
        SampleRate: SampleRate,
        BufferSize: AnalyseBufferSize,
    }
    
    ra := &analyser.Analyser{
        SampleRate: SampleRate,
        BufferSize: AnalyseBufferSize,
    }
    
    var lb [InputBufferSize]float64
    var rb [InputBufferSize]float64
    var an analyser.Analysis
    
    cycles := 0
    start := time.Now()
    
    for {
        err = s.Read(lb[:], rb[:])
        if err != nil {
            panic(err)
        }
        
        la.Write(lb[:])
        ra.Write(rb[:])
        
        la.Analyse(&an)
        
        err = d.Draw(an)
        if err != nil {
            panic(err)
        }
        
        cycles++
        expectedTime := start.Add(time.Duration(float64(cycles) * ExpectedCycleSecs * float64(time.Second)))
        lag := time.Now().Sub(expectedTime)
        
        var c string
        if lag > 0 {
            c = "\x1b[1;31m"
        }
        
        fmt.Printf("overall lag: %s%s\x1b[21;39m\n", c, lag)
    }
}
