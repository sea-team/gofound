package utils

import (
	"bytes"
	"compress/flate"
	"encoding/binary"
	"encoding/gob"
	"fmt"
	"hash/crc32"
	"io/ioutil"
	"log"
	"os"
	"time"
)

func ExecTime(fn func()) float64 {
	start := time.Now()
	fn()
	tc := float64(time.Since(start).Nanoseconds())
	return tc / 1e6
}

// Write 写入二进制数据到磁盘文件
func Write(data interface{}, filename string) {
	buffer := new(bytes.Buffer)
	encoder := gob.NewEncoder(buffer)
	err := encoder.Encode(data)
	if err != nil {
		panic(err)
	}

	log.Println("Write:", filename)
	compressData := Compression(buffer.Bytes())
	err = ioutil.WriteFile(filename, compressData, 0600)
	if err != nil {
		panic(err)
	}
}

func Encoder(data interface{}) []byte {
	if data == nil {
		return nil
	}
	buffer := new(bytes.Buffer)
	encoder := gob.NewEncoder(buffer)
	err := encoder.Encode(data)
	if err != nil {
		panic(err)
	}
	return buffer.Bytes()
}

func Decoder(data []byte, v interface{}) {
	if data == nil {
		return
	}
	buffer := bytes.NewBuffer(data)
	decoder := gob.NewDecoder(buffer)
	err := decoder.Decode(v)
	if err != nil {
		panic(err)
	}
}

// Compression 压缩数据
func Compression(data []byte) []byte {
	buf := new(bytes.Buffer)
	write, err := flate.NewWriter(buf, flate.DefaultCompression)
	defer write.Close()

	if err != nil {
		panic(err)
	}

	write.Write(data)
	write.Flush()
	log.Println("原大小：", len(data), "压缩后大小：", buf.Len(), "压缩率：", fmt.Sprintf("%.2f", float32(buf.Len())/float32(len(data))), "%")
	return buf.Bytes()
}

//Decompression 解压缩数据
func Decompression(data []byte) []byte {
	return DecompressionBuffer(data).Bytes()
}

func DecompressionBuffer(data []byte) *bytes.Buffer {
	buf := new(bytes.Buffer)
	read := flate.NewReader(bytes.NewReader(data))
	defer read.Close()

	buf.ReadFrom(read)
	return buf
}

// Read 从磁盘文件加载二进制数据
func Read(data interface{}, filename string) {
	raw, err := ioutil.ReadFile(filename)
	if err != nil {
		if os.IsNotExist(err) {
			//忽略
			return
		}
		panic(err)
	}
	//解压
	decoData := Decompression(raw)

	buffer := bytes.NewBuffer(decoData)
	dec := gob.NewDecoder(buffer)
	err = dec.Decode(data)
	if err != nil {
		panic(err)
	}
}

// StringToInt 字符串转整数
func StringToInt(value string) uint32 {
	return crc32.ChecksumIEEE([]byte(value))
}

func Uint32Comparator(a, b interface{}) int {
	aAsserted := a.(uint32)
	bAsserted := b.(uint32)
	switch {
	case aAsserted > bAsserted:
		return 1
	case aAsserted < bAsserted:
		return -1
	default:
		return 0
	}
}

func Uint32ToBytes(i uint32) []byte {
	var buf = make([]byte, 4)
	binary.BigEndian.PutUint32(buf, i)
	return buf
}

func BytesToUint32(buf []byte) uint32 {
	return binary.BigEndian.Uint32(buf)
}

// QuickSortAsc 快速排序
func QuickSortAsc(arr []int, start, end int, cmp func(int, int)) {
	if start < end {
		i, j := start, end
		key := arr[(start+end)/2]
		for i <= j {
			for arr[i] < key {
				i++
			}
			for arr[j] > key {
				j--
			}
			if i <= j {
				arr[i], arr[j] = arr[j], arr[i]
				if cmp != nil {
					cmp(i, j)
				}
				i++
				j--
			}
		}

		if start < j {
			QuickSortAsc(arr, start, j, cmp)
		}
		if end > i {
			QuickSortAsc(arr, i, end, cmp)
		}
	}
}
func DeleteArray(array []uint32, index int) []uint32 {
	return append(array[:index], array[index+1:]...)
}
