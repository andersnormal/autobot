package cmd

import (
	"context"
	"log"
	"time"

	pb "github.com/andersnormal/autobot/proto"
	"github.com/andersnormal/autobot/cli/config"
	"github.com/andersnormal/autobot/cli/dialer"

	"github.com/spf13/cobra"
	"google.golang.org/grpc"
)

var ListPlugins = &cobra.Command{
	Use:   "list-plugins",
	Short: "Lists Plugins",
	Run: func(cmd *cobra.Command, args []string) {
		var opts []grpc.CallOption
		var err error

		dial, err := dialer.NewDialer(config.C)
		if err != nil {
			log.Fatal(err)
		}
		defer dial.Close()

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		client := pb.NewAutobotClient(dial)
		in := &pb.ListDomains_Request{}

		resp, err := client.ListPlugins(ctx, in, opts...)
		if err != nil {
			log.Fatal(err)
		}

		log.Println(resp.Plugins)
	},
}
