package main

import (
	"time"

	log "github.com/Sirupsen/logrus"
	"gopkg.in/alecthomas/kingpin.v2"

	"github.bus.zalan.do/teapot/mate/consumers"
	"github.bus.zalan.do/teapot/mate/controller"
	"github.bus.zalan.do/teapot/mate/producers"
)

const (
	// TODO: to run synchronize form time to time
	defaultInterval = 10 * time.Minute

	// defaultDomain = "oolong.gcp.zalan.do"

)

var params struct {
	producer string
	consumer string
	sync     bool
	debug    bool
	syncOnly bool
}

var version = "Unknown"

func init() {
	kingpin.Flag("producer", "The endpoints producer to use.").Required().StringVar(&params.producer)
	kingpin.Flag("consumer", "The endpoints consumer to use.").Required().StringVar(&params.consumer)
	kingpin.Flag("sync", "Keep hosted zone in sync and remove records that don't belong to a service.").BoolVar(&params.sync)
	kingpin.Flag("debug", "Enable debug logging.").BoolVar(&params.debug)
	kingpin.Flag("sync-only", "Disable event watcher").BoolVar(&params.syncOnly)
}

func main() {
	kingpin.Version(version)
	kingpin.Parse()

	if params.debug {
		log.SetLevel(log.DebugLevel)
	}

	p, err := producers.New(params.producer)
	if err != nil {
		log.Fatalf("Error creating producer: %v", err)
	}

	c, err := consumers.NewSynced(params.consumer)
	if err != nil {
		log.Fatalf("Error creating consumer: %v", err)
	}

	ctrl := controller.New(p, c, nil)

	if params.sync {
		err = ctrl.Synchronize()
		if err != nil {
			log.Fatalln(err)
		}
	}

	if !params.syncOnly {
		err = ctrl.Watch()
		if err != nil {
			log.Fatalln(err)
		}
	}

	ctrl.Wait()
}
