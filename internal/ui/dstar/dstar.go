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
	// Setting requestCancellation allows requests to finish even if they were cancelled, which
	// happens if the element that triggered the event is removed from the DOM.
	return fmt.Sprintf("@%s('%s', { requestCancellation: 'disabled' })", name, url)
}
