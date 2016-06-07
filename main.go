package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/jessevdk/go-flags"
)

const version = "0.1.0"

type opts struct {
	Host    string `short:"H" long:"host" value-name:"hostname" description:"ikachan hostname"`
	Port    string `short:"p" long:"port" value-name:"port" default:"4979" description:"ikachan port"`
	Channel string `short:"c" long:"channel" value-name:"'#channel'" description:"destination channel"`
	MsgType string `short:"t" long:"type" value-name:"msgtype" default:"notice" description:"message type notice/privmsg)"`
	Stream  bool   `short:"s" long:"stream" description:"messages to Ikachan continuously"`
	Version bool   `short:"V" long:"version" description:"show version"`
}

func readIn(lines chan string) {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		lines <- scanner.Text()
	}
	close(lines)
}

func printVersion() {
	fmt.Printf("ikachancat version v%s\n", version)
}

func main() {
	o := &opts{}
	p := flags.NewParser(o, flags.Default)
	p.Usage = "--host HOSTNAME --channel '#CHANNEL' [--port=PORT] [--type=MSGTYPE] \n\nVerion: " + version
	_, err := p.ParseArgs(os.Args)
	if err != nil {
		p.WriteHelp(os.Stderr)
		os.Exit(-1)
	}

	if o.Version {
		printVersion()
		os.Exit(0)
	}

	if o.Host == "" {
		fmt.Fprintf(os.Stderr, "--host required\n")
		p.WriteHelp(os.Stderr)
		os.Exit(1)
	}
	if o.Channel == "" {
		fmt.Fprintf(os.Stderr, "--channel required\n")
		p.WriteHelp(os.Stderr)
		os.Exit(1)
	}

	lines := make(chan string)
	go readIn(lines)

	if o.Stream {
		ikachancat := newIkachanCat(o.Host, o.Port, o.Channel, o.MsgType)
		go ikachancat.addToStreamQ(lines)
		go ikachancat.processStreamQ()
		go ikachancat.trap()
		select {}
	} else {
		fmt.Fprintf(os.Stderr, "currentry --stream required")
		os.Exit(-1)
	}

	os.Exit(0)
}
