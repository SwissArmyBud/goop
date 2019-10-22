package rewriters

import (
  "regexp"
  "logger"
)

var channelToken string = "<<";

func ChannelTokenRewriter(data string, log logger.LevelLogger) string{

  // Pattern match for channel token
  // (mandatory) <-
  ctRegex := regexp.MustCompile(channelToken);

  if ( ctRegex.MatchString(data) ){
    // Use regex to replace channel token
    data = ctRegex.ReplaceAllString(data, "<-");
    log.Logln(5, "[SYNC] Rewriter - Channel Token Update");
  }

  // Return string to caller
  return data;

}
