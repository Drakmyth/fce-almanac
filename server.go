package main

import (
	"context"
	"encoding/xml"
	"errors"
	"fmt"
	"html/template"
	"io"
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
			fmt.Printf("Checking for exe: %s\n", path2)
		}
		_, err := os.Stat(path2)
		return !errors.Is(err, fs.ErrNotExist)
	})

	if !found_exe {
		return errors.New("could not find FortressCraft Evolved at supplied path")
	}

	dataPath := filepath.Join(path, "Default/Data")
	gacPath := filepath.Join(dataPath, "GenericAutoCrafter")
	// TODO: Parse Default/Lang files
	// TODO: Maybe parse Default/Handbook

	dataFiles, err := getFileDirEntries(dataPath)
	if err != nil {
		return err
	}

	gacFiles, err := getFileDirEntries(gacPath)
	if err != nil {
		return err
	}

	var (
		genericAutoCrafterDataEntries []GenericAutoCrafterDataEntry
		craftDataRecords              []CraftData
		itemEntries                   []ItemEntry
		recipeSets                    []RecipeSet
		researchDataEntries           []ResearchDataEntry
		terrainDataEntries            []TerrainDataEntry
	)

	for _, fileEntry := range gacFiles {
		fmt.Printf("Scanning file: %v\n", fileEntry.Name())
		xmlFilePath := filepath.Join(gacPath, fileEntry.Name())
		data, err := os.ReadFile(xmlFilePath)
		if err != nil {
			return err
		}

		dataObj := &GenericAutoCrafterDataEntry{}
		err = xml.Unmarshal([]byte(data), &dataObj)
		if err != nil {
			return err
		}
		genericAutoCrafterDataEntries = append(genericAutoCrafterDataEntries, *dataObj)
	}

	for _, fileEntry := range dataFiles {
		fmt.Printf("Scanning file: %v\n", fileEntry.Name())
		xmlFilePath := filepath.Join(dataPath, fileEntry.Name())
		data, err := os.ReadFile(xmlFilePath)
		if err != nil {
			return err
		}

		switch fileEntry.Name() {
		case "Items.xml":
			dataObj := &ItemEntryDocument{}
			err = xml.Unmarshal([]byte(data), &dataObj)
			if err != nil {
				return err
			}
			itemEntries = dataObj.ItemEntries
		case "RecipeSets.xml":
			dataObj := &RecipeSetDocument{}
			err = xml.Unmarshal([]byte(data), &dataObj)
			if err != nil {
				return err
			}
			recipeSets = dataObj.RecipeSets
		case "Research.xml":
			dataObj := &ResearchDataEntryDocument{}
			err = xml.Unmarshal([]byte(data), &dataObj)
			if err != nil {
				return err
			}
			for _, researchDataEntry := range dataObj.ResearchDataEntries {
				researchDataEntry.ResearchRequirements = append(researchDataEntry.ResearchRequirements, researchDataEntry.MoreResearchRequirements...)
			}
			researchDataEntries = dataObj.ResearchDataEntries
		case "TerrainData.xml":
			dataObj := &TerrainDataEntryDocument{}
			err = xml.Unmarshal([]byte(data), &dataObj)
			if err != nil {
				return err
			}
			terrainDataEntries = dataObj.TerrainDataEntries
		default:
			machine := strings.Split(fileEntry.Name(), "Recipes")[0]
			dataObj := &CraftDataDocument{}
			err = xml.Unmarshal([]byte(data), &dataObj)
			if err != nil {
				return err
			}
			for _, craftData := range dataObj.Data {
				craftData.Machine = machine
				craftData.ScanRequirements = append(craftData.ScanRequirements, craftData.MoreScanRequirements...)
			}
			craftDataRecords = append(craftDataRecords, dataObj.Data...)
		}
	}

	fmt.Println()
	fmt.Printf("Generic Auto Crafter Data Entries: %d\n", len(genericAutoCrafterDataEntries))
	fmt.Printf("CraftData Records: %d\n", len(craftDataRecords))
	fmt.Printf("Item Entries: %d\n", len(itemEntries))
	fmt.Printf("Recipe Sets: %d\n", len(recipeSets))
	fmt.Printf("Research Data Entries: %d\n", len(researchDataEntries))
	fmt.Printf("Terrain Data Entries: %d\n", len(terrainDataEntries))

	return nil
}

func runserver(ctx context.Context, args []string) error {
	mux := muxpatterns.NewServeMux()

	// TODO: Separate partial endpoints into dedicated muxes. Example:
	// handbookMux := muxpatterns.NewServeMux()
	// handbookMux.HandleFunc("GET /{taskId}", getHandbook)
	// handbookMux.HandleFunc("DELETE /{taskId}", deleteHandbook)
	// handbookMux.HandleFunc("GET /{taskId}/edit", getEditHandbook)
	// handbookMux.HandleFunc("POST /", createHandbook)
	// handbookMux.HandleFunc("GET /{$}", getTaskHandbook)
	// mux.Handle("/handbook/", http.StripPrefix("/handbook", handbookMux))

	mux.HandleFunc("/handbook", getHandbook)
	mux.HandleFunc("/", getIndex)
	fmt.Printf("Applicaton listening on port %d...\n", port)
	err := http.ListenAndServe(fmt.Sprintf(":%d", port), mux)
	if err != nil {
		return err
	}

	return nil
}

func getFileDirEntries(path string) ([]fs.DirEntry, error) {
	dirEntries, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}

	files := []fs.DirEntry{}
	for _, entry := range dirEntries {
		if !entry.IsDir() {
			files = append(files, entry)
		}
	}

	return files, nil
}

func getHandbook(w http.ResponseWriter, r *http.Request) {
	err := executeTemplate(w, "./templates/handbook.tmpl.html", nil)
	if err != nil {
		log.Println(err)
	}
}

func getIndex(w http.ResponseWriter, r *http.Request) {
	err := executeTemplate(w, "./templates/index.tmpl.html", nil)
	if err != nil {
		log.Println(err)
	}
}

func executeTemplate(wr io.Writer, path string, data any) error {
	templatePaths, _ := filepath.Glob("./templates/partials/*.tmpl.html")
	templatePaths = append(templatePaths, path)

	tmpl, err := template.ParseFiles(templatePaths...)
	if err != nil {
		return err
	}

	return tmpl.ExecuteTemplate(wr, filepath.Base(path), data)
}
