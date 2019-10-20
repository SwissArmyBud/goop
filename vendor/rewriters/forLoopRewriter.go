package rewriters

import (
  "regexp"
  "fmt"
)

func ForLoopRewriter (data string) string {

  // Pattern match for golang for loop
  // (mandatory) for(
  // (optional, parameterization) i := 0; i < 5; i++
  // (mandatory) ){
  // NOTE - This will match up to the FIRST opening bracket
  // after the loop is opened - this means function CALLS are
  // legitimate syntax but function BODIES are NOT
  flRegex := regexp.MustCompile(`\bfor\s*\(\s*(.*?)\s*\)\s*{`);

  if(flRegex.MatchString(data)){
    matches := flRegex.FindAllStringIndex(data, -1);
    fmt.Printf("[SYNC] Worker - Unwrapping %d for loop(s)...\n", len(matches));
    // Use regex to unwrap the parameter capture groups
    data = flRegex.ReplaceAllString(data, "for $1 {")
  }

  // Do something
  return data;

}
