%{
// vim: set filetype=go:
//
// go-sonos
// ========
//
// Copyright (c) 2012, Ian T. Richards <ianr@panix.com>
// All rights reserved.
//
// Redistribution and use in source and binary forms, with or without
// modification, are permitted provided that the following conditions
// are met:
//
//   * Redistributions of source code must retain the above copyright notice,
//     this list of conditions and the following disclaimer.
//   * Redistributions in binary form must reproduce the above copyright
//     notice, this list of conditions and the following disclaimer in the
//     documentation and/or other materials provided with the distribution.
//
// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS
// "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT
// LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR
// A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENTING SHALL THE COPYRIGHT
// HOLDER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL,
// SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED
// TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR
// PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF
// LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING
// NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS
// SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
//
package main

import (
	"bufio"
	"bytes"
	"io"
	"log"
	"os"
	"strconv"
	"unicode"
	"unicode/utf8"
)

type Arg struct {
	id int
	bv bool
	sv string
	rv float64
}

type Cmd struct {
	name string
	args []Arg
}

type Script struct {
	cmds []Cmd
}

var SCRIPT Script

%}

%token BOOL
%token COMMA
%token IDENT
%token LPAREN
%token NUMBER
%token RPAREN
%token SEMICOLON
%token STRING

%start script

%type<av> argument
%type<al> argument_list
%type<al> argument_list_tail
%type<bv> BOOL
%type<cv> command
%type<cl> command_list
%type<cl> command_list_tail
%type<sv> IDENT
%type<rv> NUMBER
%type<xv> script
%type<sv> STRING

%union {
	bv bool
	sv string
	rv float64
	av Arg
	al []Arg
	cv Cmd
	cl []Cmd
	xv Script
}

%%

script: command_list
	{
		$$ = Script{cmds: $1}
		execute($$)
	}

command_list :
	{
		$$ = nil
	} | command command_list_tail {
		$$ = append(make([]Cmd, 0), $1)
		$$ = append($$, $2 ...)
	}

command : IDENT LPAREN argument_list RPAREN {
	$$ = Cmd{name: $1, args: $3}
}

command_list_tail :
	{
		$$ = nil
	} | SEMICOLON command command_list_tail {
		$$ = append(make([]Cmd, 0), $2)
		$$ = append($$, $3 ...)
	}

argument_list :
	{
		$$ = nil
	}
	| argument argument_list_tail {
		$$ = append(make([]Arg, 0), $1)
		$$ = append($$, $2 ...)
	}

argument :
	IDENT {
		$$ = Arg{id: IDENT, sv: $1}
	}
	| NUMBER {
		$$ = Arg{id: NUMBER, rv: $1}
	}
	| STRING {
		$$ = Arg{id: STRING, sv: $1}
	}
	| BOOL {
		$$ = Arg{id: BOOL, bv: $1}
	}

argument_list_tail :
	{
		$$ = nil
	}
	| COMMA argument argument_list_tail {
		$$ = append(make([]Arg, 0), $2)
		$$ = append($$, $3 ...)
	}

%%

type Lexer struct {
	line []byte
	peek rune
}

func (this *Lexer) Lex(yylval *yySymType) int {
	for {
		c := this.next()
		if 0 == c {
			return 0
		} else if '(' == c {
			return LPAREN
		} else if ')' == c {
			return RPAREN
		} else if ';' == c {
			return SEMICOLON
		} else if ',' == c {
			return COMMA
		} else if unicode.IsLetter(c) {
			return this.getIdent(c, yylval)
		} else if unicode.IsDigit(c) {
			return this.getNumber(c, yylval)
		} else if '"' == c {
			return this.getString(c, yylval)
		}
	}
}

func (this *Lexer) getString(c rune, yylval *yySymType) int {
	var buffer bytes.Buffer
	escaped := false
	for {
		c = this.next()
		if escaped {
			buffer.WriteRune(c)
			escaped = false
		} else if '\\' == c {
			escaped = true
		} else if '"' == c {
			break
		} else {
			buffer.WriteRune(c)
		}
	}
	yylval.sv = buffer.String()
	return STRING
}

func (this *Lexer) getNumber(c rune, yylval *yySymType) int {
	var buffer bytes.Buffer
	buffer.WriteRune(c)
	has_point := false
	for {
		c = this.next()
		if '.' == c {
			if has_point {
				panic("malformed literal")
			} else {
				buffer.WriteRune(c)
				has_point = true
			}
		} else if unicode.IsDigit(c) {
			buffer.WriteRune(c)
		} else {
			this.peek = c
			break
		}
	}
	var err error
	s := buffer.String()
	if yylval.rv, err = strconv.ParseFloat(s, 64); nil != err {
		panic(err)
	}
	return NUMBER
}

func (this *Lexer) getIdent(c rune, yylval *yySymType) int {
	var buffer bytes.Buffer
	buffer.WriteRune(c)
	for {
		c = this.next()
		if unicode.IsLetter(c) || unicode.IsDigit(c) || '_' == c {
			buffer.WriteRune(c)
		} else {
			this.peek = c
			break
		}
	}
	var err error
	s := buffer.String()
	if yylval.bv, err = strconv.ParseBool(s); nil == err {
		return BOOL
	} else {
		yylval.sv = s
		return IDENT
	}
}

func (this *Lexer) next() rune {
	if 0 != this.peek {
		r := this.peek
		this.peek = 0
		return r
	} else if 0 == len(this.line) {
		return 0
	} else {
		c, size := utf8.DecodeRune(this.line)
		this.line = this.line[size:]
		if utf8.RuneError == c && 1 == size {
			panic("invalid utf8 character")
		}
		return c
	}
}

func (this *Lexer) Error(s string) {
	panic(s)
}

func execute(script Script) {
	log.Printf("%#v", script)
}

func main() {
	bin := bufio.NewReader(os.Stdin)
	for {
		if bytes, err := bin.ReadBytes('\n'); nil != err {
			if io.EOF == err {
				break
			} else {
				panic(err)
			}
		} else {
			yyParse(&Lexer{line: bytes})
		}
	}
}
