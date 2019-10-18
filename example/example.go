package main

import (
  vtx "vertex"
  "fmt"
)

func main(){
  v := vtx.Vertex{3, 4};
  fmt.Println(v.Abs());
  fmt.Println(v.EightAbs());
  v.Hello();
}
