module github.com/atlassian/terraform-provider-artifactory

require (
	github.com/atlassian/go-artifactory/v2 v2.5.0
	github.com/hashicorp/terraform v0.12.29
	github.com/jfrog/jfrog-client-go v0.13.1
	github.com/stretchr/testify v1.5.1
)

replace github.com/jfrog/jfrog-client-go => /Users/christianb/dev/jfrog/jfrog-client-go

go 1.15
