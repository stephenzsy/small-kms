package acr

import (
	"fmt"
	"net/url"
	"strings"
)

func ExtractACRName(imageRef string) (string, error) {
	host, err := ExtractACRLoginServer(imageRef)
	return strings.Split(host, ".")[0], err
}

func ExtractACRLoginServer(imageRef string) (string, error) {
	if u, err := url.Parse(fmt.Sprintf("https://%s", imageRef)); err != nil {
		return "", err
	} else {
		return u.Host, nil
	}
}
