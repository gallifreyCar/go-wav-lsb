package main

import (
	"fmt"
	"github.com/go-audio/wav"
	"os"
)

func main() {
	f, err := os.Open("testInput/Avicii - The Nights.wav")
	if err != nil {
		panic(fmt.Sprintf("couldn't open audio file - %v", err))
	}

	// Decode the original audio file
	// and collect audio content and information.
	d := wav.NewDecoder(f)
	buf, err := d.FullPCMBuffer()
	if err != nil {
		panic(err)
	}
	f.Close()
	fmt.Println("Old file ->", d)

	// Destination file
	out, err := os.Create("testOutput/Avicii - The Nights.wav")
	if err != nil {
		panic(fmt.Sprintf("couldn't create output file - %v", err))
	}

	// setup the encoder and write all the frames
	e := wav.NewEncoder(out,
		buf.Format.SampleRate,
		int(d.BitDepth),
		buf.Format.NumChannels,
		int(d.WavAudioFormat))

	oldData := buf.Data
	newData := make([]int, len(oldData))
	for i := 0; i < len(oldData); i++ {
		//把这个数字的二进制形式的最后一位变成1
		newData[i] = oldData[i] | 1
	}
	buf.Data = newData
	if err = e.Write(buf); err != nil {
		panic(err)
	}
	// close the encoder to make sure the headers are properly
	// set and the data is flushed.
	if err = e.Close(); err != nil {
		panic(err)
	}
	out.Close()

	// reopen to confirm things worked well
	out, err = os.Open("testOutput/Avicii - The Nights.wav")
	if err != nil {
		panic(err)
	}
	d2 := wav.NewDecoder(out)
	d2.ReadInfo()
	fmt.Println("New file ->", d2)
	out.Close()

}
