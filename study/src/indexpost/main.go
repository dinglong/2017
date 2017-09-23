package main

import (
	"bufio"
	"flag"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

var eol byte = byte('\n')

// PostInfo prepare some data to insert into the template.
type PostInfo struct {
	FileName  string
	Title     string
	Timestamp string
}

// command paramates
var (
	templateFile string
	postdir      string
	help         bool
)

func init() {
	flag.StringVar(&templateFile, "template", "readme.template", "template file")
	flag.StringVar(&postdir, "postdir", "posts", "the post file dir")
	flag.BoolVar(&help, "help", false, "display this information")
	flag.Parse()
}

func main() {
	// display help information
	if help {
		flag.Usage()
		return
	}

	// load template
	readme, err := ioutil.ReadFile(templateFile)
	if err != nil {
		log.Println("read template file:", err)
		return
	}

	postinfos, err := readPostInfos(postdir)
	if err != nil {
		log.Println("read post dir:", err)
		return
	}

	// Create a new template and parse the letter into it.
	t := template.Must(template.New("letter").Parse(string(readme)))

	// Execute the template for each postinfos.
	err = t.Execute(os.Stdout, postinfos)
	if err != nil {
		log.Println("executing template:", err)
	}
}

// readPostInfos read PostInfos from dir
func readPostInfos(dir string) ([]PostInfo, error) {
	postinfos := make([]PostInfo, 0)

	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	for _, f := range files {
		p, err := readPostInfo(filepath.Join(dir, f.Name()))
		if err != nil {
			return nil, err
		}
		postinfos = append(postinfos, *p)
	}

	return postinfos, nil
}

// readPostInfo read PostInfo from file header
// # title
// > timestamp
func readPostInfo(fn string) (*PostInfo, error) {
	f, err := os.Open(fn)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	r := bufio.NewReaderSize(f, 1024)

	// read title
	title, err := r.ReadString(eol)
	if err != nil {
		return nil, err
	}

	// read timestamp
	timestamp, err := r.ReadString(eol)
	if err != nil {
		return nil, err
	}

	// generate PostInfo
	p := &PostInfo{
		FileName:  filepath.Base(f.Name()),
		Title:     strings.Trim(title, "#\t\n\r "),
		Timestamp: strings.Trim(timestamp, ">\t\n\r "),
	}

	return p, nil
}
