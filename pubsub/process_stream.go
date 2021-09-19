package pubsub

import (
	"antelope-stream-chunking/messaging"
	log "github.com/sirupsen/logrus"
	"gocv.io/x/gocv"
	"sync"
)

// ProcessStreamRTSP TODO: Implement RTSP to minimize decoding process
//func ProcessStreamRTSP(signals <-chan int, source string) {
//	defer GetSourceManager().TerminateSource(source)
//
//	s := struct {
//		mu        sync.Mutex
//		frameLeft int
//		locked    bool
//	}{}
//
//	// Receive first signal
//	v, _ := <-signals
//	s.frameLeft = v
//
//	// Parse streaming url source
//	u, err := base.ParseURL(source)
//	if err != nil {
//		log.Errorf("Something happened for %s: %s", source, err)
//		return
//	}
//
//	// Connecting to RTSP server
//	conn, err := gortsplib.Dial(u.Scheme, u.Host)
//	defer func(conn *gortsplib.ClientConn) {
//		_ = conn.Close()
//	}(conn)
//
//	tracks, baseUrl, _, err := conn.Describe(u)
//	if err != nil {
//		log.Errorf("Something happened for %s: %s", source, err)
//		return
//	}
//
//	var h264TrackId int
//
//	// Select video track only
//	for id, track := range tracks {
//		if track.Media.MediaName.Media == "video" && track.IsH264() {
//			_, err := conn.Setup(headers.TransportModePlay, baseUrl, track, 0, 0)
//			h264TrackId = id
//			if err != nil {
//				log.Errorf("Something happened for %s: %s", source, err)
//				return
//			}
//		}
//	}
//
//	_, err = conn.Play(nil)
//
//	// Stop stream and reconnect stream
//	go func() {
//		for {
//			v, more := <-signals
//			if more {
//				s.mu.Lock()
//				s.frameLeft = v
//				s.mu.Unlock()
//			} else {
//				log.Infof("Closing %s stream", source)
//				s.mu.Lock()
//				s.locked = true
//				_ = conn.Close()
//				s.mu.Unlock()
//				return
//			}
//		}
//	}()
//	order := 0
//
//	h264Dec := rtph264.NewDecoder()
//
//	err = conn.ReadFrames(func(trackID int, streamType gortsplib.StreamType, payload []byte) {
//		if s.frameLeft > 0 {
//			if streamType != gortsplib.StreamTypeRTP || trackID != h264TrackId {
//				return
//			}
//			nalus, _, err := h264Dec.Decode(payload)
//			if err != nil {
//				return
//			}
//			s.frameLeft--
//
//			for _, nalu := range nalus {
//				go messaging.SendToKafka(nalu, order, source)
//			}
//			order++
//		} else if order != 0 {
//			order = 0
//		}
//	})
//	log.Errorf("Stop streaming: %s", err)
//}

func ProcessStreamCV(signals <-chan int, source string) {
	defer GetSourceManager().TerminateSource(source)
	s := struct {
		mu        sync.Mutex
		frameLeft int
		locked    bool
	}{}

	// Receive first signal
	v, _ := <-signals
	s.frameLeft = v

	vCap, err := gocv.OpenVideoCapture(source)

	if err != nil {
		log.Error(err)
		return
	}

	// Stop stream and reconnect stream
	go func() {
		for {
			v, more := <-signals
			if more {
				s.mu.Lock()
				s.frameLeft = v
				s.mu.Unlock()
			} else {
				log.Infof("Closing %s stream", source)
				s.mu.Lock()
				s.locked = true
				if err != nil {
					log.Error(err)
				}
				s.mu.Unlock()
				return
			}
		}
	}()

	frame := gocv.NewMat()
	order := int64(0)
	consecutiveFails := 0
	for {
		if vCap.IsOpened() && !s.locked && consecutiveFails < 100 {
			vCap.Grab(1)
			if s.frameLeft <= 0 {
				continue
			}
			ok := vCap.Read(&frame)

			if !ok {
				log.Errorf("Failed to read frame %d", order)
				consecutiveFails++
				continue
			} else {
				consecutiveFails = 0
			}
			go messaging.SendToKafka(&frame, order, source)
			s.frameLeft--
			order++
		} else {
			log.Infof("%s closed, consecutive failure: %d", source, consecutiveFails)
			_ = vCap.Close()
			break
		}
	}
}
