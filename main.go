package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
)

type Options struct {
	From      string
	To        string
	Offset    uint
	Limit     uint
	blocksize uint
	conv      string
}

func ParseFlags() (*Options, error) {

	var opts Options

	flag.StringVar(&opts.From, "from", "", "file to read. by default - stdin")
	flag.StringVar(&opts.To, "to", "", "file to write. by default - stdout")
	flag.UintVar(&opts.Offset, "offset", 0, "Number of bytes to skip before copying")
	flag.UintVar(&opts.Limit, "limit", 0, "Max bytes to read (0 for all)")
	flag.UintVar(&opts.blocksize, "blocksize", 4096, "Block size for reading/writing")
	flag.StringVar(&opts.conv, "conv", "", "upper_case, lower_case, trim_spaces")

	flag.Parse()

	conv := strings.Split(opts.conv, ",")
	upper := false
	lower := false
	for _, str := range conv {
		str = strings.TrimSpace(str)
		if str == "" {
			continue
		}
		switch str {
		case "upper_case":
			upper = true
		case "lower_case":
			lower = true
		}
	}
	if upper && lower {
		return nil, errors.New("cannot use both upper_case and lower_case")
	}
	return &opts, nil
}

func main() {
	opts, err := ParseFlags()
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, "can not parse flags:", err)
		os.Exit(1)
	}

	fileFrom, err := os.OpenFile(opts.From, os.O_RDONLY, 0644)
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, "can not open file:", err)
		os.Exit(1)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
		}
	}(fileFrom)

	fileTo, err := os.OpenFile(opts.To, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, "can not open/create file:", err)
		os.Exit(1)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
		}
	}(fileTo)

	//offset
	offset := opts.Offset

	fi, err := fileFrom.Stat()
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, "cannot accuses file info:", err)
	}

	if offset > 0 {
		if int64(offset) > fi.Size() {
			err := errors.New("offset is greater than file length")
			_, _ = fmt.Fprintln(os.Stderr, "offset error:", err)
		}
		_, err = fileFrom.Read(make([]byte, offset))
	}

	buf := make([]byte, opts.blocksize)
	//limit
	limit := opts.Limit
	var cnt uint
	bs := opts.blocksize
	//Запись в файл
	for {
		cnt += bs
		if cnt > limit {
			limBuf := make([]byte, limit%bs)
			_, err = fileFrom.Read(limBuf)
			if err != nil {
				if err == io.EOF {
					break
				}
				_, _ = fmt.Fprintln(os.Stderr, "cannot read file:", err)
			}
			_, err := fileTo.Write(limBuf)
			if err != nil {
				_, _ = fmt.Fprintln(os.Stderr, "can not write to file:", err)
			}
			break
		}
		_, err = fileFrom.Read(buf)
		if err != nil {
			if err == io.EOF {
				break
			}
			_, _ = fmt.Fprintln(os.Stderr, "cannot read file:", err)
		}
		_, err := fileTo.Write(buf)
		if err != nil {
			_, _ = fmt.Fprintln(os.Stderr, "can not write to file:", err)
		}
		fmt.Println(string(buf))
	}
}
