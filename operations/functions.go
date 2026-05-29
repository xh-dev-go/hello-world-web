package operations

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/xh-dev-go/hello-world-web/interfaces"
	"gopkg.in/yaml.v2"
)

// getResponseBody is a helper function to fetch and read the body from a URL.
func GetResponseBody(url string) []byte {
	resp, err := http.Get(url)
	if err != nil {
		log.Fatalf("Error making request to %s: %v", url, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("Warning: Received non-OK status code (%s) from %s", resp.Status, url)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Error reading response body: %v", err)
	}
	return body
}

func GetIpFromHeaders(resp interfaces.ResponseBody) (string, error) {
	if cfIp, ok := resp.Headers["CF-Connecting-IP"]; ok {
		return cfIp[0], nil
	} else if forwardedIps, ok := resp.Headers["X-Forwarded-For"]; ok && len(forwardedIps) > 0 {
		// The value can be a comma-separated list of IPs. The client is the first one.
		clientIp := strings.Split(forwardedIps[0], ",")[0]
		return clientIp, nil
	} else {
		// Fallback to the IP field which is the direct connection IP (RemoteAddr).
		return resp.Ip, nil	
	}
}

func GetIp(url string) (string, error){
	body:=GetResponseBody(url)
	var response interfaces.ResponseBody
	err := yaml.Unmarshal(body, &response)
	if err != nil {
		return "", fmt.Errorf("Error parsing YAML response from %s: %v\nBody: %s", url, err, string(body))
	}
	
	return GetIpFromHeaders(response);
	// if forwardedIps, ok := response.Headers["X-Forwarded-For"]; ok && len(forwardedIps) > 0 {
	// 	// The value can be a comma-separated list of IPs. The client is the first one.
	// 	clientIp := strings.Split(forwardedIps[0], ",")[0]
	// 	return clientIp, nil
	// } else {
	// 	// Fallback to the IP field which is the direct connection IP (RemoteAddr).
	// 	return response.Ip, nil	
	// }
}
