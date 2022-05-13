package utils

import (
	"bytes"
	"compress/flate"
	"encoding/binary"
	"encoding/gob"
	"fmt"
	"io/fs"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
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

const (
	c1 = 0xcc9e2d51
	c2 = 0x1b873593
	c3 = 0x85ebca6b
	c4 = 0xc2b2ae35
	r1 = 15
	r2 = 13
	m  = 5
	n  = 0xe6546b64
)

var (
	Seed = uint32(1)
)

func Murmur3(key []byte) (hash uint32) {
	hash = Seed
	iByte := 0
	for ; iByte+4 <= len(key); iByte += 4 {
		k := uint32(key[iByte]) | uint32(key[iByte+1])<<8 | uint32(key[iByte+2])<<16 | uint32(key[iByte+3])<<24
		k *= c1
		k = (k << r1) | (k >> (32 - r1))
		k *= c2
		hash ^= k
		hash = (hash << r2) | (hash >> (32 - r2))
		hash = hash*m + n
	}

	var remainingBytes uint32
	switch len(key) - iByte {
	case 3:
		remainingBytes += uint32(key[iByte+2]) << 16
		fallthrough
	case 2:
		remainingBytes += uint32(key[iByte+1]) << 8
		fallthrough
	case 1:
		remainingBytes += uint32(key[iByte])
		remainingBytes *= c1
		remainingBytes = (remainingBytes << r1) | (remainingBytes >> (32 - r1))
		remainingBytes = remainingBytes * c2
		hash ^= remainingBytes
	}

	hash ^= uint32(len(key))
	hash ^= hash >> 16
	hash *= c3
	hash ^= hash >> 13
	hash *= c4
	hash ^= hash >> 16

	// 出发吧，狗嬷嬷！
	return
}

// StringToInt 字符串转整数
func StringToInt(value string) uint32 {
	return Murmur3([]byte(value))
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

func ReleaseAssets(file fs.File, out string) {
	if file == nil {
		return
	}

	if out == "" {
		panic("out is empty")
	}

	//判断out文件是否存在
	if _, err := os.Stat(out); os.IsNotExist(err) {
		//读取文件信息
		fileInfo, err := file.Stat()
		if err != nil {
			panic(err)
		}
		buffer := make([]byte, fileInfo.Size())
		_, err = file.Read(buffer)
		if err != nil {
			panic(err)
		}

		// 读取输出文件目录
		outDir := filepath.Dir(out)
		err = os.MkdirAll(outDir, os.ModePerm)
		if err != nil {
			panic(err)
		}

		//创建文件
		outFile, _ := os.Create(out)
		defer func(outFile *os.File) {
			err := outFile.Close()
			if err != nil {
				panic(err)
			}
		}(outFile)

		err = ioutil.WriteFile(out, buffer, os.ModePerm)
		if err != nil {
			panic(err)
		}
	}

}

// DirSizeB DirSizeMB getFileSize get file size by path(B)
func DirSizeB(path string) int64 {
	var size int64
	filepath.Walk(path, func(_ string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			size += info.Size()
		}
		return err
	})

	return size
}

//getFileSize get file size by path(B)
func getFileSize(path string) int64 {
	if !exists(path) {
		return 0
	}
	fileInfo, err := os.Stat(path)
	if err != nil {
		return 0
	}
	return fileInfo.Size()
}

//exists Whether the path exists
func exists(path string) bool {
	_, err := os.Stat(path)
	return err == nil || os.IsExist(err)
}

// RemovePunctuation 移除所有的标点符号
func RemovePunctuation(str string) string {
	reg := regexp.MustCompile(`\p{P}+`)
	return reg.ReplaceAllString(str, "")
}

// RemoveSpace 移除所有的空格
func RemoveSpace(str string) string {
	reg := regexp.MustCompile(`\s+`)
	return reg.ReplaceAllString(str, "")
}

func contains(s *[]string, e string, skipIndex int) bool {
	for index, a := range *s {
		if index != skipIndex && strings.Contains(a, e) {
			return true
		}
	}
	return false
}

// GetLongWords 获取长词
func GetLongWords(words *[]string) []string {

	var newWords = make([]string, 0)
	for index, w := range *words {
		if !contains(words, w, index) {
			newWords = append(newWords, w)
		}
	}
	return newWords
}
