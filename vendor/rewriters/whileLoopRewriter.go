package rewriters

import (
  "regexp"
  "logger"
)
var whileKeywordToken string = "while";
var newWhileKeyword string = "for";

func WhileLoopRewriter(data string, log logger.LevelLogger) string {

  // Pattern match for keyword token
  keywordRegex := regexp.MustCompile(`\b` + whileKeywordToken + `\b`);

  if( keywordRegex.MatchString(data) ){
    log.Logln(5, "[SYNC] Rewriter - Converting while loop...");
    // Use regex to replace keyword
    data = keywordRegex.ReplaceAllString(data, newWhileKeyword);
  }

  // Return string to caller
  return data;

}
