package webssh

import (
	"io/ioutil"
	"log"
	"time"
)

// https://github.com/fanb129/Kube-CC

var (
	DefaultTerm        = TermXterm
	DefaultConnTimeout = 15 * time.Second
	DefaultLogger      = log.New(ioutil.Discard, "[webssh] ", log.Ltime|log.Ldate)
	DefaultBuffSize    = uint32(8192)
)
