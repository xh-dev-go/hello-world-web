package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"gopkg.in/yaml.v2"
	"github.com/spf13/cobra"
	"github.com/xh-dev-go/hello-world-web/interfaces"
	"github.com/xh-dev-go/hello-world-web/server"
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

var testCmd = &cobra.Command{
	Use:   "test",
	Short: "Make a test request to a URL and print the response",
	Long:  `Sends a GET request to the specified URL and prints the response body to the console.`,
	Run: func(cmd *cobra.Command, args []string) {
		url, _ := cmd.Flags().GetString("url")
		if url == "" {
			log.Fatal("Error: --url flag is required")
		}

		resp, err := http.Get(url)
		if err != nil {
			log.Fatalf("Error making request to %s: %v", url, err)
		}
		defer resp.Body.Close()

		fmt.Printf("Response from %s (Status: %s):\n", url, resp.Status)
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Fatalf("Error reading response body: %v", err)
		}
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

		body := getResponseBody(url)
		var response interfaces.ResponseBody
		err := yaml.Unmarshal(body, &response)
		if err != nil {
			log.Fatalf("Error parsing YAML response from %s: %v\nBody: %s", url, err, string(body))
		}

		// Check for X-Forwarded-For header first, as it indicates the original client IP when behind a proxy.
		// The header value can be a comma-separated list; the first IP is the client.
		if forwardedIps, ok := response.Headers["X-Forwarded-For"]; ok && len(forwardedIps) > 0 {
			// The value can be a comma-separated list of IPs. The client is the first one.
			clientIp := strings.Split(forwardedIps[0], ",")[0]
			fmt.Println(strings.TrimSpace(clientIp))
		} else {
			// Fallback to the IP field which is the direct connection IP (RemoteAddr).
			fmt.Println(response.Ip)
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

		body := getResponseBody(url)
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

		body := getResponseBody(url)
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
			fmt.Printf("%s -> [ no proxy ]-> %s\n", response.Ip, destination)
		}
	},
}

// getResponseBody is a helper function to fetch and read the body from a URL.
func getResponseBody(url string) []byte {
	resp, err := http.Get(url)
	if err != nil {
		log.Fatalf("Error making request to %s: %v", url, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("Warning: Received non-200 status code (%s) from %s", resp.Status, url)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Error reading response body: %v", err)
	}
	return body
}

func main() {
	rootCmd.AddCommand(serverCmd)
	rootCmd.AddCommand(testCmd)
	rootCmd.AddCommand(getIpCmd)
	rootCmd.AddCommand(getHeadersCmd)
	rootCmd.AddCommand(proxyChainCmd)

	testCmd.Flags().StringP("url", "u", "", "URL to make a test request to (required)")
	getIpCmd.Flags().StringP("url", "u", "", "URL of the hello-world-web server (required)")
	getHeadersCmd.Flags().StringP("url", "u", "", "URL of the hello-world-web server (required)")
	proxyChainCmd.Flags().StringP("url", "u", "", "URL of the hello-world-web server (required)")

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
