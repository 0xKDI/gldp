package parser

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"regexp"
	"strings"
	"text/template"
)


type gradleArtifactMetadata struct {
	GroupId    string
	ArtifactId string
	Version    string
	Repository    string
}


type Parser struct {
	repos       []string
	outTemplate *template.Template
	regex *regexp.Regexp
}


func (p *Parser) ParseFlags() error {
	r := flag.String("repos", "https://repo.maven.apache.org/maven2/", "Comma separated list of repositories")
	f := flag.String("f", `{{ printf "%v:%v:%v\n" .GroupId .ArtifactId .Version }}`, "artifact output format")
	rgx := flag.String("regex", `Downloading.*\.jar.*`, "regex expretion to search for deps")

	flag.Parse()

	// split string to list
	p.repos = strings.Split(*r, ",")

	// Parse template
    t, err := template.New("output-format").Parse(*f)
    if err != nil { return err }
	p.outTemplate = t

	// define regex
	regex, err := regexp.Compile(*rgx)
    if err != nil { return err }
	p.regex = regex

    return nil
}


func (p *Parser) Run() error {
    err := p.ParseFlags()
	if err != nil { return err }

	err = p.Process()
	if err != nil { return err }

	return nil
}


func (p *Parser) GetArtifactMetadata(line string, uries map[string]bool) (gradleArtifactMetadata, error) {
	var metadata gradleArtifactMetadata

	// return empty metadata and error if line doesn't match regex
	if !p.regex.MatchString(line) {
		return metadata, fmt.Errorf("Line doesn't match given regex")
	}

	// get artifact uri from the line
	uri := strings.Split(line, " ")[1]

	var trimmedUri string

	// cut repo from uri if exists
	for _, repo := range p.repos {
		if strings.Contains(uri, repo) {
			trimmedUri = strings.TrimPrefix(uri, repo)
			metadata.Repository = repo
			break
		}
	}

	if metadata.Repository == "" {
		return metadata, fmt.Errorf("not repo found in a string")
	}

	trimmedUri = strings.TrimPrefix(trimmedUri, "/")

	// continue if uri exists in the map
	if uries[trimmedUri] {
		return metadata, fmt.Errorf("Duplicate uri")
	}

	// add uri to the map
	uries[trimmedUri] = true

	// split artifact uri with a "/"
	splitedUri := strings.Split(trimmedUri, "/")

	// get slice length
	sliceLen := len(splitedUri)

	// assign metadata
	metadata.GroupId = strings.Join(splitedUri[:sliceLen-3], ".")
	metadata.ArtifactId = splitedUri[sliceLen-3]
	metadata.Version = splitedUri[sliceLen-2]

	return metadata, nil
}


func (p *Parser) Process() error {

	s := bufio.NewScanner(os.Stdin)

	// trimmed artifact uris
	uniqueUries := make(map[string]bool)

	for s.Scan() {

		// get line from stdint
		line := s.Text()

		// get artifact metadata
		artifactMetadata, err := p.GetArtifactMetadata(line, uniqueUries)

		// continue if we can't get artifact metadata from line
		if err != nil { continue }

		// Execute template and pring to Stdout
		err = p.outTemplate.Execute(os.Stdout, artifactMetadata)

		if err != nil { return err }
	}

	return nil
}
