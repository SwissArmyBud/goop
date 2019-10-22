package main

import (
  str "strings"
  "os"
  "path/filepath"
  file "io/ioutil"
  strUtil "strconv"
  "rewriters"
  "regexp"
  "reflect"
  "logger"
)

// Configuration options
var KEEP_ENGINE_CMD_IN_OUTPUT bool = false;
var MAX_LOGGING_LEVEL int = 5;

// Condensed function pointers
var itoa = strUtil.Itoa
var atoi = strUtil.Atoi
var areEqual = reflect.DeepEqual
func regexReplaceCleanSplit(re *regexp.Regexp, in string, pt string) []string {
  return str.Split(str.TrimSpace(re.ReplaceAllString(in, pt)), " ");
}
func Panic(er error){ 
  if( er != nil ) { panic(er); }
}
var log = logger.LevelLogger{ MAX_LOGGING_LEVEL };


func main(){ 
  // Greet
  log.Logln(1, "\n[ENGINE] Starting...");

  // Get
  var input string;
  if( len(os.Args) < 2 ){
    name, err := os.Getwd();
    Panic(err);
    input = name;
  } else {
    input = os.Args[1];
  }
  log.Logln(1, "[ENGINE] Searching directory: " + input);

  // Work
  pathWalker(input);

  // Wave
  log.Logln(1, "[ENGINE] Finished...\n");
}

func pathWalker(path string){ 

  // Setup an array for all the files that need processing
  var parsedFiles []string;

  // Walk the filepath and gather files
  err := filepath.Walk(path, func(fname string, finfo os.FileInfo, err error) error { 
    // Error out on error (capture??)
    Panic(err);
    // Gather files with matching extensions
    if ( !finfo.IsDir() ) && ( fname[len(fname)-3:] == "goo" ) {
      parsedFiles = append(parsedFiles, fname);
    }
    // Return nil error
    return nil;
  });
  Panic(err);
  log.Logln(2, "[ENGINE] Found " + itoa(len(parsedFiles)) + " goo file(s)...");

  // Create a channel, fill it, and then drain with reports from workers
  channel := make(chan int);
  tLines := 0;

  // Fill channel to begin processing
  log.Logln(2, "[ENGINE] Dispatching " + itoa(len(parsedFiles)) + " rewrite workers...");
  for i:=0; i < len(parsedFiles); i++ {
    // Spin off concurrent tasks
    go walkWorker(parsedFiles[i], channel);
  }
  // Drain channel to ensure concurrent processing is finished
  for i:=0; i < len(parsedFiles); i++ {
    tLines += <-channel;
  }

  // Report reports
  log.Logf(2, "[ENGINE] Workers report scanning of %d lines...\n", tLines);

}

// ASYNC WORKER
func walkWorker(path string, channel chan int){

  // Get text from file path
  fileString, err := file.ReadFile(path);
  Panic(err);
  lines := str.Split(string(fileString), "\n");
  log.Logln(3, "[SYNC] Worker - Scanning " + itoa(len(lines)) + " lines in: " + path);

  // Compile the engine command regex pattern
  ecRegex := regexp.MustCompile(` *\/\/ *\[GOOP\] *\[([A-Z]*)\] *(\[(([0-9]*)|([A-Z]*))\])?`);
  // Compile the whitespace regex pattern
  wsRegex := regexp.MustCompile(`^( *)(.*)`);

  // Setup a list of transform functions to call
  rewriteFunctions := [](func(string, logger.LevelLogger)string){
    rewriters.MethodRewriter,
    rewriters.WhileLoopRewriter,
    rewriters.ForLoopRewriter,
    rewriters.ChannelTokenRewriter,
    rewriters.AnyTokenRewriter,
  };

  // Engine processing state, and line validity map for tracking
  processing := true;
  var lineSet = make(map[int]bool);

  // Scan all lines in the file
  for i:=0; i < len(lines); i++ {

    // Default include line in output
    lineSet[i] = true;

    // If the line is an engine command, set state accordingly
    if( ecRegex.MatchString(lines[i]) ){

      // Get the engine command from a regex groups, trim whitespaces, split
      command := regexReplaceCleanSplit(ecRegex, lines[i], "$1 $3");
      // If we aren't processing, ignore anything except a start command
      if( !processing && command[0] != "START" ) { continue; }

      // We keep or delete the command line itself and point to the next line
      lineSet[i] = KEEP_ENGINE_CMD_IN_OUTPUT;
      i++;

      // Set engine parameters based on parsed command
      switch( command[0] ){
        case "START": {
          log.Logf(4, "[SYNC] Worker - Starting transpile engine on line %d...\n", i)
          processing = true;
        }
        case "STOP": {
          log.Logf(4, "[SYNC] Worker - Stopping transpile engine on line %d...\n", i)
          processing = false;
        }
        case "SKIP": {
          // Get the int value of the command parameter
          count, err := atoi(command[1]);
          Panic(err);
          log.Logf(4, "[SYNC] Worker - Ignoring " + itoa(count) + " lines at %d...\n", i);
          // For each skip, include line in output and move to next line
          for l := 0; l < count; l++ {
            lineSet[i] = true;
            i++;
          }
        }
        case "DELETE": {
          // Try and get an integer value from the command parameter
          count, err := atoi(command[1]);
          // Failure means we need to parse as START token
          if( err != nil ){
            if( command[1] == "START" ){
              deleting := true;
              log.Logf(4, "[SYNC] Worker - Deleting from line %d...\n", i);
              for deleting {
                // While we are deleting lines, keep checking for STOP token
                if( ecRegex.MatchString(lines[i]) ){
                  cmd := regexReplaceCleanSplit(ecRegex, lines[i], "$1 $3");
                  if( areEqual(cmd, []string{"DELETE", "STOP"}) ){
                    deleting = false;
                  }
                }
                // Delete is automatic when no 'lineSet' token is generated
                // Point to next line
                i++;
              }
              log.Logf(4, "[SYNC] Worker - Deleted to line %d...\n", i);
            } else {
              panic("Bad parameter passed to engine's DELETE function!")
            }
          } else {
            // If an int is provided, pass that many lines without validating
            log.Logf(4, "[SYNC] Worker - Deleting " + itoa(count) + " lines at %d...\n", i);
            for l := 0; l < count; l++ {
              // Delete is automatic when no 'lineSet' token is generated
              // Point to next line
              i++;
            }
          }
        }
        case "COMMENT": {
          // Try and get an integer value from the command parameter
          count, err := atoi(str.TrimSpace(command[1]));
          // Failure means we need to parse as START token
          if( err != nil ){
            if( command[1] == "START" ){
              commenting := true;
              for commenting {
                // While we are deleting lines, keep checking for STOP token
                if( ecRegex.MatchString(lines[i]) ){
                  cmd := regexReplaceCleanSplit(ecRegex, lines[i], "$1 $3");
                  if( areEqual(cmd, []string{"COMMENT", "STOP"}) ){
                    commenting = false;
                  }
                }
                // Commenting requires rewriting the line and marking valid
                if ( commenting ){
                  lines[i] = wsRegex.ReplaceAllString(lines[i], "$1// $2");
                  lineSet[i] = true;
                }
                // Point to the next line
                i++;
              }
            } else {
              panic("Bad parameter passed to engine's COMMENT function!")
            }
          } else {
            log.Logf(4, "[SYNC] Worker - Commenting " + itoa(count) + " lines at %d...\n", i);
            for l := 0; l < count; l++ {
              // Commenting requires rewriting the line and marking valid
              lineSet[i] = true;
              lines[i] = wsRegex.ReplaceAllString(lines[i], "$1// $2");
              i++;
            }
          }
        }
        default: {
          log.Logln(4, "[SYNC] Worker - Unknown engine command: " + command[0])
        }
      }
      // Recycle the for loop, decrement i to avoid offset change on iteration
      i--;
      continue;
    }
    if ( processing ) {
      for _, rewriter := range(rewriteFunctions) {
        lines[i] = rewriter(lines[i], log);
      }
    }
  }

  // Count and condense valid lines
  vTokens := 0;
  for i := 0; i < len(lines); i++ {
    if ( lineSet[i] ) {
      lines[vTokens] = lines[i];
      vTokens++;
    }
  }

  // Write back to file with new extension
  path = path[:len(path)-3] + "go";
  fileBytes := []byte(str.Join(lines[:vTokens], "\n"));
  err = file.WriteFile(path, fileBytes, 0644);
  Panic(err);

  // Tick channel to signify completion
  channel <- len(lines);

}
