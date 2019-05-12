// Copyright © 2019 Daniel Hodges <hodges.daniel.scott@gmail.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"fmt"
	"log"
	"net/http"
	"os"

	exporter "github.com/hodgesds/perf_exporter"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile string
)

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "perf_exporter",
	Short: "Prometheus exporter for perf events",
	Long:  ``,
	RunE: func(cmd *cobra.Command, args []string) error {
		config := viper.GetViper()
		collector, err := exporter.NewPerfCollector(config)
		if err != nil {
			return err
		}
		go func() {
			if err := collector.Start(); err != nil {
				log.Fatal(err)
			}
		}()

		prometheus.MustRegister(
			collector.(prometheus.Collector),
		)

		metricsPath := config.GetString("metrics-path")
		http.Handle(metricsPath, prometheus.Handler())
		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte(`<html>
             <head><title>Perf Exporter</title></head>
             <body>
             <h1>Perf Exporter</h1>
             <p><a href='` + metricsPath + `'>Metrics</a></p>
             </body>
             </html>`))
		})
		return http.ListenAndServe(config.GetString("listen-address"), nil)
	},
}

// Execute adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func init() {
	RootCmd.PersistentFlags().StringVarP(
		&cfgFile,
		"config", "c",
		"",
		"config file (default is $HOME/.perf_exporter.yaml)",
	)
	RootCmd.PersistentFlags().String(
		"metrics-path",
		"/metrics",
		"Metrics endpoint",
	)
	RootCmd.PersistentFlags().StringP(
		"listen-address", "l",
		"0.0.0.0:8585",
		"Server listen address",
	)
	viper.BindPFlags(RootCmd.PersistentFlags())
	cobra.OnInitialize(initConfig)
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	viper.AutomaticEnv()
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
		if err := viper.ReadInConfig(); err != nil {
			fmt.Printf("Failed to read config file: %v", err)
		}
	}

	viper.SetConfigName("config")
	viper.AddConfigPath("/etc/perf_exporter")
	if err := viper.ReadInConfig(); err == nil {
		fmt.Printf("Using config file: %v", viper.ConfigFileUsed())
	}
}
