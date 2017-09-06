package http

import (
  "regexp"
)

// Scheme matcher
var urlish = regexp.MustCompile("^[a-z]+://")

// Determine if a string appears to be a URL
func isURL(u string) bool {
  return urlish.MatchString(u)
}
