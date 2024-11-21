package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/urfave/cli/v2"
)

const dict48Path = "./dictionary/4_8.txt"

var (
	dict48 map[string]struct{}
)

func setupDict48() {
	dict48 = make(map[string]struct{}, 0)
	f, err := os.Open(dict48Path)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	r := bufio.NewReader(f)
	for {
		lineBytes, _, err := r.ReadLine()
		if err == io.EOF {
			break // finish
		}
		if lineBytes[0] == '*' {
			lineBytes = lineBytes[1:]
		}
		line := string(lineBytes)

		firstWord := strings.Split(line, " ")[0]

		firstWord = strings.ToLower(firstWord)

		dict48[firstWord] = struct{}{}
	}

	fmt.Printf("succeed to setup map48, len=%d\n", len(dict48Path))
}

func main() {
	app := &cli.App{
		Name:  "Parser",
		Usage: "parse the article",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "input",
				Usage: "input article name",
			},
			&cli.StringFlag{
				Name:  "output",
				Usage: "output  name",
			},
		},

		Action: parserMain,
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func filterWordByDict48(in string, out string) {
	fin, err := os.Open(in)
	if err != nil {
		panic(err)
	}
	defer fin.Close()

	scanner := bufio.NewScanner(fin)
	output := make(map[string]struct{}, 0)

	j := 0
	for scanner.Scan() {
		output1 := strings.Split(strings.Trim(scanner.Text(), " "), " ")
		for i := 0; i < len(output1); i++ {
			w := strings.Trim(output1[i], " ")
			w = strings.Trim(w, ",")
			w = strings.Trim(w, ".")
			w = strings.Trim(w, "?")
			w = strings.ToLower(w)
			if _, found := dict48[w]; found {
				fmt.Printf("hit 4, word=%s\n", w)
				output[w] = struct{}{}
			}
			j++

		}

	}
	fmt.Println(j)

	fmt.Printf("output_len=%d\n", len(output))
	fout, err := os.Create(out)
	if err != nil {
		panic(err)
	}
	defer fout.Close()

	for k, _ := range output {
		// fmt.Printf("word=%s, level=%d\n", k, v)
		//fmt.Printf("%s\n", k)
		_, err = fout.WriteString(k + "\n")
		if err != nil {
			panic(err)
		}
	}
}

func parserMain(c *cli.Context) error {
	setupDict48()
	filterWordByDict48(c.String("input"), c.String("output"))
	return nil
}
