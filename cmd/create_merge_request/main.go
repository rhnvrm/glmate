package main

import (
	"fmt"
	"log"
	"os"

	"github.com/xanzy/go-gitlab"
	"github.com/rhnvrm/glmate/pkg/version"
)

// type argT struct {
// 	cli.Helper
// 	Token     string `cli:"token" usage:"gitlab private ci token" dft:"$GITLAB_TOKEN"`
// 	Root      string `cli:"root" dft:"https://gitlab.com/"`
// 	ProjectID string `cli:"*project" usage:"project"`
// }

func main() {
	git := gitlab.NewClient(nil, os.Getenv("GITLAB_TOKEN"))
	git.SetBaseURL(os.Getenv("GITLAB_BASE_URL"))
	pid := "rhnvrm/test-merge-auto"

	tl, _, err := git.Tags.ListTags(pid, nil)

	var prevTag, nextTag string
	if len(tl) == 0 {
		prevTag, nextTag = "nil", "v0.0.1"
	} else {
		prevTag = tl[0].Name
		nextTag, err = version.BumpVersion(prevTag, version.PATCH)
		if err != nil {
			log.Fatal(err)
		}
	}

	log.Printf("prevTag: %s, nextTag: %s", prevTag, nextTag)

	sourceBranch, _, err := git.Branches.GetBranch(pid, "develop")
	if err != nil {
		log.Fatal(err)
	}

	targetBranch, _, err := git.Branches.GetBranch(pid, "master")
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("source commit:\n%#v\ntarget commit:\n%#v", sourceBranch.Commit, targetBranch.Commit)

	todoMsg := "todo"
	newTag, _, err := git.Tags.CreateTag(pid, &gitlab.CreateTagOptions{
		TagName:            &nextTag,
		Ref:                &sourceBranch.Commit.ID,
		Message:            &todoMsg,
		ReleaseDescription: &todoMsg,
	})

	log.Printf("new tag: %#v", newTag)

	mr, _, err := git.MergeRequests.CreateMergeRequest(pid, &gitlab.CreateMergeRequestOptions{
		Title:        &todoMsg,
		Description:  &todoMsg,
		SourceBranch: &sourceBranch.Name,
		TargetBranch: &targetBranch.Name,
	})
	if err != nil {
		log.Fatal(err)
	}

	log.Println("created merge request:")
	log.Println(mr.WebURL)

	commentBody := fmt.Sprintf("added tag: %s", nextTag)
	_, _, err = git.Notes.CreateMergeRequestNote(pid, mr.IID,
		&gitlab.CreateMergeRequestNoteOptions{
			Body: &commentBody,
		})
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("added comment: %s", commentBody)

	log.Println("done")
}
