package main

import (
	"context"
	"errors"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/jba/muxpatterns"
	"github.com/peterbourgon/ff/v4"
	"github.com/peterbourgon/ff/v4/ffhelp"
)

var tmpl *template.Template = nil
var config string = ""
var verbose bool = false
var port int = 8000

func main() {
	rootFlags := ff.NewFlagSet("fce-almanac")
	rootFlags.StringVar(&config, 'c', "config", config, "config file")
	rootFlags.IntVar(&port, 'p', "port", port, "application port")
	rootFlags.BoolVar(&verbose, 'v', "verbose", "increase log verbosity")
	rootCmd := &ff.Command{
		Name:      "fce-almanac",
		ShortHelp: "Run almanac backend server",
		Usage:     "fce-almanac [FLAGS] [SUBCOMMAND]",
		Flags:     rootFlags,
		Exec:      runserver,
	}

	initdbFlags := ff.NewFlagSet("initdb").SetParent(rootFlags)
	initdbCmd := &ff.Command{
		Name:      "initdb",
		Usage:     "fce-almanac initdb <path>",
		ShortHelp: "Initialize db with FortressCraft Evolved data",
		Flags:     initdbFlags,
		Exec:      initdb,
	}
	rootCmd.Subcommands = append(rootCmd.Subcommands, initdbCmd)

	err := rootCmd.ParseAndRun(context.Background(), os.Args[1:],
		ff.WithEnvVarPrefix("FCE"),
		ff.WithConfigFileFlag("config"),
		ff.WithConfigFileParser(ff.PlainParser),
		ff.WithConfigIgnoreUndefinedFlags())

	switch {
	case errors.Is(err, ff.ErrHelp):
		fmt.Fprintf(os.Stderr, "%s\n", ffhelp.Command(rootCmd))
		os.Exit(0)
	case err != nil:
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func initdb(ctx context.Context, args []string) error {
	fmt.Fprintf(os.Stdout, "TODO: Initialize db\n")
	return nil
}

func runserver(ctx context.Context, args []string) error {
	mux := muxpatterns.NewServeMux()
	tmpl, _ = template.ParseGlob("./templates/*.tmpl.html")

	mux.Handle("/", http.FileServer(http.Dir("public")))
	fmt.Fprintf(os.Stdout, "Applicaton listening on port %d...\n", port)
	err := http.ListenAndServe(fmt.Sprintf(":%d", port), mux)
	log.Fatal(err)
	return err
}
