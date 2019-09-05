package nats_test

import (
	"bytes"
	"context"
	"io/ioutil"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/andersnormal/autobot/pkg/config"
	. "github.com/andersnormal/autobot/pkg/nats"

	nats "github.com/nats-io/nats.go"
	stan "github.com/nats-io/stan.go"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestNew_Start(t *testing.T) {
	ready := make(chan bool)

	var buf bytes.Buffer
	log.SetOutput(&buf)

	fn := func() {
		var readyOnce sync.Once
		readyOnce.Do(func() {
			ready <- true
		})
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cfg := config.New()

	// only will use temp dir for tests...
	cfg.DataDir, _ = ioutil.TempDir("", "")
	cfg.Debug = true
	cfg.Verbose = true
	s := New(cfg, Timeout(5*time.Second))

	defer func() { _ = os.RemoveAll(cfg.NatsFilestoreDir()) }()

	go s.Start(ctx, fn)() // todo: have a testing wrapper in pkg

	<-ready

	assert.Equal(t, s.ClusterID(), cfg.Nats.ClusterID)

	nc, err := nats.Connect(
		cfg.Nats.ClusterURL,
		nats.MaxReconnects(-1),
		nats.ReconnectBufSize(-1),
	)
	assert.NoError(t, err)

	defer nc.Close()

	sc, err := stan.Connect(
		s.ClusterID(),
		"foo",
		stan.NatsConn(nc),
	)
	assert.NoError(t, err)

	defer sc.Close()

	err = sc.Publish(cfg.Nats.Inbox, []byte("test"))

	var msg []byte
	exit := make(chan struct{})

	sub, err := sc.QueueSubscribe(cfg.Nats.Inbox, "foo", func(m *stan.Msg) {
		msg = m.Data
		m.Ack()
		exit <- struct{}{}
	}, stan.SetManualAckMode(), stan.DurableName("foo"), stan.StartWithLastReceived())
	assert.NoError(t, err)

	defer sub.Unsubscribe()

	select {
	case <-exit:
	case <-time.After(5 * time.Second):
	}

	assert.Equal(t, []byte("test"), msg)
}
