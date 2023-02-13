/*
Copyright © 2023 Julien CAGNIART

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"github.com/julien040/gut/src/controller"
	"github.com/spf13/cobra"
)

// telemetryCmd represents the telemetry command
var telemetryCmd = &cobra.Command{
	Use:   "telemetry [enable|disable]",
	Short: "Enable or disable telemetry",
	Long: `Gut collects anonymous usage data to help us improve the product.
This data is stored with BigQuery and is only accessible by the Gut team.
You can enable or disable telemetry by running the following command:
	gut telemetry enable
	gut telemetry disable

Here is the data we collect:
	- Gut version
	- Gut command
	- OS
	- Architecture
	- CPU cores
	- Country

You can find more information about telemetry here:
https://gut-cli.dev/docs

If you have any questions, please contact us at:
contact@gut-cli.dev
`,
	Run: controller.Telemetry,
}

func init() {
	rootCmd.AddCommand(telemetryCmd)
}
