package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/coreos/go-semver/semver"
	"github.com/google/go-github/github"
)

/*
	[DONE] remove pre-release versions
	add test-cases
	[DONE] don't hardcode pathname
	support pagination
*/

// LatestVersions returns a sorted slice with the highest version as its first element and the highest version of the smaller minor versions in a descending order
func LatestVersions(releases []*semver.Version, minVersion *semver.Version) []*semver.Version {
	var tempSlice semVerSlice
	var versionSlice semVerSlice
	// This is just an example structure of the code, if you implement this interface, the test cases in main_test.go are very easy to run
	for _, ver := range releases {
		if ver.LessThan(*minVersion) || ver.Equal(*minVersion) || ver.PreRelease != "" {
			continue
		} else {
			tempSlice = append(tempSlice, ver)
		}
	}
	tempSlice.sort()

	if len(tempSlice) == 0 {
		return semVerSlice{}
	}
	vComps := strings.Split(tempSlice[0].String(), ".")
	prefix := vComps[0] + "." + vComps[1]
	versionSlice = append(versionSlice, tempSlice[0])

	for i := range tempSlice {
		if strings.HasPrefix(tempSlice[i].String(), prefix) {
			continue
		} else {
			vComps = strings.Split(tempSlice[i].String(), ".")
			prefix = vComps[0] + "." + vComps[1]
			versionSlice = append(versionSlice, tempSlice[i])
		}
	}
	return versionSlice
}

// Here we implement the basics of communicating with github through the library as well as printing the version
// You will need to implement LatestVersions function as well as make this application support the file format outlined in the README
// Please use the format defined by the fmt.Printf line at the bottom, as we will define a passing coding challenge as one that outputs
// the correct information, including this line
func main() {
	b, err := ioutil.ReadFile("input.txt") // just pass the file name
	if err != nil {
		fmt.Print(err)
	}
	data := strings.Split(string(b), "\n")[1:]

	for i := range data {
		fileSlice := strings.Split(data[i], ",")
		fmt.Printf("latest versions of %s: %s\n", fileSlice[0], fetchReleases(fileSlice[0], fileSlice[1]))
	}

}
func fetchReleases(pathName string, minVersion string) []*semver.Version {
	path := strings.Split(pathName, "/")
	client := github.NewClient(nil)
	ctx := context.Background()
	opt := &github.ListOptions{PerPage: 10}
	releases, _, err := client.Repositories.ListReleases(ctx, path[0], path[1], opt)
	if err != nil {
		fmt.Printf("Request failed with error: %s\n", err)
		os.Exit(1)
	}
	minV := semver.New(minVersion)
	allReleases := make([]*semver.Version, len(releases))
	for i, release := range releases {
		versionString := *release.TagName
		if versionString[0] == 'v' {
			versionString = versionString[1:]
		}
		allReleases[i] = semver.New(versionString)
	}
	return LatestVersions(allReleases, minV)
}

type semVerSlice []*semver.Version

func (vSlice semVerSlice) swap(index1 int, index2 int) {
	temp := vSlice[index1]
	vSlice[index1] = vSlice[index2]
	vSlice[index2] = temp
}

func (vSlice semVerSlice) sort() {
	for i := range vSlice {
		for j := range vSlice {
			if !(vSlice[i].LessThan(*vSlice[j])) {
				vSlice.swap(i, j)
			}
		}
	}
}

type ListOptions struct {
	// For paginated result sets, page of results to retrieve.
	Page int `url:"page,omitempty"`

	// For paginated result sets, the number of results to include per page.
	PerPage int `url:"per_page,omitempty"`
}
