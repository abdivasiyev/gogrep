package app

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"regexp"

	"github.com/fatih/color"
)

type App struct {
	pattern    *regexp.Regexp
	srcList    []string
	whiteColor *color.Color
	redColor   *color.Color
}

func New(pattern string, srcList []string) (*App, error) {
	re, err := regexp.Compile(pattern)
	if err != nil {
		return nil, err
	}

	return &App{
		pattern:    re,
		srcList:    srcList,
		whiteColor: color.New(color.FgWhite, color.Bold),
		redColor:   color.New(color.FgRed, color.Bold).Add(color.Underline),
	}, nil
}

func (a *App) Run(ctx context.Context) error {
	for _, srcPath := range a.srcList {
		if err := a.walkSrc(srcPath); err != nil {
			return err
		}
	}

	return nil
}

func (a *App) walkSrc(path string) error {
	return filepath.WalkDir(path, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			fmt.Printf("Error while visiting %s: %s\n", path, err)
			return err
		}

		if d.IsDir() {
			return nil
		}

		f, err := os.Open(path)
		if err != nil {
			fmt.Printf("Error while opening %s: %s\n", path, err)
			return err
		}
		defer f.Close()

		matches := make([]Match, 0, 100)
		sc := bufio.NewScanner(f)
		line := 0

		for sc.Scan() {
			if sc.Err() != nil {
				fmt.Printf("Error while reading %s: %s\n", path, sc.Err())
				return sc.Err()
			}

			line++

			if len(sc.Text()) == 0 {
				continue
			}

			if !re.Match(sc.Bytes()) {
				continue
			}

			indexes := re.FindAllIndex(sc.Bytes(), -1)

			match := Match{
				Line:      line,
				Data:      sc.Bytes(),
				Positions: make([]Position, len(indexes)),
			}

			for i, idx := range indexes {
				match.Positions[i] = Position{
					Start: idx[0],
					End:   idx[1],
				}
			}

			matches = append(matches, match)
		}

		if len(matches) == 0 {
			return nil
		}

		if len(srcList) > 1 {
			fmt.Println(path)
		}
		for _, m := range matches {
			linePrinter.Printf("%d: ", m.Line)

			for i, pos := range m.Positions {
				if i > 0 {
					fmt.Printf("%s", m.Data[m.Positions[i-1].End:pos.Start])
				} else if i == 0 {
					fmt.Printf("%s", m.Data[:pos.Start])
				}

				matchPrinter.Printf("%s", m.Data[pos.Start:pos.End])

				if i == len(m.Positions)-1 {
					fmt.Printf("%s", m.Data[pos.End:])
				}
			}
			fmt.Printf("\n")
		}
		fmt.Println()

		return nil
	})
}
