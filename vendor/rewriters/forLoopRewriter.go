package rewriters

import (
  "regexp"
  "logger"
)

func ForLoopRewriter(data string, log logger.LevelLogger) string {

  // Pattern match for golang for loop
  // (mandatory) for(
  // (optional, parameterization) i := 0; i < 5; i++
  // (mandatory) ){
  // NOTE - This will match up to the FIRST opening bracket
  // after the loop is opened - this means function CALLS are
  // legitimate syntax but function BODIES are NOT
  flRegex := regexp.MustCompile(`\bfor\s*\(\s*(.*?)\s*\)\s*{`);

  // If a "for" loop is matched, unwrap and pad one space on left/right
  if( flRegex.MatchString(data) ){
    matches := flRegex.FindAllStringIndex(data, -1);
    log.Logf(5, "[SYNC] Rewriter - Unwrapping %d for loop(s)...\n", len(matches));
    // Use regex to unwrap the parameter capture groups
    data = flRegex.ReplaceAllString(data, "for $1 {");;
  }

  // Send string back to caller
  return data;

}
