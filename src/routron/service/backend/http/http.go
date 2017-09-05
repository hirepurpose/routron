package http

import (
  "os"
  "fmt"
  "time"
  "strings"
  "net/http"
  "net/http/httputil"
)

import (
  "github.com/bww/go-util/path"
  "github.com/bww/go-util/debug"
  // "github.com/davecgh/go-spew/spew"
)

// Don't wait forever
const ioTimeout = time.Second * 10

// Default to http
const defaultScheme = "http"

// HTTP proxy service
type httpService struct {
  conf    *Config
  server  *http.Server
}

// Create a new service
func New(p string) (*httpService, error) {
  
  f, err := os.Open(p)
  if err != nil {
    return nil, err
  }
  defer f.Close()
  
  conf, err := NewConfig(f)
  if err != nil {
    return nil, err
  }
  
  return &httpService{conf:conf}, nil
}

// Stringer
func (s *httpService) String() string {
  return fmt.Sprintf("http://%v", s.conf.Listen)
}

// Start the service
func (s *httpService) Start() error {
  if s.server != nil {
    return fmt.Errorf("Service is running")
  }
  
  proxy := &httputil.ReverseProxy{
    Director: s.routeRequest,
  }
  
  s.server = &http.Server{
    Addr: s.conf.Listen,
    Handler: proxy,
    ReadTimeout: ioTimeout,
    WriteTimeout: ioTimeout,
    MaxHeaderBytes: 1 << 20,
  }
  
  go func(){
    err := s.server.ListenAndServe()
    if err != nil && err != http.ErrServerClosed {
      panic(err)
    }
  }()
  
  return nil
}

// Stop the service
func (s *httpService) Stop() error {
  if s.server == nil {
    return fmt.Errorf("Service is not running")
  }
  err := s.server.Close()
  s.server = nil
  return err
}

// Handle requests
func (s *httpService) routeRequest(req *http.Request) {
  var err error
  for _, e := range s.conf.Routes {
    var match bool
    if e.Methods == nil {
      match = true
    }else if _, ok := e.Methods[strings.ToLower(req.Method)]; ok {
      match = true
    }
    if match {
      match, err = path.Match(e.Path, req.URL.Path)
      if err != nil {
        panic(fmt.Errorf("Invalid path pattern: %v: %v\n", req.URL, err))
      }else if match {
        s.rewriteRequest(req, e)
        return
      }
    }
  }
  if debug.VERBOSE {
    fmt.Printf("http: %s %s -> <no match>\n", req.Method, req.URL.Path)
  }
}

// Handle requests
func (s *httpService) rewriteRequest(req *http.Request, route Route) {
  b := route.Backend
  if debug.VERBOSE {
    fmt.Printf("http: %s %s -> %v @ %v (%s)\n", req.Method, req.URL.Path, route.Route, b, route.Descr)
  }
  
  if b.Scheme != "" {
    req.URL.Scheme = b.Scheme
  }else{
    req.URL.Scheme = defaultScheme
  }
  
  req.URL.Host = b.Host
  req.URL.Path = joinPath(b.Path, req.URL.Path)
  
  req.Host = b.Host
  req.Header.Set("Host", b.Host)
}

// Join a path without extraneous slashes
func joinPath(a, b string) string {
  aslash := strings.HasSuffix(a, "/")
  bslash := strings.HasPrefix(b, "/")
  switch {
    case aslash && bslash:
      return a + b[1:]
    case !aslash && !bslash:
      return a + "/" + b
    default:
      return a + b
  }
}
