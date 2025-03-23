package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"strings"
)

type Options struct {
	From      string
	To        string
	Offset    uint
	Limit     int
	blocksize uint
	conv      string
}

func ParseFlags() (*Options, error) {
	var opts Options

	flag.StringVar(&opts.From, "from", "", "file to read. by default - stdin")
	flag.StringVar(&opts.To, "to", "", "file to write. by default - stdout")
	flag.UintVar(&opts.Offset, "offset", 0, "Number of bytes to skip before copying")
	flag.IntVar(&opts.Limit, "limit", -1, "Max bytes to read (-1 for all)")
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

	fmt.Println(*opts)
}
