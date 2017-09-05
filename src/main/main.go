package main

import (
  "os"
  "fmt"
  "flag"
  "strings"
  "routron/service"
)

import (
  "github.com/bww/go-util/debug"
)

/**
 * You know what it does
 */
func main() {
  os.Exit(app())
}

/**
 * Actually run the app
 */
func app() int {
  var svcs int
  
  cmdline   := flag.NewFlagSet(os.Args[0], flag.ExitOnError)
  fDebug    := cmdline.Bool     ("debug",         strToBool(os.Getenv("ROUTRON_DEBUG")),        "Enable debugging mode.")
  fVerbose  := cmdline.Bool     ("verbose",       strToBool(os.Getenv("ROUTRON_VERBOSE")),      "Be more verbose.")
  cmdline.Parse(os.Args[1:])
  
  debug.DEBUG = *fDebug
  debug.VERBOSE = *fVerbose
  
  for _, e := range cmdline.Args() {
    svc, err := service.New(e)
    if err != nil {
      fmt.Printf("* * * Could not create proxy service: %v\n", err)
      return 1
    }
    
    err = svc.Start()
    if err != nil {
      fmt.Printf("* * * Could not start proxy service: %v\n", err)
      return 1
    }
    
    svcs++
    defer svc.Stop()
    fmt.Printf("----> Service %v\n", svc)
  }
  
  if svcs < 1 {
    fmt.Println("* * * No services defined")
    return 1
  }
  
  <- make(chan struct{})
  return 0
}

/**
 * Return the first non-empty string from those provided
 */
func coalesce(v... string) string {
  for _, e := range v {
    if e != "" {
      return e
    }
  }
  return ""
}

/**
 * String to bool
 */
func strToBool(s string, d ...bool) bool {
  if s == "" {
    if len(d) > 0 {
      return d[0]
    }else{
      return false
    }
  }
  return strings.EqualFold(s, "t") || strings.EqualFold(s, "true") || strings.EqualFold(s, "y") || strings.EqualFold(s, "yes")
}
