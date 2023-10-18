package mq

import (
	"testing"
	"time"

	"github.com/gage-technologies/gigo-lib/config"
	"github.com/gage-technologies/gigo-lib/logging"
	"github.com/gage-technologies/gigo-lib/mq/streams"
)

func TestJetstreamClient(t *testing.T) {
	logger, err := logging.CreateBasicLogger(logging.NewDefaultBasicLoggerOptions("/tmp/gigo-core-js-test.log"))
	if err != nil {
		t.Fatal(err)
	}

	js, err := NewJetstreamClient(config.JetstreamConfig{
		Host:        "mq://127.0.0.1:4222",
		MaxPubQueue: 256,
	}, logger)
	if err != nil {
		t.Fatal(err)
	}

	defer js.Close()

	_, err = js.Publish(streams.SubjectWorkspaceStop, []byte("test"))
	if err != nil {
		t.Fatal(err)
	}

	sub, err := js.SubscribeSync(streams.SubjectWorkspaceStop)
	if err != nil {
		t.Fatal(err)
	}

	msg, err := sub.NextMsg(time.Second)
	if err != nil {
		t.Fatal(err)
	}

	if string(msg.Data) != "test" {
		t.Fatal("incorrect message received:", string(msg.Data))
	}

	if msg.Subject != streams.SubjectWorkspaceStop {
		t.Fatal("incorrect subject:", msg.Subject)
	}

	err = msg.AckSync()
	if err != nil {
		t.Fatal(err)
	}
}
