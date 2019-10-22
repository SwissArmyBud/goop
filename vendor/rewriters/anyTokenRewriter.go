package rewriters

import (
  "regexp"
  "logger"
)

var anyKeywordToken string = "any";
var newAnyKeyword string = "interface{}";

func AnyTokenRewriter(data string, log logger.LevelLogger) string {

  // Pattern match for keyword token
  keywordRegex := regexp.MustCompile(`\b` + anyKeywordToken + `\b`);

  if ( keywordRegex.MatchString(data) ){
    log.Logln(5, "[SYNC] Rewriter - Converting any token...");
    // Use regex to replace keyword
    data = keywordRegex.ReplaceAllString(data, newAnyKeyword);
  }

  // Return string to caller
  return data;

}
