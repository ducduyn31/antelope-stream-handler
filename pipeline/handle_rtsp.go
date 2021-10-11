package pipeline

import (
	"github.com/tinyzimmer/go-gst/gst"
	"github.com/tinyzimmer/go-gst/gst/app"
)


//  rtspsrc location=SOURCE ! rtph264depay ! avdec_h264 ! jpegenc
func createPipeline(source string) (*gst.Pipeline, error) {
	gst.Init(nil)
	pipeline, err := gst.NewPipeline("")
	if err != nil {
		return nil, err
	}

	// rtspsrc location=SOURCE
	src, err := gst.NewElement("rtspsrc")
	if err != nil {
		return nil, err
	}
	src.SetArg("location", source)

	// rtph264depay
	h264Decoder, err := gst.NewElement("rtph264depay")
	if err != nil {
		return nil, err
	}

	// avdec_h264
	avDecoder, err := gst.NewElement("avdec_h264")
	if err != nil {
		return nil, err
	}

	// jpegenc
	jpegEncoder, err := gst.NewElement("jpegenc")

	// appsink
	sink, err := app.NewAppSink()
	if err != nil {
		return nil, err
	}

	sink.SetCallbacks(&app.SinkCallbacks{
		NewSampleFunc: func(appSink *app.Sink) gst.FlowReturn {
			sample := sink.PullSample()

			if sample == nil {
				return gst.FlowEOS
			}

			buffer := sample.GetBuffer()
			if buffer == nil {
				return gst.FlowError
			}

		},
	})

	// Link elements
	err = pipeline.AddMany(src, h264Decoder, avDecoder, jpegEncoder, sink.Element)
	if err != nil {
		return nil, err
	}
	err = src.Link(h264Decoder)
	if err != nil {
		return nil, err
	}
	err = h264Decoder.Link(avDecoder)
	if err != nil {
		return nil, err
	}
	err = avDecoder.Link(jpegEncoder)
	if err != nil {
		return nil, err
	}

	return pipeline, nil
}

func mainLoop(pipeline *gst.Pipeline) error {
	err := pipeline.SetState(gst.StatePlaying)
	if err != nil {
		return err
	}

	bus := pipeline.GetPipelineBus()


	return nil
}

func HandleRTSP(source string) error {
	var pipeline *gst.Pipeline
	var err error
	if pipeline, err = createPipeline(source); err != nil {
		return err
	}
	return mainLoop(pipeline)
}