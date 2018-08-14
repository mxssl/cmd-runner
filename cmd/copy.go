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
	"fmt"
	"log"
	"sync"
	"time"

	"os/exec"

	"github.com/spf13/cobra"
)

// copyCmd represents the copy command
var copyCmd = &cobra.Command{
	Use:   "copy",
	Short: "Copy file from local to remote machine over ssh",
	Long:  `Copy file from local to remote machine over ssh`,
	Run: func(cmd *cobra.Command, args []string) {
		copy()
	},
}

// SourcePath : path to source file
var SourcePath string

// DestinationPath : path to destination file
var DestinationPath string

func init() {
	rootCmd.AddCommand(copyCmd)

	// Read config file
	unmarshallCfg()

	copyCmd.Flags().StringVarP(&SourcePath,
		"source",
		"s",
		c.SourcePath,
		"Path to source file")

	copyCmd.Flags().StringVarP(&DestinationPath,
		"destination",
		"d",
		c.DestinationPath,
		"Path to destination file")
}

func copy() {
	var wg sync.WaitGroup
	start := time.Now()

	for k := range c.Hosts {
		wg.Add(1)
		host := c.Hosts[k]
		go scpCopyFile(host, &wg)
	}
	wg.Wait()
	log.Printf("Elapsed time: %.2fs\n", time.Since(start).Seconds())
}

func scpCopyFile(host string, wg *sync.WaitGroup) {
	defer wg.Done()
	scp(host)
}

func scp(host string) {
	// scp -i ~/.ssh/id_rsa.pub local_file username@remote_host:/remote/path/to/file/filename

	args := []string{}

	dst := fmt.Sprintf(c.Username + "@" + host + ":" + DestinationPath)
	args = append(args, "-i", c.SSHPublicKey, SourcePath, dst)

	log.Printf("Host: %v Copying local file: %v to remote file: %v ...\n", host, SourcePath, DestinationPath)
	out, err := exec.Command("scp", args...).Output()
	if err != nil {
		log.Printf("%s\n", err)
		log.Printf("%s\n", out)
	}
}
