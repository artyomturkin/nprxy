// Copyright © 2018 Artyom Turkin
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package cmd

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"

	"github.com/artyomturkin/nprxy"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Start nprxy",
	Long:  `Proxy requests to the target service and from the service to upstream dependencies.`,
	RunE:  run,
}

func init() {
	rootCmd.AddCommand(runCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// runCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// runCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func run(cmd *cobra.Command, args []string) error {
	c := &nprxy.Config{}
	err := viper.Unmarshal(c)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithCancel(context.Background())

	if len(c.Services) == 0 {
		fmt.Printf("Must provide at least one service in config file\n")
		os.Exit(1)
	}

	for _, service := range c.Services {
		go func(s nprxy.ServiceConfig) {
			log.Printf("Starting %s", s.Name)
			err := nprxy.ProxyService(ctx, s)
			if err != http.ErrServerClosed {
				fmt.Printf("Proxy service %s failed: %v\n", s.Name, err)
				os.Exit(2)
			}
		}(service)
	}

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)
	<-ch
	cancel()
	return nil
}
