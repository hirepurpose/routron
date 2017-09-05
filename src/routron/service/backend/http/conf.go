package http

import (
  "io"
  "fmt"
  "sort"
  "strings"
  "net/url"
  "io/ioutil"
)

import (
  "gopkg.in/yaml.v2"
)

/**
 * A route
 */
type route struct {
  Methods     []string            `yaml:"methods"`
  Path        string              `yaml:"path"`
  Headers     map[string]string   `yaml:"headers"`
  Route       string              `yaml:"route"`
}

/**
 * A backend
 */
type backend struct {
  Name        string              `yaml:"name"`
  Host        string              `yaml:"host"`
}

/**
 * A backend configuration
 */
type config struct {
  Listen      string              `yaml:"listen"`
  Backends    []backend           `yaml:"backends"`
  Routes      []route             `yaml:"routes"`
}

/**
 * A route
 */
type Route struct {
  Methods     map[string]struct{}
  Path        string
  Headers     map[string]string
  Route       string
  Backend     *url.URL
  Descr       string
}

/**
 * A backend configuration
 */
type Config struct {
  Listen      string
  Backends    map[string]*url.URL
  Routes      []Route
}

/**
 * Load a config
 */
func NewConfig(src io.ReadCloser) (*Config, error) {
  c := &config{}
  
  data, err := ioutil.ReadAll(src)
  if err != nil {
    return nil, err
  }
  
  err = yaml.Unmarshal(data, c)
  if err != nil {
    return nil, err
  }
  
  conf := &Config{}
  conf.Listen = c.Listen
  
  conf.Backends = make(map[string]*url.URL)
  for _, e := range c.Backends {
    u, err := url.Parse(e.Host)
    if err != nil {
      return nil, fmt.Errorf("Malformed URL: %v", err)
    }
    conf.Backends[e.Name] = u
  }
  
  conf.Routes = make([]Route, len(c.Routes))
  for i, e := range c.Routes {
    var m map[string]struct{}
    var d []string
    
    if len(e.Methods) > 0 {
      m = make(map[string]struct{})
      for _, x := range e.Methods {
        m[strings.ToLower(x)] = struct{}{}
        d = append(d, strings.ToUpper(x))
      }
    }else{
      d = []string{"ANY"}
    }
    
    b, ok := conf.Backends[e.Route]
    if !ok {
      return nil, fmt.Errorf("No such backend: %v", e.Route)
    }
    
    sort.Strings(d)
    conf.Routes[i] = Route{
      Methods: m,
      Path: e.Path,
      Headers: e.Headers,
      Route: e.Route,
      Backend: b,
      Descr: fmt.Sprintf("%s %s", strings.Join(d, ", "), e.Path),
    }
  }
  
  return conf, nil
}

/**
 * Return the first non-nil error or nil if there are none.
 */
func coalesce(err ...error) error {
  for _, e := range err {
    if e != nil {
      return e
    }
  }
  return nil
}
