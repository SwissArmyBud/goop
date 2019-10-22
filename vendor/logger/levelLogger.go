package logger

import (
  "fmt"
)

type LevelLogger struct {
  Level int;
}

func (this *LevelLogger) GetLevel() int { return this.Level; }
func (this *LevelLogger) SetLevel( l int ){ this.Level = l; }
func (this *LevelLogger) Log( l int, args ...interface{}){
  if( l <= this.Level ){
    fmt.Print(args...);
  }
}
func (this *LevelLogger) Logln( l int, args ...interface{}){
  if( l <= this.Level ){
    fmt.Println(args...);
  }
}
func (this *LevelLogger) Logf( l int, s string, args ...interface{}){
  if( l <= this.Level ){
    fmt.Printf(s, args...);
  }
}
