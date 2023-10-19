package cmd

import (
	"fmt"
	"strings"

	"github.com/ory/viper"
	"github.com/spf13/cobra"
	fn "knative.dev/func/pkg/functions"
)

func NewSubscribeCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:        "subscribe",
		Short:      "Subscribe to a function",
		Long:       `Subscribe to a function`,
		SuggestFor: []string{"subscribe", "subscribe"},
		PreRunE:    bindEnv("filter", "source"),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runSubscribe(cmd, args)
		},
	}

	cmd.Flags().StringArrayP("filter", "f", []string{}, "Filter for the Cloud Event metadata")

	cmd.Flags().StringP("source", "s", "default", "The source, like a Knative Broker")

	return cmd
}

func runSubscribe(cmd *cobra.Command, args []string) (err error) {
	var (
		cfg subscibeConfig
		f   fn.Function
	)
	cfg = newSubscribeConfig(cmd)

	if f, err = fn.NewFunction(""); err != nil {
		return
	}
	if !f.Initialized() {
		return fn.NewErrNotInitialized(f.Root)
	}
	if !f.Initialized() {
		return fn.NewErrNotInitialized(f.Root)
	}

	// add subscription	to function
	f.Subscription = append(f.Subscription, fn.SubscriptionSpec{
		Source:  cfg.Source,
		Filters: extractFilterMap(cfg),
	})

	// pump it
	return f.Write()

}

func extractFilterMap(cfg subscibeConfig) map[string]string {
	subscriptionFilters := make(map[string]string)
	for _, filter := range cfg.Filter {
		kv := strings.Split(filter, "=")
		if len(kv) != 2 {
			fmt.Println("Invalid pair:", filter)
			continue
		}
		key := kv[0]
		value := kv[1]
		subscriptionFilters[key] = value
	}
	return subscriptionFilters
}

type subscibeConfig struct {
	Filter []string
	Source string
}

func newSubscribeConfig(cmd *cobra.Command) (c subscibeConfig) {
	c = subscibeConfig{
		Filter: viper.GetStringSlice("filter"),
		Source: viper.GetString("source"),
	}
	// NOTE: .Filter should be viper.GetStringSlice, but this returns unparsed
	// results and appears to be an open issue since 2017:
	// https://github.com/spf13/viper/issues/380
	var err error
	if c.Filter, err = cmd.Flags().GetStringArray("filter"); err != nil {
		fmt.Fprintf(cmd.OutOrStdout(), "error reading filter arguments: %v", err)
	}

	return
}