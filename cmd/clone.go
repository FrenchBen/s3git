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
	"os"
	"path/filepath"
	"strings"

	"github.com/cheggaaa/pb/v3"
	"github.com/dustin/go-humanize"
	"github.com/s3git/s3git-go"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// cloneCmd represents the clone command
var cloneCmd = &cobra.Command{
	Use:   "clone [resource]",
	Short: "Clone a repository into a new directory",
	Long:  "Clone a repository into a new directory",
	Run: func(cmd *cobra.Command, args []string) {

		if len(args) == 0 {
			er("Missing resource to clone from")
		}

		parts := strings.Split(args[0], "//")
		if len(parts) != 2 {
			er(fmt.Sprintf("Bad resource for cloning (missing '//' separator): %s", args[0]))
		}

		dir, err := filepath.Abs(filepath.Dir("."))
		if err != nil {
			er(err)
		}

		dir += "/" + parts[1]

		// Check whether directory to clone into does not yet exist -- abort otherwise
		if _, err := os.Stat(dir); err == nil {
			er(fmt.Sprintf("Cannot clone into existing directory: %s", dir))
		}

		// Output directory and create it
		fmt.Println("Cloning into", dir)
		err = os.MkdirAll(dir, os.ModePerm)
		if err != nil {
			er(err)
		}

		var barDownloading, barProcessing *pb.ProgressBar

		progressDownload := func(total int64) {
			if barDownloading == nil {
				barDownloading = pb.New64(total).Start()
				barDownloading.Set("prefix", "Downloading ")
			}
			barDownloading.Increment()
			if barDownloading.Current() == total {
				barDownloading.Finish()
			}
		}

		progressProcessing := func(total int64) {
			if barProcessing == nil {
				barProcessing = pb.New64(total).Start()
				barDownloading.Set("prefix", "Processing ")
			}
			barProcessing.Increment()
			if barProcessing.Current() == total {
				barProcessing.Finish()
			}
		}

		options := []s3git.CloneOptions{}
		options = append(options, s3git.CloneOptionSetAccessKey(viper.GetString(ACCESS_KEY)))
		options = append(options, s3git.CloneOptionSetSecretKey(viper.GetString(SECRET_KEY)))
		options = append(options, s3git.CloneOptionSetEndpoint(viper.GetString(ENDPOINT)))
		options = append(options, s3git.CloneOptionSetDownloadProgress(progressDownload))
		options = append(options, s3git.CloneOptionSetProcessingProgress(progressProcessing))

		repo, err := s3git.Clone(args[0], dir, options...)
		if err != nil {
			er(err)
		}

		outputStats(repo)
	},
}

func outputStats(repo *s3git.Repository) {

	stats, err := repo.Statistics()
	if err == nil {
		fmt.Printf("Done. Totaling %s object%s.\n", humanize.Comma(int64(stats.Objects)), pluralize(stats.Objects))
	}
}

func pluralize(number uint64) string {
	if number != 1 {
		return "s"
	}
	return ""
}

func init() {
	RootCmd.AddCommand(cloneCmd)

	// Add local message flags
}
