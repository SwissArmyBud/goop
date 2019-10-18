package vertex

import (
	"math"
    "fmt"
)

type Vertex struct{
	X, Y float64
}

func ( this *Vertex ) Abs() float64 {
	return math.Sqrt( this.X*this.X + this.Y*this.Y );
}

func ( this *Vertex ) Scale(f float64) {
	this.X = this.X * f
	this.Y = this.Y * f
}

func ( this *Vertex ) EightAbs() float64 {

  // Named function
  two := func () float64 {
    return 2;
  }
  
  // Assigned closure
  four := func () float64 {
    return 4;
  }
  
  // Anonymous closure
  return func () float64 {
    return this.Abs() * two() * four();
  }();
  
}

// Default void return type
func ( this *Vertex ) Hello() {
    fmt.Println("Hello from vertex:", this);
}