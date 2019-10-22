package rewriters

import (
  "regexp"
  "fmt"
)

func WhileLoopRewriter(data string) string{

  // Pattern match for golang for loop
  // (mandatory) for
  flRegex := regexp.MustCompile(`\bwhile\b`);

  if( flRegex.MatchString(data) ){
    matches := flRegex.FindAllStringIndex(data, -1);
    fmt.Printf("[SYNC] Rewriter - Converting %d for loop(s)...\n", len(matches));
    // Use regex to unwrap the parameter capture groups
    data = flRegex.ReplaceAllString(data, "for")
  }

  // Return string to caller
  return data;

}
