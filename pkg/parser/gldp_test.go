package parser

import "testing"

func TestGetArticatMetadata(t *testing.T) {

	p := Parser{}
	_ = p.Run()

	uniqueUries := make(map[string]bool)

	line_valid := "Downloading https://repo.maven.apache.org/maven2/org/codehaus/woodstox/stax2-api/4.1/stax2-api-4.1.jar to /tmp/gradle_download14056974681211842918bin"

	got, _ := p.GetArtifactMetadata(line_valid, uniqueUries)

	want := gradleArtifactMetadata{
		Version: "4.1",
		ArtifactId: "stax2-api",
		GroupId: "org.codehaus.woodstox",
		Repository: "https://repo.maven.apache.org/maven2/",
	}

	if got != want {
		t.Errorf("got %q, wanted %q", got, want)
	}
}
