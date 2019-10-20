<p align="center">
    <img
      alt="Goop"
      src="goop.png"
      width="400"
    />
</p>

# Golang for Object Oriented Programmers
**Native Golang transpiling from C++ style method declarations, plus other syntax fixes.**

## Motivation
Golang has made several unfortunate syntax choices, none of which contribute to what is an otherwise useful language, and Goop fixes all of it while staying low-profile and straightforward to use. This project is a simple transpiling engine that allows:
 - The writing of C++ style method definitions which are then transpiled into native/legal Golang syntax
 - Converts `while` loops into single-condition `for` loops
 - Unwraps `for` loop parameterizations from parenthesis
 - Changes the channel token from `<-` to `<<`, a pre-existing C++ token for feeding pipes.

## Overview
Goop doesn't change anything about the existing build process, and attempts to interfere as little as possible with the language's ecosystem. The engine looks for files recursively - so it simply needs to be pointed at the root of a project or called from inside the root folder. Working from dedicated `.goo` files, the engine rewrites function definitions and then outputs a standard golang `.go` file, which can be imported/built/run etc as usual.

## Style
#### While-Loop Conversion and Unwrapping
Golang forces programmers to use single-condition `for` loops to duplicate `while` loop behaviors that have been standardized for over 30 years. Goop allows for `while` loops and because Goop will also unwrap `for` loops, `while` loops can be wrapped in Goop syntax.
```
// GOLANG LEGAL
let i := 0;
for i++ < 5 {
  fmt.Println(i);
}

// GOOP LEGAL
let i := 0;
while (i++ < 5) {
  fmt.Println(i);
}
```
#### For-Loop Unwrapping
Golang allows `if` statements to have their test condition be wrapped or unwrapped, but forces `for` loops to be unwrapped. Goop will unwrap any `for` loop parameterization in place without complaint.
```
// GOLANG LEGAL
for i := 0; i < 5; i++ {
  fmt.Println(i);
}

// GOOP LEGAL
for ( i := 0; i < 5; i++ ){
  fmt.Println(i);
}
```
#### Channel Token
Golang does not include `<<` as a syntax token, even though the usage of Golang's `<-` token is conceptually similar to the use of `<<` in C++. Goop allows the use of the '<<' token and will replace each instance it finds.
```
// GOLANG LEGAL
ch := make(chan bool);
for i := 0; i < 5; i++ {
  ch <- true;
}
for i := 0; i < 5; i++ {
  <- ch;
}

// GOOP LEGAL
ch := make(chan bool);
for (i := 0; i < 5; i++) {
  ch << true;
}
for (i := 0; i < 5; i++) {
  << ch;
}
```
#### Standard Function Rewriting
Go introduces a new syntax token when declaring variables, it is `:=` and is called a [short variable declaration](https://tour.golang.org/basics/10). Goop follows this syntax hinting closely, introducing the "short function declaration" thanks to C++ member functions having member and class names separated by the `::` token - in Goop _all_ functions are declared with this syntax token. Keeping with standard Go function syntax, in Goop all return types are optional and default to the `void` type.
<br><br>
**Named Function**
```
// GOLANG
func Hello(){ fmt.Println("Hello World!"); }
func Double(val int) int { return 2 * val; }

// GOOP
::Hello(){ fmt.Println("Hello World!"); }
int ::Double(val int){ return 2 * val; }
```
**Named Closure**
```
// GOLANG
Hello := func(){ fmt.Println("Hello World!"); }
Double := func(val int) int { return 2 * val; }

// GOOP
Hello := ::(){ fmt.Println("Hello World!"); }
Double := int ::(val int){ return 2 * val; }
```
**Anonymous Closure**
```
// GOLANG
func (){ fmt.Println("Hello World!"); }();
func (val int) int { return 2 * val; }();

// GOOP
::(){ fmt.Println("Hello World!"); }();
int ::(val int){ return 2 * val; }();
```
#### Member Function Rewriting / Reference Injection
Goop can cut down on boilerplate for standard functions, but the real power is when writing member functions. Class member pointers are automatically injected as `this` and function bodies can then be written as if in C++.

**Member Function**
```
// Type Definition
type Vertex struct {
  X, Y float64
}

// GOLANG
func (v *Vertex) Scale(f float64) {
  v.X = v.X * f
  v.Y = v.Y * f
}
func (v Vertex) Abs() float64 {
  return math.Sqrt(v.X*v.X + v.Y*v.Y)
}

// GOOP
Vertex::Scale(f float64){
  this.X = this.X * f
  this.Y = this.Y * f
}
float64 Vertex::Abs(){
  return math.Sqrt(this.X*this.X + this.Y*this.Y);
}
```
<br>

## Limitations in Goop Syntax
The following (known) limitations exist inside the Goop ecosystem. These issues should be simple to work around, as vanilla Golang syntax is still valid Goop syntax - and any programmer well versed enough in Golang to hit an issue in Goop syntax should have no problem understanding how to avoid it in the first place.

**Method Declaration Rewriting**
 - Method declarations MUST be started and finished on the same line of code:
 ```
// ILLEGAL Declaration Spans Multiple Lines
int Vertex::MultMultX( firstMult integer,
                        secndMult integer,
                        thirdMult integer ) {
  return this.X * firstMult * secndMult * thirdMult;
}

// Legal Goop Syntax Function Declaration
int Vertex::MultMultX( firstMult integer, secndMult integer, thirdMult integer ) {
  return this.X *
          firstMult *
          secndMult *
          thirdMult;
}
```
 - Some complicated return types may not parse correctly:
```
// ILLEGAL Declaration, Return Type Breaks at Named Return
func(int)(a int, b int) Vertex::SomeFunc(){ return getFuncPtr(); }

// Legal Goop Syntax - Use Return Type Coercion
type newType func(int)(a int, b int)
newType Vertex::SomeFunc(){ return getFuncPtr(); }
```

**For Loop Unwrapping**
 - The matching system to find for loops will match a `for` against the first opening bracket (`}`) it sees, which means function CALLS are legal syntax inside `for` loops, but function BODIES are not.
 ```
 // ILLEGAL for Loop
 for( i:=0; i < ::(){ return 5; }; i++) { fmt.Println(i); }

 // Legal Goop Syntax - Function body outside loop parameters
 f := int ::(){ return 5; }
 for( i:=0; i < f(); i++) { fmt.Println(i); }
 ```

**Channel Token Conversion**
 - ANY token in source matching `<<` will be transpiled to `<-` INCLUDING logging string, regex strings, and code comments. This _could_ break doc gen or build systems.

## Example / Demo
The repo contains a complete demo based on the [Tour of Go](https://tour.golang.org/methods/4), once Goop is built call it and point it at the example directory `./goop.exe ./example` to transpile the example project. Once transpiling is complete, build the example as normal and run it.
