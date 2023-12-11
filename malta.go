package malta

import (
	"bytes"
	_ "embed"
	"encoding/json"
	"html/template"
	"os"
	"path/filepath"
	"strings"

	"github.com/adrg/frontmatter"
	"github.com/alecthomas/chroma/v2/formatters/html"
	"github.com/yuin/goldmark"
	highlighting "github.com/yuin/goldmark-highlighting/v2"
	"github.com/yuin/goldmark/parser"
)

var config Config
var markdownFilePaths []string

//go:embed assets/template.html
var htmlTemplate []byte

//go:embed assets/main.css
var mainCss []byte

//go:embed assets/markdown.css
var markdownCss []byte

func main() {
	configJson, _ := os.ReadFile("config.json")
	json.Unmarshal(configJson, &config)

	navSections := []NavSection{}
	for _, sidebarSection := range config.Sidebar {
		navSection := NavSection{sidebarSection.Title, []NavPage{}}
		for _, sidebarSectionPage := range sidebarSection.Pages {
			navPage := NavPage{Title: sidebarSectionPage[0], Href: sidebarSectionPage[1]}
			navSection.Pages = append(navSection.Pages, navPage)
		}
		navSections = append(navSections, navSection)
	}

	if err := filepath.Walk("pages", walkPagesDir); err != nil {
		panic(err)
	}

	markdown := goldmark.New(
		goldmark.WithExtensions(
			highlighting.NewHighlighting(highlighting.WithFormatOptions(html.WithClasses(true))),
		),
	)

	os.RemoveAll("dist")

	for _, markdownFilePath := range markdownFilePaths {
		var matter struct {
			Title string `yaml:"title"`
		}

		markdownFile, _ := os.Open(markdownFilePath)
		defer markdownFile.Close()
		pageMarkdown, _ := frontmatter.MustParse(markdownFile, &matter)

		var markdownHtmlBuf bytes.Buffer

		if err := markdown.Convert(pageMarkdown, &markdownHtmlBuf, parser.WithContext(parser.NewContext())); err != nil {
			panic(err)
		}

		tmpl, _ := template.New("html").Parse(string(htmlTemplate))

		dstPath := strings.Replace(strings.Replace(markdownFilePath, "pages", "dist", 1), ".md", ".html", 1)

		if err := os.MkdirAll(filepath.Dir(dstPath), os.ModePerm); err != nil {
			panic(err)
		}
		dstHtmlFile, err := os.Create(dstPath)
		if err != nil {
			panic(err)
		}
		defer dstHtmlFile.Close()

		urlPathname := strings.Replace(strings.Replace(dstPath, "dist", "", 1), ".html", "", 1)
		urlPathname = strings.Replace(urlPathname, "/index", "", 1)

		err = tmpl.Execute(dstHtmlFile, Data{
			Markdown:    template.HTML(markdownHtmlBuf.String()),
			Name:        config.Name,
			Description: config.Description,
			Url:         config.Domain + urlPathname,
			Twitter:     config.Twitter,
			Title:       matter.Title,
			NavSections: navSections,
		})
		if err != nil {
			panic(err)
		}
	}

	os.WriteFile("dist/main.css", mainCss, os.ModePerm)
	os.WriteFile("dist/markdown.css", markdownCss, os.ModePerm)
}

type Data struct {
	Markdown    template.HTML
	Title       string
	Description string
	Twitter     string
	Url         string
	Name        string
	NavSections []NavSection
}

type NavSection struct {
	Title string
	Pages []NavPage
}

type NavPage struct {
	Title string
	Href  string
}

type Config struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Domain      string                 `json:"domain"`
	Twitter     string                 `json:"twitter"`
	Sidebar     []SidebarSectionConfig `json:"sidebar"`
}

type SidebarSectionConfig struct {
	Title string     `json:"title"`
	Pages [][]string `json:"pages"`
}

func walkPagesDir(path string, info os.FileInfo, err error) error {
	if err != nil {
		return err
	}
	if info.IsDir() {
		return nil
	}
	markdownFilePaths = append(markdownFilePaths, path)
	return nil
}