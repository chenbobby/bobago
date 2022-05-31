package main

import (
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

type Recipe struct {
	filePath string
	template string
	bindings map[string]any
}

func main() {
	siteDirPath := ""
	if len(os.Args) > 1 {
		siteDirPath = os.Args[1]
	}

	if err := os.Chdir(siteDirPath); err != nil {
		panic(err)
	}

	if err := run(); err != nil {
		panic(err)
	}
}

func run() error {
	// Find all files.
	pagesDirPath := "pages"
	pagePaths, err := filePaths(pagesDirPath)
	if err != nil {
		return err
	}

	// Parse each file into a recipe.
	var recipes []Recipe
	for _, pagePath := range pagePaths {
		recipe, err := recipe(pagePath)
		if err != nil {
			return err
		}
		recipes = append(recipes, recipe)
	}

	// Prepare the output directory.
	buildDirPath := "_build"
	if err := os.RemoveAll(buildDirPath); err != nil {
		return err
	}
	if err := os.MkdirAll(buildDirPath, fs.ModePerm); err != nil {
		return err
	}

	// Render the recipe and save the output into a new file.
	for _, recipe := range recipes {
		output, err := recipe.render()
		if err != nil {
			return err
		}

		outputPath := filepath.Join(buildDirPath, recipe.name()+".html")
		os.MkdirAll(filepath.Dir(outputPath), fs.ModePerm)
		if err := os.WriteFile(outputPath, []byte(output), 0644); err != nil {
			return err
		}
	}

	return nil
}

func filePaths(dirPath string) (paths []string, err error) {
	fn := func(path string, info fs.DirEntry, _err error) error {
		if info.IsDir() {
			return nil
		}
		paths = append(paths, path)
		return nil
	}

	err = filepath.WalkDir(dirPath, fn)
	return
}

func recipe(filePath string) (recipe Recipe, err error) {
	templateBytes, err := os.ReadFile(filePath)
	if err != nil {
		return
	}

	bindings := map[string]any{}
	return Recipe{filePath, string(templateBytes), bindings}, nil
}

func (r Recipe) render() (string, error) {
	return "booop.", nil
}

func (r Recipe) name() string {
	return strings.TrimSuffix(r.filePath, filepath.Ext(r.filePath))
}
