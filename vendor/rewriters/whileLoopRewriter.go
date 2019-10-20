package rewriters

import (
  "regexp"
  "fmt"
)

func WhileLoopRewriter (data string) string {

  // Pattern match for golang while loop
  // (mandatory) while
  flRegex := regexp.MustCompile(`\bwhile\b`);

  if(flRegex.MatchString(data)){
    matches := flRegex.FindAllStringIndex(data, -1);
    fmt.Printf("[SYNC] Worker - Converting %d while loop(s)...\n", len(matches));
    // Use regex to unwrap the parameter capture groups
    data = flRegex.ReplaceAllString(data, "for")
  }

  // Do something
  return data;

}
