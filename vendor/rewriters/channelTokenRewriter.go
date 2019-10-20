package rewriters

import (
  "regexp"
  "fmt"
)

func ChannelTokenRewriter (data string) string {

  // Pattern match for channel token
  // (mandatory) <-
  ctRegex := regexp.MustCompile(`<<`);

  if( ctRegex.Match( []byte(data) ) ){
    // Use regex to replace channel token
    data = ctRegex.ReplaceAllString(data, "<-");
    fmt.Println("[SYNC] Worker - Channel Token Updated");
  }

  // Finished, return updated string
  return data;

}
