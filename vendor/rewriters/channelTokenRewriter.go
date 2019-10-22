package rewriters

import (
  "regexp"
  "fmt"
)

var channelToken string = "<<";

func ChannelTokenRewriter(data string) string{

  // Pattern match for channel token
  // (mandatory) <-
  ctRegex := regexp.MustCompile(channelToken);

  if ( ctRegex.MatchString(data) ){
    // Use regex to replace channel token
    data = ctRegex.ReplaceAllString(data, "<-");
    fmt.Println("[SYNC] Rewriter - Channel Token Update");
  }

  // Return string to caller
  return data;

}
