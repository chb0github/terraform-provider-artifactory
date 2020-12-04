package main

import (
	"fmt"
	"github.com/atlassian/terraform-provider-artifactory/pkg/artifactory"
	artifactorynew "github.com/jfrog/jfrog-client-go/artifactory"
	"github.com/jfrog/jfrog-client-go/artifactory/auth"
	"github.com/jfrog/jfrog-client-go/artifactory/services"
	"github.com/jfrog/jfrog-client-go/config"
	"github.com/jfrog/jfrog-client-go/utils/log"
)

func main() {
	provider := artifactory.Provider()
	details := auth.NewArtifactoryDetails()
	details.SetUser("admin")
	details.SetPassword("Password1")
	details.SetUrl("http://localhost:8081/artifactory/")
	cfg, err := config.NewConfigBuilder().
		SetServiceDetails(details).
		SetDryRun(false).
		Build()
	log.SetLogger(log.NewLogger(log.DEBUG,nil))
	rt, err := artifactorynew.New(&details, cfg)
	gs := services.NewGroupService(rt.Client())
	gs.SetArtifactoryDetails(details)

	if rt == nil || err != nil {
		panic("oh shit")
	}

	sources := provider.DataSources()
	description := "hello"
	name := "ABC123"
	err = gs.CreateGroup(services.Group{
		Name:            &name,
		Description:     &description,
		AutoJoin:        new(bool),
		AdminPrivileges: new(bool),
		Realm:           nil,
		RealmAttributes: nil,
	})
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("hello %s",sources)
}
