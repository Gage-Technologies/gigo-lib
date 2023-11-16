package mq

import (
	"fmt"
	"reflect"
	"sort"
	"strings"
	"time"

	"github.com/gage-technologies/gigo-lib/config"
	"github.com/gage-technologies/gigo-lib/logging"
	"github.com/gage-technologies/gigo-lib/mq/streams"
	"github.com/nats-io/nats.go"
)

// JetstreamClient
//
//	Wrapper for nats.JetstreamContext to handle errors,
//	closure of the server connection and stream initialization.
type JetstreamClient struct {
	nats.JetStreamContext
	conn   *nats.Conn
	logger logging.Logger
}

// NewJetstreamClient
//
//	Create a net NewJetstreamClient
func NewJetstreamClient(cfg config.JetstreamConfig, logger logging.Logger) (*JetstreamClient, error) {
	// format options for jetstream connection
	opts := nats.Options{
		Url:            cfg.Host,
		AllowReconnect: true,
		User:           cfg.Username,
		Password:       cfg.Password,
		Compression:    true,
	}

	// connect to jetstream server
	natsConn, err := opts.Connect()
	if err != nil {
		return nil, fmt.Errorf("could not connect to jetstream server: %v", err)
	}

	// create jetstream client
	client := &JetstreamClient{
		conn:   natsConn,
		logger: logger,
	}

	// create jetstream context
	js, err := natsConn.JetStream(
		nats.PublishAsyncMaxPending(cfg.MaxPubQueue),
		nats.PublishAsyncErrHandler(client.errorHandler),
	)
	if err != nil {
		natsConn.Close()
		return nil, fmt.Errorf("could not create jetstream context: %v", err)
	}

	// set jetstream context
	client.JetStreamContext = js

	// initialize streams
	err = client.init()
	if err != nil {
		natsConn.Close()
		return nil, fmt.Errorf("could not initialize streams: %v", err)
	}

	return client, nil
}

// Close
//
//	Closes the internal connection to the jetstream server.
func (c *JetstreamClient) Close() {
	c.conn.Close()
}

// errorHandler
//
//	Callback function to log jetstream publish errors
func (c *JetstreamClient) errorHandler(js nats.JetStream, msg *nats.Msg, err error) {
	// skip for nil error as this should never happen
	if err == nil {
		return
	}
	c.logger.Errorf("failed to publish message to %q: %v", msg.Subject, err)
}

// init
//
//	Initializes streams for the Gigo Core system
func (c *JetstreamClient) init() error {
	// initialize workspace stream
	err := c.initStream(
		streams.StreamWorkspace,
		streams.StreamSubjectsWorkspace,
		streams.RetentionPolicyWorkspace,
		streams.DuplicateFilterWindowWorkspace,
	)
	if err != nil {
		return fmt.Errorf("could not initialize workspace stream: %v", err)
	}

	// initialize session stream
	err = c.initStream(
		streams.StreamMisc,
		streams.StreamSubjectsMisc,
		streams.RetentionPolicyMisc,
		streams.DuplicateFilterWindowMisc,
	)
	if err != nil {
		return fmt.Errorf("could not initialize session stream: %v", err)
	}

	// initialize streak stream
	err = c.initStream(
		streams.StreamStreakXP,
		streams.StreamSubjectsStreak,
		streams.RetentionPolicyStreak,
		streams.DuplicateFilterWindowStreak,
	)
	if err != nil {
		return fmt.Errorf("could not initialize streak stream: %v", err)
	}

	// initialize workspace status stream
	err = c.initStream(
		streams.StreamWorkspaceStatus,
		streams.StreamSubjectsWorkspaceStatus,
		streams.RetentionPolicyWorkspaceStatus,
		streams.DuplicateFilterWindowWorkspaceStatus,
	)
	if err != nil {
		return fmt.Errorf("could not initialize workspace status stream: %v", err)
	}

	// initialize broadcast event stream
	err = c.initStream(
		streams.StreamBroadcastEvent,
		streams.StreamSubjectsBroadcastEvent,
		streams.RetentionPolicyBroadcastEvent,
		streams.DuplicateFilterWindowBroadcastEvent,
	)
	if err != nil {
		return fmt.Errorf("could not initialize broadcast event stream: %v", err)
	}

	err = c.initStream(
		streams.StreamNemesis,
		streams.StreamSubjectsNemesis,
		streams.RetentionPolicyNemesis,
		streams.DuplicateFilterWindowNemesis,
	)
	if err != nil {
		return fmt.Errorf("could not initialize nemesis stream: %v", err)
	}

	err = c.initStream(
		streams.StreamTailscale,
		streams.StreamSubjectsTailscale,
		streams.RetentionPolicyTailscale,
		streams.DuplicateFilterWindowTailscale,
	)
	if err != nil {
		return fmt.Errorf("could not initialize tailscale stream: %v", err)
	}

	err = c.initStream(
		streams.StreamWsConnCache,
		streams.StreamSubjectsWsConnCache,
		streams.RetentionPolicyWsConnCache,
		streams.DuplicateFilterWindowWsConnCache,
	)
	if err != nil {
		return fmt.Errorf("could not initialize tailscale stream: %v", err)
	}

	err = c.initStream(
		streams.StreamChat,
		streams.StreamSubjectsChat,
		streams.RetentionPolicyChat,
		streams.DuplicateFilterWindowChat,
	)
	if err != nil {
		return fmt.Errorf("could not initialize chat stream: %v", err)
	}

	// initialize email stream
	err = c.initStream(
		streams.StreamEmail,
		streams.StreamSubjectsEmail,
		streams.RetentionPolicyEmail,
		streams.DuplicateFilterWindowEmail,
	)
	if err != nil {
		return fmt.Errorf("could not initialize email stream: %v", err)
	}

	return nil
}

// initStream
//
//	Creates or updates a stream for the Gigo Core system
func (c *JetstreamClient) initStream(stream string, subjects []string, retentionPolicy nats.RetentionPolicy,
	duplicateWindow time.Duration) error {
	// attempt to retrieve stream info
	s, _ := c.StreamInfo(stream)

	// create stream if it doesn't exist
	if s == nil {
		_, err := c.AddStream(&nats.StreamConfig{
			Name:      stream,
			Subjects:  subjects,
			Retention: retentionPolicy,
		})
		if err != nil && !strings.Contains(err.Error(), "stream name already in use") {
			return fmt.Errorf("could not create stream: %v", err)
		}

		// exit since we created the stream
		return nil
	}

	// check if we need to reconfigure the stream
	reConfigure := false

	// sort subjects so we can compare them
	sort.Slice(s.Config.Subjects, func(i, j int) bool {
		return s.Config.Subjects[i] < s.Config.Subjects[j]
	})

	// check if subjects match
	if !reflect.DeepEqual(s.Config.Subjects, subjects) {
		reConfigure = true
	}

	// check if retention policy matches
	if s.Config.Retention != retentionPolicy {
		reConfigure = true
	}

	// conditionally reconfigure stream
	if reConfigure {
		_, err := c.UpdateStream(&nats.StreamConfig{
			Name:       stream,
			Subjects:   subjects,
			Retention:  retentionPolicy,
			Duplicates: duplicateWindow,
		})
		if err != nil {
			return fmt.Errorf("could not reconfigure stream: %v", err)
		}
	}

	return nil
}
