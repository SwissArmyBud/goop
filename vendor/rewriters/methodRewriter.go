package rewriters

import (
  str "strings"
  "regexp"
  "logger"
)

var methodToken string = "::";

func MethodRewriter(data string, log logger.LevelLogger) string {

  // Pattern match for golang return type capturing:
  // (optional) map
  // (optional, mandatory if "map") [*]
  // (optional, excluded if "map") func(*)
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
  // oktype) func(int, int) map[]func(int, int)void
  rtRegex := regexp.MustCompile(`((((map)?(\[.*?\]))?|(func\(.*?\))?\s?)?([A-Za-z][A-Za-z0-9]*)\ *)$`);

  // Pattern match for interface{} number of whitespaces and a function opening
  foRegex := regexp.MustCompile(`^ *\{`);

  // Grab all tokens, init containers and validity map
  tokens := str.Split(data, " ");
  var cToken int;
  var tokenSet = make(map[int]bool);

  // Main definition for rewriting - rewrites nodes as needed
  tokenRewriter := func(){

    // Setup the needed container variables for transpiling definitions
    var returnType, className, methodName string;
    var methodDef string;

    // Method definition token, gather and invalidate
    token := tokens[cToken];
    tokenSet[cToken] = false;
    // Class name
    className = str.Split(token, methodToken)[0];
    // Method name
    methodName = str.Split( str.Split(token, methodToken)[1], "(")[0];

    // Default return type is "void"
    returnType = "void";
    // Return type is in previous token, if it exists
    if( cToken > 0 ){
      pIdx := cToken - 1;
      // Gather previous line tokens and look for return type match
      linePrefix := str.Join(tokens[:cToken], " ");
      mIdx := rtRegex.FindStringIndex(linePrefix);
      // If match exists, extract it and handle remaining token bits
      if ( mIdx != nil ){
        // If a match starts past zero, save the shard in the previous token
        returnType = str.TrimSpace(linePrefix[mIdx[0]:]);
        tokens[ pIdx ] = linePrefix[:mIdx[0]];
        tokenSet[ pIdx ] = len(tokens[pIdx]) > 0;
        // Invalidate all the gathered tokens
        for i := 0; i < pIdx; i++ {
          tokenSet[i] = false;
        }
      } // else - No regex match, return type stays default at "void"
    }

    // Start at cToken and go looking for end of param defs
    var parameters []string;
    // If token isn't self-closing
    if( !str.Contains(token, ")") ){
      // Grab parameter section of token and push
      parameters = append(parameters, str.Split(token, "(")[1]);
      // Keep pushing tokens until parameters finished
      for idx := cToken + 1; idx < len(tokens); idx++ {
        // Check for end of parameters
        if( !str.Contains(tokens[idx], ")") ){
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
    if( len(parameters) > 0 ){
      parameters = append(parameters, splitToken[0]);
    }

    // Empty and rebuild method def
    methodDef = "func";

    if( len(className) > 0 ){
      // Format for object pointer/receiver
      methodDef += " (this *" + className + ")";
    }
    if( len(methodName) > 0 ){
      methodDef += " " + methodName;
    } else {
      // Anonymous closure
    }
    // Add parenthesis to joined parameter set, add parameters to def
    methodDef += "(" + str.Join(parameters, " ") + ")";
    // Add return type if not void
    if( returnType != "void" ) { methodDef += " " + returnType + " "; }

    // We need to open the function, do it now
    methodDef += "{";
    // Then remove function opening in next token, if it exists
    if( !str.Contains(splitToken[1], "{") ){
      tokens[cToken + 1] = foRegex.ReplaceAllString(tokens[cToken + 1], "")
    }
    // Write the new method definition to the previous slot and mark valid
    tokens[cToken] = methodDef;
    tokenSet[cToken] = true;

    // Notify function was rewriten
    if( len(methodName) == 0 ){ methodName = "<ANONYMOUS>"; }
    log.Logln(5, "[SYNC] Rewriter - Method Rewrite: " + methodName);
  }

  // Check all tokens for method def signal
  for cToken = 0; cToken < len(tokens); cToken++ {
    tokenSet[cToken] = true;
    if( str.Contains(tokens[cToken],methodToken) ){
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
