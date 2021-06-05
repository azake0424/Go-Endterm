package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"sort"
	"time"
)

type reader struct {
	words *[]record
}

type record struct {
	word    []byte
	counter int
}

func (r *reader) contains(element []byte) (bool, int) {
	for index, v := range *r.words {
		if bytes.Equal(v.word, element) {
			return true, index
		}
		index = index + 1
	}
	return false, 0
}

func (r *reader) read_from_chan(ch chan []byte) {
	for node := range ch {
		state, index := r.contains(node)
		if state {
			(*r.words)[index].counter++
		} else {
			record := record{node, 1}
			*r.words = append(*r.words, record)
		}
	}
}

func (r *reader) get20mostfrequentwords() {
	sort.Slice(*r.words, func(i, j int) bool {
		return (*r.words)[i].counter > (*r.words)[j].counter
	})
	for i := 0; i < 20; i++ {
		fmt.Println((*r.words)[i].counter, " ", string((*r.words)[i].word))
	}
}

func main() {
	start := time.Now()
	file, err := os.Open("mobydick.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err = file.Close(); err != nil {
			log.Fatal(err)
		}
	}()
	readingBuf := make([]byte, 1)

	words := make([]record, 0)
	reader := reader{words: &words}

	writingBuf := make([]byte, 0)

	ch := make(chan []byte)

	r := bufio.NewReader(file)
	go func() {
		for {
			n, err := r.Read(readingBuf)

			if n > 0 {
				byteVal := readingBuf[0]

				if byteVal >= 65 && byteVal <= 90 {

					byteVal = byteVal + 32
					writingBuf = append(writingBuf, byteVal)

				} else if byteVal >= 97 && byteVal <= 122 {

					writingBuf = append(writingBuf, byteVal)

				} else if byteVal == 32 && len(writingBuf) != 0 {

					ch <- writingBuf
					writingBuf = nil

				} else if ((byteVal > 122 || byteVal < 65) || (byteVal > 90 && byteVal < 97)) && len(writingBuf) != 0 {

					ch <- writingBuf
					writingBuf = nil

				} else {
					continue
				}
			}

			if err == io.EOF {
				ch <- writingBuf
				writingBuf = nil
				break
			}
		}
		close(ch)
	}()

	reader.read_from_chan(ch)
	reader.get20mostfrequentwords()
	fmt.Printf("Process took %s\n", time.Since(start))
}
