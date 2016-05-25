package main

import (
	"fmt"
	"os"
	"os/signal"
	"strings"
	"time"
	"net/http"
	"net/url"
)

type IkachanCat struct {
	host		string
	port		string
	channel		string
	msgType		string
	queue		*StreamQ
	shutdown	chan os.Signal
}

func newIkachanCat(host, port, channel, msgType string) *IkachanCat {
	ic := &IkachanCat{
		host:		host,
		port:		port,
		queue:		newStreamQ(),
		shutdown:	make(chan os.Signal, 1),
		channel:	channel,
		msgType:	msgType,
	}

	signal.Notify(ic.shutdown, os.Interrupt)
	return ic
}

func (ic *IkachanCat) trap() {
	sigcount := 0
	for sig := range ic.shutdown {
		if sigcount > 0 {
			fmt.Fprintln(os.Stderr, "aborted")
			os.Exit(1)
		}
		fmt.Fprintf(os.Stderr, "got signal: %s\n", sig.String())
		fmt.Fprintln(os.Stderr, "press ctrl+c again to exit immediately\n")
		sigcount++
		os.Exit(0)
	}
}

func (ic *IkachanCat) exit() {
	for {
		if ic.queue.isEmpty() {
			// os.Exit(0)
		} else {
			time.Sleep(3 * time.Second)
		}
	}
}

func (ic *IkachanCat) addToStreamQ(lines chan string) {
	for line := range lines {
		ic.queue.add(line)
	}
	ic.exit()
}

//TODO: handle messages with length exceeding maximum for Ikachan chat
func (ic *IkachanCat) processStreamQ() {
	if !(ic.queue.isEmpty()) {
		msglines := ic.queue.flush()
		if err := ic.postMsg(msglines); err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err)
		}
	}
	time.Sleep(3 * time.Second)
	ic.processStreamQ()
}

func (ic *IkachanCat) postMsg(msglines []string) error {
	msg := strings.Join(msglines, "\n")
	client := &http.Client{
		Timeout: time.Duration(5) * time.Second,
	}
	values := url.Values{"channel": {ic.channel}, "message": {msg}}
	_, err := client.PostForm(fmt.Sprintf("http://%s:%s/%s", ic.host, ic.port, ic.msgType), values)
	if err != nil {
		return err
	}

	return nil
}

