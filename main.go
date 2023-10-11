package main

import (
	"fmt"
	"github.com/go-audio/wav"
	"math/rand"
	"os"
	"strconv"
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

	// modify the data
	oldData := buf.Data
	x := buildX(len(oldData))
	newData := make([]int, len(oldData))
	for i := 0; i < len(x); i++ {
		newData[i] = handler(oldData[i], x[i])
	}
	buf.Data = newData

	// write the modified data
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

func handler(n int, x int) int {
	oldStr := fmt.Sprintf("%b", n)
	newStr := oldStr[:len(oldStr)-1] + strconv.Itoa(x)
	newInt, _ := strconv.ParseInt(newStr, 2, 64)
	return int(newInt)
}

func buildX(len int) []int {
	x := make([]int, len)
	for i := 0; i < len; i++ {
		x[i] = rand.Intn(2)
	}
	return x
}
