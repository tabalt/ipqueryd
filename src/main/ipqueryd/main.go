package main

import (
	"flag"
	"log"
	"os"

	"github.com/tabalt/pidfile"
)

var (
	cf string
)

func main() {
	// parse flag
	flag.StringVar(&cf, "c", "", "config file for ipqueryd")
	flag.Parse()

	// parse conf
	conf, err := parseIpquerydConf(cf)
	if err != nil {
		log.Printf("parse ipqueryd conf error: %v", err)
		os.Exit(1)
	}

	//create pidfile
	pf, err := pidfile.CreatePidFile(conf.PidFile)
	if err != nil {
		log.Printf("create pid file %s failed, error: %v", conf.PidFile, err)
		os.Exit(1)
	}

	// start server
	log.Printf("try to start ipqueryd server. ")
	srvErr := startServer(conf)
	if srvErr != nil {
		log.Printf("ipqueryd server stopped, error: %v", srvErr)
		os.Exit(1)
	}

	//clear pidfile
	if err := pidfile.ClearPidFile(pf); err != nil {
		log.Printf("clear pid file failed, error: %v", err)
	}
}
