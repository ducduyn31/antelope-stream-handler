package pubsub

import (
	"antelope-stream-chunking/common"
	"antelope-stream-chunking/pipeline"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/tinyzimmer/go-glib/glib"
	"sync"
)

func ProcessStream(signals <-chan int, source string) {
	defer GetSourceManager().TerminateSource(source)
	s := struct {
		mu        sync.Mutex
		frameLeft int
		locked    bool
	}{}

	// Receive first signal for initializing
	v, _ := <-signals
	s.frameLeft = v

	gstElementQuery := []string{
		fmt.Sprintf("rtspsrc location=\"%s\"", source),
		"rtph264depay",
		"avdec_h264",
		"jpegenc",
		"appsink",
	}

	mainLoop := glib.NewMainLoop(glib.MainContextDefault(), false)

	go func() {
		for {
			v, more := <-signals
			if more {
				s.mu.Lock()
				s.frameLeft = v
				s.mu.Unlock()
			} else {
				log.Infof("Closing %s stream", source)
				mainLoop.Quit()
				return
			}

		}
	}()

	go func() {
		if err := pipeline.HandleRTSP(source); err != nil {
			log.Error(err)
		}
		mainLoop.Quit()
	}()

	mainLoop.Run()

	log.Printf("gst-launch-1.0 %s", common.ConcatStringArr(gstElementQuery, " ! "))

}
