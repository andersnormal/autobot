package plugins

import (
	"context"
	"fmt"

	pb "github.com/andersnormal/autobot/proto"

	"github.com/nats-io/stan.go"
)

// Plugin ...
type Plugin interface {
	Subscribe(context.Context, func(e *pb.Event)) error
}

type plugin struct {
	env *Env
}

// Plugins ...
func New(e *Env) Plugin {
	p := new(plugin)
	p.env = e

	return p
}

// Subscribe ...
func (p *plugin) Subscribe(ctx context.Context, fn func(e *pb.Event)) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// connect to NATs
	sc, err := stan.Connect(p.env.Get(AutobotClusterID), p.env.Get(AutobotClusterURL))
	if err != nil {
		return err
	}

	defer sc.Close()

	// subscribe
	sub, _ := sc.Subscribe(p.env.Get(AutobotTopic), func(m *stan.Msg) {
		fmt.Printf("Received a message: %s\n", string(m.Data))
	})

	defer sub.Unsubscribe()

	<-ctx.Done()

	return nil
}
