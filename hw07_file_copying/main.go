package main

import (
	"errors"
	"flag"
	"log"
)

var (
	from, to          string
	limit, offset     int64
	ErrInvalidArgFrom = errors.New("invalid file to read from")
	ErrInvalidArgTo   = errors.New("invalid file to write to")
)

func init() {
	flag.StringVar(&from, "from", "", "file to read from")
	flag.StringVar(&to, "to", "", "file to write to")
	flag.Int64Var(&limit, "limit", 0, "limit of bytes to copy")
	flag.Int64Var(&offset, "offset", 0, "offset in input file")
}

func main() {
	flag.Parse()

	if len(from) == 0 {
		log.Println(ErrInvalidArgFrom.Error())

		return
	}

	if len(to) == 0 {
		log.Println(ErrInvalidArgTo.Error())

		return
	}

	err := Copy(from, to, offset, limit)

	if err != nil {
		log.Println(err)
	}
}
