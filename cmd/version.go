/*
 * Copyright 2016 Frank Wessels <fwessels@xs4all.nl>
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

// versionCmd represents the ls command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "List files in the repository",
	Long:  "List files in the repository",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("S3Git")
		fmt.Println("")
		fmt.Println(fmt.Sprintf("%-10s %s", "Version:", version))
		fmt.Println(fmt.Sprintf("%-10s %s", "Commit:", commit))
		fmt.Println(fmt.Sprintf("%-10s %s", "Date:", date))
	},
}

func init() {
	RootCmd.AddCommand(versionCmd)
}
