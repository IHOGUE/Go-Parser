package main

import (
    "fmt"
    "io/ioutil"
    "os"
)

func check(e error) {
    if e != nil {
        panic(e)
    }
}

type token struct{
    name string
    lexeme string
}
type point struct{
    x string
    y string
    name string
}
type test struct{
    option string
    pointList []*point
    id int
}

func main() {
   
    args := os.Args[1:]
    
    tokenList := make([]token, 0)
   
    terminals := map[byte]string{'=': "ASSIGN",  ';': "SEMICOLON", ',': "COMMA", '.': "PERIOD", '(': "LPAREN", ')': "RPAREN"} 
    keyWords := map[string]string{"point": "POINT", "test": "TEST", "square": "SQUARE", "triangle": "TRIANGLE"}
   
    content, err := ioutil.ReadFile(args[0])
    check(err)

    var input string 
    //------------------------------------------------LEXER-------------------------------------------------------------------------
    for i := 0; i <= len(content) - 1; i++{
        input = ""
        if terminal, ok := terminals[content[i]]; ok {
            token := token{name: terminal, lexeme: string(content[i])} 
            tokenList = append(tokenList, token)
        } else if content[i] > 47 && content[i] < 58{   // if digit, 0 - 9
            token := token{name: "NUM", lexeme: string(content[i])} 
            tokenList = append(tokenList, token)
        } else if content[i] > 96 && content[i] < 123 {   // if lower case letter
            //loop until offending byte or terminal
            input = input + string(content[i])
            for (i + 1 < len(content) && content[i+1] > 96 && content[i + 1] < 123){
                input = input + string(content[i+1])
                i++
            }            
            if word, ok := keyWords[input]; ok {  // this is for adding keywords like "square"
                token := token{name: word, lexeme: input}
                tokenList = append(tokenList, token)
            } else{ // right now, an ID will be created up to the offending byte, so abc? will make token abc, and then ? won't be recognized
                token := token{name: "ID", lexeme: input}
                tokenList = append(tokenList, token)
            }         
        } else if content[i] == 32 || content[i] == 10 || content[i] == 13 { // catch spaces, line feeds, and carriage returns
            continue
        } else {
            fmt.Println("Lexical error", string(content[i]), "not recognized")
            return
        }
    }
    // ------------------------------------------PARSER--------------------------------------------------------------------------------
    var goodParse bool
    for i:= 0; i < len(tokenList); i++{
        match, k, ex := STMT(tokenList, i)
        i = k
        goodParse = match
        if !match{
            fmt.Println("Syntax error", tokenList[i].lexeme, "found", ex, "expected")
            return
        }
    }
    //------------------------------------------------POINT PARSER-----------------------------------------------------------------------
    testList := make([]*test, 0)
    points := map[string]*point{}
    var k int

    for i := 0; i < len(tokenList); i++{
         if (tokenList[i].name == "TEST"){
             pointList := make([]*point,0)
             temp := new(test)
             i = i + 2
             temp.option = tokenList[i].lexeme 
             temp.id = i
             i = i + 2
             for k = i; k < len(tokenList); k = k + 2{
                 if match, ok := points[tokenList[k].lexeme]; ok{
                     pointList = append(pointList, match )
                 } else if tokenList[k].name == "RPAREN"{   // breaks early when test stmt is finished
                    break;
                 } else if tokenList[k].name == "ID" && !ok{ // catch points with no definition
                     temp := new(point)
                     temp.name = tokenList[k].lexeme
                     temp.x = ""
                     temp.y = ""
                     pointList = append(pointList, temp)
                 }
             }
             temp.pointList = pointList
             testList = append(testList, temp)
         }
         if (tokenList[i].name == "ID" && tokenList[i+1].name == "ASSIGN"){
             temp := new(point)
             temp.name = tokenList[i].lexeme
             i = i + 4
             temp.x = tokenList[i].lexeme
             i = i + 2
             temp.y = tokenList[i].lexeme
             points[temp.name] = temp
         }
     }
         // -------------------------------------------Code Generation----------------------------------------------
    output := ""
    if  len(args) > 1 && args[1] == "-s" && goodParse{
        fmt.Println("; Processing input file", args[0], "\n; Lexical and Syntax analysis passed\n; Generating Scheme Code")
        if (len(points) == 0){
            fmt.Println("\nError: No points defined")
            return
        }
        if (len(testList) == 0){
            fmt.Println("\nError: No tests defined")
            return
        }
        for i:= range testList{
            output = "(process-" + testList[i].option
            for k := range testList[i].pointList{
                if Empty(testList[i].pointList[k], &output){
                    fmt.Println(output)
                    return
                }
                output += " (make-point " + testList[i].pointList[k].x + " " + testList[i].pointList[k].y + ")"   
            }
            output += ")."
            fmt.Println(output)
        }
        return
    }
    if len(args) > 1 && args[1] == "-p" && goodParse {
        fmt.Println("/* processing input file", args[0], "\n   Lexical and Syntax analysis passed\n   Generating Prolog Code */\n")
        if (len(points) == 0){
            fmt.Println("Error: no points defined")
            return
        }
        if (len(testList) == 0){
            fmt.Println("Error: no tests defined")
            return
        }
        triangleTests := [10]string{"line", "triangle", "vertical", "horizontal", "equilateral", "isosceles", "right", "scalene", "acute", "obtuse"}
        for i := range testList{
            if (testList[i].option == "square"){
                output = "/* Processing test(square"
                for k := range testList[i].pointList{
                    output += ", " + testList[i].pointList[k].name
                }    
                output += ") */"
                fmt.Println(output)
                if Empty(testList[i].pointList[0], &output){ // have to place an extra here, to catch the first point in the test
                    fmt.Println(output)
                    return
                }
                output = "query(" + testList[i].option
                output += "(point2d("+ testList[i].pointList[0].x + "," + testList[i].pointList[0].y + ")"
                for j := 1; j < len(testList[i].pointList); j++{
                    if Empty(testList[i].pointList[j], &output){
                        fmt.Println(output)
                        return
                    }
                    output += ",(point2d("+ testList[i].pointList[j].x + "," + testList[i].pointList[j].y + ")"
                }
                output += "))."
                fmt.Println(output, "\n")
            }   else if (testList[i].option == "triangle"){
                    output = "/* Processing test(triangle"
                    for k := range testList[i].pointList{
                        output += ", " + testList[i].pointList[k].name 
                    }   
                    output += ") */" 
                    fmt.Println(output)
                    output = ""    
                    for j := range triangleTests{
                        if Empty(testList[i].pointList[0], &output){ // have to place an extra here, to catch the first point in the test
                            fmt.Println(output)
                            return
                        }
                        output += "query(" + triangleTests[j]
                        output += "(point2d(" + testList[i].pointList[0].x + "," + testList[i].pointList[0].y + ")"
                        for l := 1; l < len(testList[i].pointList); l++{
                            if Empty(testList[i].pointList[l], &output){
                                fmt.Println(output)
                                return
                            }
                            output += ", point2d(" + testList[i].pointList[l].x + "," + testList[i].pointList[l].y + ")"
                        }
                        output += ")).\n"
                    }
                    fmt.Println(output)
                }
        }
        output = "/* Query Processing */\nwriteln(T) :- write(T), nl.\nmain:- forall(query(Q), Q-> (writeln('yes')) ; (writeln('no))),\n      halt."
        fmt.Println(output)
        return
    }
}
func Empty(Point *point, output *string) bool{
    if Point.x == "" && Point.y == ""{
        *output = "\nError: point " + Point.name + " not defined"
        return true;
    }
    return false
}
//--------------------------------------------------Function List for Syntax Analyzer------------------------------------------------------------
func STMT(tokenList []token, i int) (bool, int, string){
    if tokenList[i].name == "ID"{
        return POINT_DEF(tokenList, i)
    } 
    if tokenList[i].name == "TEST"{
        return TEST(tokenList, i)
    }
    return false, i, "identifier or test"
}
func TEST(tokenList []token, i int) (bool,int,string){
    i++
    return LPAREN(tokenList, i, 1) // sends parameter 1 so that code knows that I want a test statement
}
func POINT_DEF(tokenList []token, i int) (bool, int, string){
    return ID(tokenList, i, 0) // sends parameter 0 so that code knows that I want an ID assignment statement
}
func POINT_LIST(tokenList []token, i int) (bool, int, string){
    return ID(tokenList, i, 1)
}
func ID(tokenList []token, i int, stmt int) (bool, int, string){
    if ("ID" == tokenList[i].name){
        i++
        if (stmt == 0){
            return ASSIGN(tokenList, i)
        }
        if (stmt == 1){
            return COMMA(tokenList, i, 1)
        }
    }   
    return false, i, "identifier"
}
func ASSIGN(tokenList []token, i int) (bool, int, string){
    if ("ASSIGN" == tokenList[i].name){
        i++
        return POINT(tokenList, i)
    }
    return false, i, "="
}
func POINT(tokenList []token, i int) (bool, int, string){
    if ("POINT" == tokenList[i].name){
        i++
        return LPAREN(tokenList, i, 0)
    }
    return false, i, "point"
}
func OPTION(tokenList []token, i int) (bool, int, string){
    if (tokenList[i].name == "SQUARE" || tokenList[i].name == "TRIANGLE"){
        i++
        return COMMA(tokenList, i, 1)
    }
    return false, i, "sqaure or triangle"
}
func LPAREN(tokenList []token, i int, stmt int) (bool, int, string){
    if ("LPAREN" == tokenList[i].name){
        i++
        if (stmt == 1){
            return OPTION(tokenList, i)
        }
        return NUM(tokenList, i)
    }
    return false, i, "("
}
func NUM(tokenList []token, i int) (bool, int, string){
    if (tokenList[i-1].name == "COMMA" && "NUM" == tokenList[i].name){
        i++
        return RPAREN(tokenList, i, 0)
    } else if ("NUM" == tokenList[i].name){
        i++
        return COMMA(tokenList, i, 0)
    }
    return false, i, "number"
}
func COMMA(tokenList []token, i int, stmt int) (bool, int, string){
    if ("COMMA" == tokenList[i].name){
        i++
        if (stmt == 0){
            return NUM(tokenList, i)
        }
        return POINT_LIST(tokenList, i)
    }
    if (stmt == 1){
        return RPAREN(tokenList, i, 1)
    }
    return false, i, ","
}
func RPAREN(tokenList []token, i int, stmt int) (bool, int, string){
    if ("RPAREN" == tokenList[i].name){
        i++
        return END(tokenList, i)
    }
    if (stmt == 1){
        return false, i, ","
    }
    return false, i, ")"
}
func END(tokenList []token, i int) (bool, int, string){
    if ("SEMICOLON" != tokenList[i].name && i + 1 < len(tokenList)){
        return false, i, ";"
    }
    if ("PERIOD" != tokenList[i].name && i + 1 == len(tokenList)){    
        return false, i, "."
    }
    return true, i, tokenList[i].name
}