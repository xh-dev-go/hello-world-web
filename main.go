package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/xh-dev-go/hello-world-web/operations"
	"github.com/xh-dev-go/hello-world-web/interfaces"
	"github.com/xh-dev-go/hello-world-web/server"
	"gopkg.in/yaml.v2"
)

var rootCmd = &cobra.Command{
	Use:   "hello-world-web",
	Short: "A simple Go web server and client for inspecting HTTP requests.",
	Long: `hello-world-web is a versatile tool for web development.

It can run as a server to inspect incoming request details,
or act as a client to test endpoints.`,
}

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Launch the web server",
	Long:  `Starts an HTTP server that echoes request information back to the client.`,
	Run: func(cmd *cobra.Command, args []string) {
		server.LaunchServer()
	},
}

var clientCmd = &cobra.Command{
	Use:   "client",
	Short: "Client-side commands to interact with a hello-world-web server",
	Long:  `A collection of commands to test and inspect responses from a running hello-world-web server instance.`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		// This check runs before any subcommand of 'client'
		url, _ := cmd.Flags().GetString("url")
		if url == "" {
			return fmt.Errorf("Error: --url flag is required for all client commands")
		}
		return nil
	},
}

var testCmd = &cobra.Command{
	Use:   "test",
	Short: "Make a test request to a URL and print the response",
	Long:  `Sends a GET request to the specified URL and prints the response body to the console.`,
	Run: func(cmd *cobra.Command, args []string) {
		url, _ := cmd.Flags().GetString("url")
		if url == "" {
			// This check is handled by clientCmd.PersistentPreRunE,
			// but we can keep it as a safeguard.
			log.Fatal("Error: --url flag is required")
		}

		body := operations.GetResponseBody(url)
		fmt.Println(string(body))
	},
}

var getIpCmd = &cobra.Command{
	Use:   "get-ip",
	Short: "Calls a hello-world-web server and prints the client IP address",
	Long:  `Makes a request to a URL (expected to be a hello-world-web server), parses the response, and prints the IP address from the response body.`,
	Run: func(cmd *cobra.Command, args []string) {
		url, _ := cmd.Flags().GetString("url")
		if url == "" {
			log.Fatal("Error: --url flag is required")
		}

		ip, err := operations.GetIp(url)
		if err != nil {
			log.Fatalf("Error: %v", err)
		} else {
			fmt.Println(ip)
		}
	},
}

var getHeadersCmd = &cobra.Command{
	Use:   "get-headers",
	Short: "Calls a hello-world-web server and prints the request headers",
	Long:  `Makes a request to a URL (expected to be a hello-world-web server), parses the response, and prints the headers from the response body as a JSON object.`,
	Run: func(cmd *cobra.Command, args []string) {
		url, _ := cmd.Flags().GetString("url")
		if url == "" {
			log.Fatal("Error: --url flag is required")
		}

		body := operations.GetResponseBody(url)
		var response interfaces.ResponseBody
		err := yaml.Unmarshal(body, &response)
		if err != nil {
			log.Fatalf("Error parsing YAML response from %s: %v\nBody: %s", url, err, string(body))
		}

		headersJson, err := json.MarshalIndent(response.Headers, "", "  ")
		if err != nil {
			log.Fatalf("Error converting headers to JSON: %v", err)
		}

		fmt.Println(string(headersJson))
	},
}

var proxyChainCmd = &cobra.Command{
	Use:   "proxy-chain",
	Short: "Displays the proxy chain of a request to a hello-world-web server",
	Long:  `Makes a request to a URL, parses the response, and displays the visualized proxy chain from the original client to the destination server.`,
	Run: func(cmd *cobra.Command, args []string) {
		url, _ := cmd.Flags().GetString("url")
		if url == "" {
			log.Fatal("Error: --url flag is required")
		}

		body := operations.GetResponseBody(url)
		var response interfaces.ResponseBody
		err := yaml.Unmarshal(body, &response)
		if err != nil {
			log.Fatalf("Error parsing YAML response from %s: %v\nBody: %s", url, err, string(body))
		}

		destination := response.Host

		if forwardedIpsHeader, ok := response.Headers["X-Forwarded-For"]; ok && len(forwardedIpsHeader) > 0 {
			// Case: X-Forwarded-For header exists
			ipList := strings.Split(forwardedIpsHeader[0], ",")
			for i, ip := range ipList {
				ipList[i] = strings.TrimSpace(ip)
			}

			originIp := ipList[0]
			// The full chain includes intermediate proxies and the final connecting IP
			fullChain := append(ipList[1:], response.Ip)

			fmt.Printf("%s -> [ %s ] -> %s\n", originIp, strings.Join(fullChain, " -> "), destination)
		} else {
			// Case: No X-Forwarded-For header
			fmt.Printf("%s -> [ no proxy ] -> %s\n", response.Ip, destination)
		}
	},
}

func main() {
	rootCmd.AddCommand(serverCmd)
	rootCmd.AddCommand(clientCmd)

	// Add the --url flag to the parent 'client' command
	clientCmd.PersistentFlags().StringP("url", "u", "", "URL of the target hello-world-web server (required)")

	// Add the client-side commands as subcommands of 'client'
	clientCmd.AddCommand(testCmd)
	clientCmd.AddCommand(getIpCmd)
	clientCmd.AddCommand(getHeadersCmd)
	clientCmd.AddCommand(proxyChainCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
