package service

import (
  "fmt"
  "strings"
  "routron/service/backend/http"
)

var (
  ErrMalformed    = fmt.Errorf("Malformed descriptor")
  ErrUnsupported  = fmt.Errorf("Unsupported service")
)

const (
  typeHttp = "http"
)
  

// A service
type Service interface {
  Start()(error)
  Stop()(error)
}

// Create a new service
func New(s string) (Service, error) {
  d, err := parseDef(s)
  if err != nil {
    return nil, err
  }
  switch d.Type {
    case typeHttp:
      return http.New(d.Path)
    default:
      return nil, ErrUnsupported
  }
}

// Service definition
type definition struct {
  Type  string
  Path  string
}

/**
 * Parse a descriptor
 */
func parseDef(s string) (definition, error) {
  var def definition
  
  p := strings.SplitN(s, "=", 2)
  if len(p) != 2 {
    return def, fmt.Errorf("Invalid service: %v", s)
  }
  
  def.Type = strings.TrimSpace(p[0])
  def.Path = strings.TrimSpace(p[1])
  
  return def, nil
}
