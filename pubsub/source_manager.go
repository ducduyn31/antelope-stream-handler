package pubsub

import (
	"antelope-stream-chunking/common"
	log "github.com/sirupsen/logrus"
	"os"
	"strconv"
	"sync"
)

var lock = &sync.Mutex{}

type sourceManager struct {
	Sources        []string
	SourceStatus   map[string]bool
	handlerManager *PubSub
	batchSize      int
}

var instance *sourceManager

func GetSourceManager() *sourceManager {
	if instance == nil {
		lock.Lock()
		defer lock.Unlock()
		if instance == nil {
			batchSize, err := strconv.Atoi(os.Getenv("BATCH_SIZE"))
			if err != nil {
				log.Fatal("BATCH_SIZE must be a number")
				panic(err)
			}
			instance = &sourceManager{
				Sources:        nil,
				SourceStatus:   make(map[string]bool),
				handlerManager: CreatePubSub(),
				batchSize:      batchSize,
			}
		}
	}
	return instance
}

func (s *sourceManager) LoadSourcesFromFile(path string) {
	sources, f, err := common.OpenFileAndReadLines(path)
	if err != nil {
		panic(err)
	}
	defer func(f *os.File) {
		_ = f.Close()
	}(f)

	s.Sources = sources
	for _, source := range sources {
		s.SourceStatus[source] = false
	}
}

func (s *sourceManager) AddSourceToMem(source string) {
	if _, ok := s.SourceStatus[source]; !ok {
		s.Sources = append(s.Sources, source)
		s.SourceStatus[source] = false
	}
}

func (s *sourceManager) ActivateAllSources() {
	for _, source := range s.Sources {
		log.Infof("Activating %s", source)
		s.ActivateSource(source)
	}
}

func (s *sourceManager) ActivateSource(source string) {
	if running, ok := s.SourceStatus[source]; ok {
		go s.handlerManager.Signal(source, s.batchSize)
		if !running {
			ch := s.handlerManager.Subscribe(source)
			go ProcessStreamCV(ch, source)
			s.SourceStatus[source] = true
		}
	} else {
		log.Infof("%s is not available", source)
	}
}

func (s *sourceManager) TerminateSource(source string) {
	if running, ok := s.SourceStatus[source]; ok && running {
		s.handlerManager.CloseOne(source)
		s.SourceStatus[source] = false
	}
}

func (s *sourceManager) RemoveSource(source string) {
	if _, ok := s.SourceStatus[source]; ok {
		s.TerminateSource(source)
		delete(s.SourceStatus, source)
	}
}
