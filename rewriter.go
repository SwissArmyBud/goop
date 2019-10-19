package main

import (
  str "strings"
  "regexp"
  "fmt"
)

func Rewriter (data string) string {

  // Pattern match for golang return type capturing:
  // (optional) map
  // (optional, unless "map") []
  // (mandatory) *single alpha* + *unlimited alpha or numeric* @ END OF LINE
  // https://regexr.com/
  // Tested against:
  // oktype
  // oktype(map[]float64
  // oktype(map[ ]float64
  // oktype(map[integer]float64
  // oktype(float64
  // oktype([]float64
  // oktype([10]float64
  // oktype{[]float64
  // oktype}[]float64
  // oktype};[]float64
  rtRegex := regexp.MustCompile(`((map\[.*?\])?(\[.*?\])?[A-Za-z][A-Za-z0-9]*)+$`);


  // Pattern match for any number of whitespaces and a function opening
  foRegex := regexp.MustCompile(`(\s*)\{`);
  // Use regex to replace extra whitespace groups,
  // before function opening bracket (style)
  data = foRegex.ReplaceAllString(data, "{");

  // Grab all tokens, init containers and validity map
  tokens := str.Split(data, " ");
  var cToken int;
  var tokenSet = make(map[int]bool);

  // Main definition for rewriting - traverses nodes as needed
  tokenRewriter := func (){

    // Setup the needed container variables for transpiling definitions
    var returnType, className, methodName string;
    var methodDef []string;

    // Method definition token
    token := tokens[cToken];
    // Class name
    className = str.Split(token, "::")[0];
    // Method name
    methodName = str.Split( str.Split(token, "::")[1], "(")[0];

    // Default return type is "void"
    returnType = "void";
    // Return type is in previous token, if it exists
    if(cToken > 0){
      pIdx := cToken - 1;
      // Look for return type match in token
      idx := rtRegex.FindIndex([]byte(tokens[ pIdx ]));
      // If match exists, extract it and handle remaining token bits
      if (idx != nil) {
        if idx[0] > 0 {
          // If a match starts past zero, save the shard in the previous token
          returnType = tokens[ pIdx ][ idx[0] : ]
          tokens[ pIdx ] = tokens[ pIdx ][ : idx[0] ]
        } else {
          // Otherwise keep the value and invalidate the previous token
          returnType = tokens[ pIdx ];
          tokenSet[ pIdx ] = false;
        }
      } // else - No regex match, return type stays default at "void"
    }
    // First token, no return type, but token is now invalid
    tokenSet[ cToken ] = false;

    // Start at cToken and go looking for end of param defs
    var parameters []string;
    // If token isn't self-closing
    if(!str.Contains(token, ")")){
      // Grab parameter section of token and push
      parameters = append(parameters, str.Split(token, "(")[1]);
      // Keep pushing tokens until parameters finished
      for idx := cToken + 1; idx < len(tokens); idx++ {
        // Check for end of parameters
        if(!str.Contains(tokens[idx], ")")){
          // Gather value and invalidate token
          parameters = append(parameters ,tokens[idx]);
          tokenSet[idx] = false;
        } else {
          // Update tracked token position
          cToken = idx;
          // Break parameter fill loop
          break;
        }
      }
    }

    // Split up the last parameter token
    splitToken := str.Split(tokens[cToken], ")");
    // Only need to append to parameters when they exist
    if(len(parameters) > 0){
      parameters = append(parameters, splitToken[0]);
    }
    // Leave function opening in current token
    tokens[cToken] = splitToken[1];

    // Empty and rebuild method def
    methodDef = []string{
        "func",
    };

    if(len(className) > 0){
      // Format for object pointer/receiver
      methodDef = append(methodDef, []string{
        " (",
        "this",
        " *" + className,
        ")",
      }...)
    }
    if(len(methodName) > 0){
      methodDef = append(methodDef, " " + methodName);
    } else {
      // Anonymous closure
    }
    methodDef = append(methodDef, "(" + str.Join(parameters, " ") + ")")
    // Add return type if not void
    if(returnType != "void") { methodDef = append(methodDef, " " + returnType + " "); }
    // If we stored a shard in the current slot, add it to the new line
    methodDef = append(methodDef, tokens[cToken])
    // Write the new method definition to the previous slot and mark valid
    tokens[ cToken ] = str.Join(methodDef, "");
    tokenSet[ cToken ] = true;

    // Notify function was rewriten
    if(len(methodName) == 0){ methodName = "<ANONYMOUS>"; }
    fmt.Println("[SYNC] Worker - Method Rewrite: " + methodName);
  }

  // Check all tokens for method def signal
  for cToken = 0; cToken < len(tokens); cToken++ {
    tokenSet[cToken] = true;
    if(str.Contains(tokens[cToken],"::")){
      // Process token on signal
      tokenRewriter();
    }
  }
  // Count and condense valid tokens
  vTokens := 0;
  for cToken = 0; cToken < len(tokens); cToken++ {
    if tokenSet[cToken] {
      tokens[vTokens] = tokens[cToken];
      vTokens++;
    }
  }

  // Return string made from a slice with valid tokens
  return str.Join(tokens[:vTokens], " ");

}
