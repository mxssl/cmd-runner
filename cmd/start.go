// Copyright © 2018 mxssl
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
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"sync"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/crypto/ssh"
)

// Config struct for toml config file
type Config struct {
	Username     string   `mapstructure:"username"`
	Password     string   `mapstructure:"password"`
	Hosts        []string `mapstructure:"hosts"`
	CommandsFile string   `mapstructure:"commands_file"`
}

var c Config

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start run commands on hosts from config file",
	Long:  `Read config and start run commands`,
	Run: func(cmd *cobra.Command, args []string) {
		start()
	},
}

func init() {
	rootCmd.AddCommand(startCmd)
}

func unmarshallCfg() {

	v := viper.New()

	v.SetConfigName("config")
	v.AddConfigPath(".")

	if err := v.ReadInConfig(); err != nil {
		fmt.Printf("couldn't load config: %s", err)
		os.Exit(1)
	}

	if err := v.Unmarshal(&c); err != nil {
		fmt.Printf("couldn't read config: %s", err)
	}
}

func start() {
	unmarshallCfg()
	log.Println("Running...")
	log.Printf("Hosts: %v\n", c.Hosts)
	sshRun()
}

func sshRun() {
	config := &ssh.ClientConfig{
		User: c.Username,
		Auth: []ssh.AuthMethod{
			ssh.Password(c.Password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         time.Second * 5,
	}

	commands := readCommands()
	log.Printf("Commands: %v", commands)

	var wg sync.WaitGroup
	start := time.Now()
	for k := range c.Hosts {
		wg.Add(1)
		host := c.Hosts[k] + ":22"
		go sendCommands(config, host, commands, &wg)
	}
	wg.Wait()
	log.Printf("%.2fs elapsed\n", time.Since(start).Seconds())
}

func sendCommands(config *ssh.ClientConfig, host string, commands string, wg *sync.WaitGroup) {
	log.Printf("Running commands on host: %v", host)
	defer wg.Done()
	client, err := ssh.Dial("tcp", host, config)
	if err != nil {
		log.Println(err.Error())
		log.Printf("Cannot connect to host: %v", host)
		return
	}
	session, err := client.NewSession()
	if err != nil {
		log.Println(err)
	}
	defer session.Close()

	var b bytes.Buffer
	session.Stdout = &b

	if err := session.Run(commands); err != nil {
		log.Println(err)
	}

	outputFileName := host[:len(host)-3] + "-output.txt"
	f, err := os.Create(outputFileName)
	if err != nil {
		log.Print(err)
	}
	defer func() {
		if err := f.Close(); err != nil {
			panic(err)
		}
	}()
	if _, err := f.WriteString(b.String()); err != nil {
		log.Print(err)
	}
}

func readCommands() (commands string) {
	b, err := ioutil.ReadFile(c.CommandsFile)
	if err != nil {
		log.Fatal(err)
	}
	commands = string(b)
	return commands
}
