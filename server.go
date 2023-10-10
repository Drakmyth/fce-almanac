package main

import (
	"context"
	"encoding/xml"
	"errors"
	"fmt"
	"html/template"
	"io/fs"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/jba/muxpatterns"
	"github.com/peterbourgon/ff/v4"
	"github.com/peterbourgon/ff/v4/ffhelp"
	"golang.org/x/exp/slices"
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
	if len(args) == 0 {
		return errors.New("must supply FortressCraft Evolved installation path")
	}

	path := args[0]
	exe_names := []string{"FortressCraft.exe", "FC_64.exe"}
	found_exe := slices.ContainsFunc(exe_names, func(s string) bool {
		path2 := filepath.Join(path, s)
		if verbose {
			fmt.Fprintf(os.Stdout, "Checking for exe: %s\n", path2)
		}
		_, err := os.Stat(path2)
		return !errors.Is(err, fs.ErrNotExist)
	})

	if !found_exe {
		return errors.New("could not find FortressCraft Evolved at supplied path")
	}

	dataPath := filepath.Join(path, "Default/Data")
	gacPath := filepath.Join(dataPath, "GenericAutoCrafter")

	dataFiles := getFileDirEntries(dataPath)
	gacFiles := getFileDirEntries(gacPath)

	genericAutoCrafterDataEntries := []GenericAutoCrafterDataEntry{}
	var (
		craftDataRecords    []CraftData
		itemEntries         []ItemEntry
		recipeSets          []RecipeSet
		researchDataEntries []ResearchDataEntry
		terrainDataEntries  []TerrainDataEntry
	)

	for _, fileEntry := range gacFiles {
		fmt.Printf("Scanning file: %v\n", fileEntry.Name())
		xmlFilePath := filepath.Join(gacPath, fileEntry.Name())
		data, err := os.ReadFile(xmlFilePath)
		if err != nil {
			log.Fatal(err)
			os.Exit(1)
		}

		dataObj := &GenericAutoCrafterDataEntry{}
		err = xml.Unmarshal([]byte(data), &dataObj)
		if err != nil {
			log.Fatal(err)
			os.Exit(1)
		}
		genericAutoCrafterDataEntries = append(genericAutoCrafterDataEntries, *dataObj)
	}

	for _, fileEntry := range dataFiles {
		fmt.Printf("Scanning file: %v\n", fileEntry.Name())
		xmlFilePath := filepath.Join(dataPath, fileEntry.Name())
		data, err := os.ReadFile(xmlFilePath)
		if err != nil {
			log.Fatal(err)
			os.Exit(1)
		}

		switch fileEntry.Name() {
		case "Items.xml":
			dataObj := &ItemEntryDocument{}
			err = xml.Unmarshal([]byte(data), &dataObj)
			if err != nil {
				log.Fatal(err)
				os.Exit(1)
			}
			itemEntries = dataObj.ItemEntries
		case "RecipeSets.xml":
			dataObj := &RecipeSetDocument{}
			err = xml.Unmarshal([]byte(data), &dataObj)
			if err != nil {
				log.Fatal(err)
				os.Exit(1)
			}
			recipeSets = dataObj.RecipeSets
		case "Research.xml":
			dataObj := &ResearchDataEntryDocument{}
			err = xml.Unmarshal([]byte(data), &dataObj)
			if err != nil {
				log.Fatal(err)
				os.Exit(1)
			}
			researchDataEntries = dataObj.ResearchDataEntries
		case "TerrainData.xml":
			dataObj := &TerrainDataEntryDocument{}
			err = xml.Unmarshal([]byte(data), &dataObj)
			if err != nil {
				log.Fatal(err)
				os.Exit(1)
			}
			terrainDataEntries = dataObj.TerrainDataEntries
		default:
			machine := strings.Split(fileEntry.Name(), "Recipes")[0]
			dataObj := &CraftDataDocument{}
			err = xml.Unmarshal([]byte(data), &dataObj)
			if err != nil {
				log.Fatal(err)
				os.Exit(1)
			}
			for _, craftData := range dataObj.Data {
				craftData.Machine = machine
				craftDataRecords = append(craftDataRecords, craftData)
			}
		}
	}

	fmt.Println()
	fmt.Fprintf(os.Stdout, "Generic Auto Crafter Data Entries: %d\n", len(genericAutoCrafterDataEntries))
	fmt.Fprintf(os.Stdout, "CraftData Records: %d\n", len(craftDataRecords))
	fmt.Fprintf(os.Stdout, "Item Entries: %d\n", len(itemEntries))
	fmt.Fprintf(os.Stdout, "Recipe Sets: %d\n", len(recipeSets))
	fmt.Fprintf(os.Stdout, "Research Data Entries: %d\n", len(researchDataEntries))
	fmt.Fprintf(os.Stdout, "Terrain Data Entries: %d\n", len(terrainDataEntries))

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

func getFileDirEntries(path string) []fs.DirEntry {
	dirEntries, err := os.ReadDir(path)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	files := []fs.DirEntry{}
	for _, entry := range dirEntries {
		if !entry.IsDir() {
			files = append(files, entry)
		}
	}

	return files
}
