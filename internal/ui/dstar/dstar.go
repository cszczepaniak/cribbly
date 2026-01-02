package dstar

import "fmt"

func SendGetf(url string, args ...any) string {
	return send("get", url, args...)
}

func SendPostf(url string, args ...any) string {
	return send("post", url, args...)
}

func SendDeletef(url string, args ...any) string {
	return send("delete", url, args...)
}

func SendPutf(url string, args ...any) string {
	return send("put", url, args...)
}

func send(name, url string, args ...any) string {
	url = fmt.Sprintf(url, args...)
	return fmt.Sprintf("@%s('%s')", name, url)
}
