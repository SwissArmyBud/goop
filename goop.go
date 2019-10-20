package main

import (
  "fmt"
  str "strings"
  "os"
  "path/filepath"
  file "io/ioutil"
  strUtil "strconv"
  "rewriters"
)

var itoa = strUtil.Itoa

func main() {
  // Greet
  fmt.Println("\n[ENGINE] Starting...");

  // Get
  var input string;
  if(len(os.Args) < 2){
    name, err := os.Getwd();
    if err != nil { panic(err); }
    input = name;
  } else {
    input = os.Args[1];
  }
  fmt.Println("[ENGINE] Searching directory: " + input);

  // Work
  pathWalker(input);

  // Wave
  fmt.Println("[ENGINE] Finished...\n");
}

func pathWalker(path string){

  // Setup an array for all the files that need processing
  var parsedFiles []string;

  // Walk the filepath and gather files
  err := filepath.Walk(path, func(fname string, finfo os.FileInfo, err error) error {
    // Error out on error (capture??)
    if err != nil {
      fmt.Println("[ERROR] -> Filepath error for: " + fname);
      return err;
    }
    // Gather files with matching extensions
    if ( !finfo.IsDir() ) && ( fname[len(fname)-3:] == "goo" ) {
      parsedFiles = append(parsedFiles, fname);
    }
    // Return nil error
    return nil;
  });
  if err != nil { panic(err) }
  fmt.Println("[ENGINE] Found " + itoa(len(parsedFiles)) + " goo file(s)...")

  // Create a channel, fill it, and then drain with workers
  channel := make(chan bool);
  // Fill channel to begin processing
  for i:=0; i<len(parsedFiles); i++ {
    // Spin off concurrent tasks
    go walkWorker(parsedFiles[i], channel);
  }
  // Drain channel to ensure concurrent processing is finished
  for i:=0; i<len(parsedFiles); i++ {
    <-channel;
  }

}

// ASYNC WORKER
func walkWorker(path string, channel chan bool){

  // Get text from file path
  fileString, err := file.ReadFile(path);
  if err != nil { panic(err) }
  lines := str.Split(string(fileString), "\n");
  fmt.Println("[ASYNC] Dispatcher - Scanning " + itoa(len(lines)) + " lines in: " + path);

  // Scan all lines in the file
  // NOTE - This scan method mandates single-line goop function definitions
  // TODO - Implement a marking system for lines, since line scanning is sync inside worker
  rewriteFunctions := [](func(string)string){
    rewriters.MethodRewriter,
    rewriters.WhileLoopRewriter,
    rewriters.ForLoopRewriter,
    rewriters.ChannelTokenRewriter,
  };
  for i:=0;i<len(lines);i++{
    for _, rewriter := range(rewriteFunctions) {
      lines[i] = rewriter(lines[i]);
    }
  }

  // Write back to file with new extension
  path = path[:len(path)-3] + "go";
  fileString = []byte(str.Join(lines, "\n"));
  err = file.WriteFile(path, fileString, 0644);
  if err != nil { panic(err) }

  // Tick channel to signify completion
  channel<-true;

}
