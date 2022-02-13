package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strings"
)

func getIP() (string, error) {
	var ip string

	resp, err := http.Get("https://icanhazip.com")
	if err != nil {
		log.Print(err)
		return ip, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("failed to get IP address, response code: %d", resp.StatusCode)
		return ip, err
	}

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Print(err)
		return ip, err
	}
	ip = strings.TrimSpace(string(b))

	if !validIP(ip) {
		maxLen := 10
		if len(ip) < maxLen {
			maxLen = len(ip)
		}
		return ip, fmt.Errorf("not a valid ip: %s", sanitise(ip[0:maxLen-1]))
	}

	return ip, nil
}

func validIP(ip string) bool {
	re, _ := regexp.Compile(`^(([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])\.){3}([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])$`)
	return re.MatchString(ip)
}

func sanitise(s string) string {
	// limit valid characters
	reg, _ := regexp.Compile("[^a-zA-Z0-9<>!#='\"()]+")
	return reg.ReplaceAllString(s, " ")
}
