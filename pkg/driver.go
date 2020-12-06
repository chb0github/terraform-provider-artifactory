package main

import (
	"fmt"
	artifactorynew "github.com/jfrog/jfrog-client-go/artifactory"
	"github.com/jfrog/jfrog-client-go/artifactory/auth"
	"github.com/jfrog/jfrog-client-go/artifactory/services"
	"github.com/jfrog/jfrog-client-go/config"
	"github.com/jfrog/jfrog-client-go/utils/log"
	"math/rand"
)

func main() {
	details := auth.NewArtifactoryDetails()
	details.SetUser("admin")
	details.SetPassword("password")
	details.SetUrl("http://localhost:8081/artifactory/")
	cfg, err := config.NewConfigBuilder().
		SetServiceDetails(details).
		SetDryRun(false).
		Build()
	log.SetLogger(log.NewLogger(log.DEBUG, nil))
	rt, err := artifactorynew.New(&details, cfg)
	gs := services.NewGroupService(rt.Client())
	gs.SetArtifactoryDetails(details)

	if rt == nil || err != nil {
		panic("oh shit")
	}

	description := "hello"

	name := fmt.Sprintf("test%d",  rand.Int())
	for exists, _ := gs.GroupExits(name); ; {
		if !exists {
			break
		}
	}
	group := services.Group{
		Name:            name,
		Description:     description,
		AutoJoin:        false,
		AdminPrivileges: true,
		Realm:           "internal",
		RealmAttributes: "",
	}
	err = gs.CreateGroup(group)
	if err != nil {
		fmt.Println(err)
	}
	g, err := gs.GetGroup(name)
	if err != nil {
		panic(err)
	}
	if g == nil {
		panic("no group")
	}
	if *g != group {
		panic("not equal")
	}
	group.Description = "Changed"
	err = gs.CreateGroup(group)
	g, err = gs.GetGroup(name)
	if err != nil {
		panic(err)
	}
	if g == nil {
		panic("no group")
	}
	if *g != group {
		panic("not equal")
	}
	if err := gs.DeleteGroup(name); err != nil {
		panic("boom")
	}

}
