# gldp
Gradle log dependency parser. Gldp Basically reads gradle log from stdin, finds and parces artifacts metadata and output unique dependency information to Stdout in desired format

## Examples
```
# Specify your internal repo
cat /path/to/gradle/log | gldp -repos https://your.internal.maven.repo/maven2/

# Override output format
cat /path/to/gradle/log | gldp -repos https://your.internal.maven.repo/maven2/ -f '{{ printf "%v %v %v %v \n" .Repository .GroupId .ArtifactId .Version }}'
```
