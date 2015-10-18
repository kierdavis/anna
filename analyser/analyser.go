package analyser

import (
	"math"
	"math/cmplx"
)

type Analysis struct {
	Data []complex128
	SampleRate float64
}

// Return the number of analysed frequencies.
func (a Analysis) Len() int {
	return len(a.Data) / 2
}

// Return information about analysed frequency k.
func (a Analysis) Info(k int) (frequency float64, amplitude float64, phase float64) {
	n := float64(len(a.Data))
	frequency = a.SampleRate * float64(k) / n
	x := a.Data[k]
	return frequency, cmplx.Abs(x) * 2 / n, -cmplx.Phase(x)
}

type Analyser struct {
	SampleRate float64
	BufferSize int
	
	buffer []float64
	window []float64
	permute []int
	twiddle []complex128
	bufferPos int
}

func (a *Analyser) Write(samples []float64) {
	if a.buffer == nil {
		// Initialise the buffer.
		a.buffer = make([]float64, a.BufferSize)
	}
	
	// Feed input samples to circular buffer. Only the last len(a.buffer) samples
	// will actually end up being analysed.
	for len(samples) > 0 {
		n := copy(a.buffer[a.bufferPos:], samples)
		a.bufferPos += n
		samples = samples[n:]
		
		if a.bufferPos == len(a.buffer) {
			a.bufferPos = 0
		}
	}
}

func (a *Analyser) Analyse(dest *Analysis) {
	dest.SampleRate = a.SampleRate
	
	if dest.Data == nil {
		dest.Data = make([]complex128, a.BufferSize)
	}
	
	a.fft(dest.Data)
}

func (a *Analyser) fft(output []complex128) {
	if a.window == nil {
		// Precompute window function
		n := a.BufferSize - 1
		a.window = make([]float64, a.BufferSize)
		for i := range a.window {
			a.window[i] = hamming(float64(i) / float64(n))
		}
	}
	
	if a.permute == nil {
		// Precompute bit-reversal permutation indices
		a.permute = make([]int, a.BufferSize)
		a.initPermute(a.permute, 0, 1)
	}
	
	// Permute input
	for i, j := range a.permute {
		x := a.window[i] * a.buffer[(a.bufferPos + i) % len(a.buffer)]
		output[j] = complex(x, 0)
	}
	
	if a.twiddle == nil {
		// Precompute twiddle factors
		h := a.BufferSize/2
		a.twiddle = make([]complex128, h)
		for k := range a.twiddle {
			θ := -math.Pi * float64(k) / float64(h)
			a.twiddle[k] = complex(math.Cos(θ), math.Sin(θ))
		}
	}
	
	// Apply twiddle factors
	p := 1
	q := a.BufferSize/2
	for q > 0 {
		for j := 0; j < q; j++ {
			b := j * p * 2
			for k := 0; k < p; k++ {
				x, y := output[b + k], output[b + k + p]
				ty := a.twiddle[k * q] * y
				output[b + k] = x + ty
				output[b + k + p] = x - ty
			}
		}
		p *= 2
		q /= 2
	}
}

func (a *Analyser) initPermute(p []int, j int, s int) {
	if len(p) == 1 {
		p[0] = j
	} else {
		a.initPermute(p[:len(p)/2], j, s*2)
		a.initPermute(p[len(p)/2:], j + s, s*2)
	}
}

func hann(θ float64) float64 {
	return 0.5 * (1 - math.Cos(2 * math.Pi * θ))
}

func hamming(θ float64) float64 {
	return 0.54 - 0.46 * math.Cos(2 * math.Pi * θ)
}
