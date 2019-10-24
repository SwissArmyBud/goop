<p align="center">
    <img
      alt="Goop"
      src="goop.png"
      width="400"
    />
</p>

# Golang for Object Oriented Programmers
**Native Golang function transpiling from C++ style method declarations, and other simple syntax fixes.**

## Motivation
Golang has made several unfortunate syntax choices, all of which slow down the migration of programmers into an otherwise useful language. Goop aims to fix all of Golang's shortcomings while staying low-profile and straightforward to use. This project is a simple transpiler engine that:
 - Enables the use of C++ style method declaration syntax
 - Converts C++ style `while` loops into Golang's single-condition `for` loops
 - Unwraps C++ style `for` loop parameterizations from parenthesis
 - Changes the channel operations token to `<<`, a C++ stream operator
 - Adds the `any` keyword from C++ as a substitute for Golang's `interface{}`

## Basic Use
Goop doesn't change anything about the existing build process, and attempts to interfere as little as possible with the language's ecosystem. The engine looks for files recursively - so it simply needs to be pointed at the root of a project or called from inside the root folder. Working from dedicated `.goo` files, the engine rewrites the syntax and outputs a standard Golang `.go` file, which can then be imported/built/run etc as usual.

## Style
### Standard Function Rewriting
Go introduces a new syntax token when declaring variables, it is `:=` and is called a [short variable declaration](https://tour.golang.org/basics/10). Goop follows this syntax hinting closely, introducing `::` as the "short function declaration". Just as C++ member functions have member and class names separated by the `::` token - in Goop _all_ functions are declared with this syntax token. Keeping with standard Go function syntax, in Goop all return types are optional and default to the `void` type.

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
### Member Function Rewriting / Reference Injection
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

### `while` Loop Conversion and Unwrapping
Golang forces programmers to use single-condition `for` loops to duplicate `while` loop behavior that has been standardized in C++ for over 30 years. Goop syntax supports traditional `while` loops - and because Goop will also unwrap `for` loops, `while` loops can be wrapped in parenthesis.
```
// GOLANG
let i := 0;
for i < 5 {
  i += 2;
  fmt.Println(i);
}

// GOOP
let i := 0;
while (i < 5) {
  i += 2;
  fmt.Println(i);
}
```
### `for` Loop Unwrapping
Golang allows `if` statements to have their test condition be wrapped or unwrapped, but forces `for` loops to be unwrapped. Goop will unwrap any `for` loop parameterization, in place and without complaint.
```
// GOLANG
for i := 0; i < 5; i++ {
  fmt.Println(i);
}

// GOOP
for ( i := 0; i < 5; i++ ){
  fmt.Println(i);
}
```
### Channel Token Conversion
Golang uses `<-` as a channel operation syntax token, even though those operations are conceptually similar to the use of `<<` in C++. Goop allows the use of the `<<` token and will replace each instance it finds with the Golang equivalent.
```
// GOLANG
ch := make(chan bool);
for i := 0; i < 5; i++ {
  ch <- true;
}
for i := 0; i < 5; i++ {
  <- ch;
}

// GOOP
ch := make(chan bool);
for (i := 0; i < 5; i++) {
  ch << true;
}
for (i := 0; i < 5; i++) {
  << ch;
}
```
### `any` Type Conversion
Golang uses `interface{}` as a generic type, but maintains memory allocation and lifecycle information for objects instantiated or passed as an `interface{}` - just like the C++ `any` class. By replacing `any` tokens in source, Goop can cut down on boilerplate, increase readability, and ease transition for existing C++ developers.
```
// GOLANG
type NamedPrinter struct {
  name string;
}
func (np *NamedPrinter) Print(args ...interface{}){
  fmt.Println("Printer: " + np.name);
  fmt.Println(args...);
}

// GOOP
type NamedPrinter struct {
  name string;
}
NamedPrinter::Print(args ...any){
  fmt.Println("Printer: " + this.name);
  fmt.Println(args...);
}
```
## Limitations in Goop Syntax
The following (known) limitations exist inside the Goop ecosystem. These issues should be simple to work around, as vanilla Golang syntax is still valid Goop syntax - and any programmer well versed enough in Golang to hit an issue in Goop syntax should be able to avoid it in the first place.

**Method Declaration Rewriting**
 - Method declarations MUST be started and finished on the same line of code:
```
// ILLEGAL Function Declaration
int Vertex::MultMultX( firstMult integer,
                        secndMult integer,
                        thirdMult integer ) {
  return this.X * firstMult * secndMult * thirdMult;
}

// Legal Goop Syntax - Single Line Function Declaration
int Vertex::MultMultX( firstMult integer, secndMult integer, thirdMult integer ) {
  return this.X *
          firstMult *
          secndMult *
          thirdMult;
}
```
 - Some complicated return types may not parse correctly:
```
// ILLEGAL Declaration
func(int)(a int, b int) Vertex::SomeFunc(){ return getFuncPtr(); }

// Legal Goop Syntax - Use Return Type Coercion
type newType func(int)(a int, b int)
newType Vertex::SomeFunc(){ return getFuncPtr(); }
```

**For Loop Unwrapping**
 - The matching system to find `for` loops will match a `for` against the first opening bracket (`{`) it sees, which means function CALLS are legal syntax inside Goop `for` loops, but function BODIES are not.
```
 // ILLEGAL for Loop
 for( i:=0; i < int ::(){ return 5; }(); i++) { fmt.Println(i); }

 // Legal Goop Syntax - Function body outside loop parameters
 f := int ::(){ return 5; }
 for( i:=0; i < f(); i++) { fmt.Println(i); }
```

**Channel Token Conversion**
 - ANY token in source matching `<<` will be transpiled to `<-` INCLUDING logging strings, regex strings, and code comments. This _could_ break doc gen or other build systems, and _will_ break bit-shifting operations. Areas that break can be ignored by using the engine override commands.

## Engine Overriding
Because being able to ignore problems is often easier than trying to fix them, Goop includes a comment parsing system that can start and stop the engine when it causes trouble, as well as delete or comment other code lines. Use the following comment formats to control the Goop engine _for each file_:

| action | syntax |
| :---: | --- |
| **Start Transpiling (Default)** | `// [GOOP][START]` |
| **Stop Transpiling** | `// [GOOP][STOP]` |
| **Skip (n) Line** | `// [GOOP][SKIP][n]` |
| **Start Deleting Lines** | `// [GOOP][DELETE][START]` |
| **Stop Deleting Lines** | `// [GOOP][DELETE][STOP]` |
| **Delete (n) Lines** | `// [GOOP][DELETE][n]` |
| **Start Commenting Lines** | `// [GOOP][COMMENT][START]` |
| **Stop Commenting Lines** | `// [GOOP][COMMENT][STOP]` |
| **Comment (n) Lines** | `// [GOOP][COMMENT][n]` |

## Example / Demo
The repo contains a complete demo based on the [Tour of Go](https://tour.golang.org/methods/4), once Goop is built call it and point it at the example directory `./goop.exe ./example` to transpile the example project. Once transpiling is complete, build the example as normal and run it. Goop itself is written in Goop, as all good compilers should be - and a full copy is included alongside the Golang files they produce.
