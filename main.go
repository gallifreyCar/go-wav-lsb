package main

import (
	"fmt"
	"github.com/go-audio/wav"
	"os"
	"strconv"
	"strings"
)

func main() {

	input := "你好，世界"
	x := StringToX(input)
	fmt.Println("x ->", x)
	err := EncodeLSB("testdata/test.wav", x)
	if err != nil {
		fmt.Println(err)
	}
	x2, err := DecodeLSB("testdata/test_encoded.wav")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("x2 ->", x2)
	fmt.Println("x2 ->", XToString(x2))
	//fmt.Println("len(x2) ->", len(x2))

}

// HandleData 加密Data
func HandleData(n int, x int) int {
	oldStr := fmt.Sprintf("%b", n)
	newStr := oldStr[:len(oldStr)-1] + strconv.Itoa(x)
	newInt, _ := strconv.ParseInt(newStr, 2, 64)
	return int(newInt)
}

// StringToX 字符串转换为二进制随机序列（LSB）
func StringToX(s string) []int {
	//输入：你好，世界
	//输出：
	b := []byte(s)
	//用两个字节存储字符串的长度
	xHeader := make([]int, 16)
	for i := 0; i < 16; i++ {
		xHeader[i] = len(b) >> uint(15-i) & 1
	}
	fmt.Println("xHeader ->", xHeader)
	fmt.Println("len(b) ->", len(b))

	//下面的循环是将每个字节转换为8个bit，存入x中
	x := make([]int, len(b)*8)
	for i := 0; i < len(b); i++ {
		for j := 0; j < 8; j++ {
			x[i*8+j] = int(b[i] >> uint(7-j) & 1)
		}
	}

	res := append(xHeader, x...)
	return res
}

// XToString 二进制随机序列（LSB）转换为字符串
func XToString(x []int) string {

	var b []byte
	for i := 0; i < len(x)/8; i++ {
		var tmp byte
		for j := 0; j < 8; j++ {
			//tmp是一个字节，每次左移一位，然后加上x[i*8+j]的最后一位
			tmp = tmp | byte(x[i*8+j])<<uint(7-j)
		}
		b = append(b, tmp)
	}
	return string(b)
}

func getXLen(xHeader []int) int {
	var length int
	for i := 0; i < 16; i++ {
		length = length | xHeader[i]<<uint(15-i)
	}
	return length * 8
}

func EncodeLSB(file string, x []int) error {
	// Destination file
	f, err := os.Open(file)
	if err != nil {
		return err
	}

	// Decode the original audio file
	// and collect audio content and information.
	d := wav.NewDecoder(f)
	buf, err := d.FullPCMBuffer()
	if err != nil {
		return err
	}
	f.Close()

	newFile := strings.Replace(file, ".wav", "_encoded.wav", 1)

	// Destination file
	out, err := os.Create(newFile)
	if err != nil {
		return err
	}

	// setup the encoder and write all the frames
	e := wav.NewEncoder(out,
		buf.Format.SampleRate,
		int(d.BitDepth),
		buf.Format.NumChannels,
		int(d.WavAudioFormat))

	// modify the data
	data := buf.Data
	for i := 0; i < len(x); i++ {
		data[i] = HandleData(data[i], x[i])
	}
	buf.Data = data

	// write the modified data
	if err = e.Write(buf); err != nil {
		return err
	}
	// close the encoder to make sure the headers are properly
	// set and the data is flushed.
	if err = e.Close(); err != nil {
		return err
	}
	out.Close()

	//// reopen to confirm things worked well
	//out, err = os.Open(newFile)
	//if err != nil {
	//	return err
	//}
	//d2 := wav.NewDecoder(out)
	//d2.ReadInfo()
	//fmt.Println("New file ->", d2)
	//out.Close()
	return nil
}

func DecodeLSB(file string) (x []int, err error) {
	// Destination file
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}

	// Decode the original audio file
	// and collect audio content and information.
	d := wav.NewDecoder(f)
	buf, err := d.FullPCMBuffer()
	if err != nil {
		return nil, err
	}
	f.Close()

	// get modified X
	data := buf.Data
	xLen := getXLen(data[:16])
	x = make([]int, xLen)
	for i := 0; i < xLen; i++ {
		x[i] = data[i+16] & 1
	}
	return x, nil
}
