module github.com/tencentess/ess-skills/toolkit/contract-atoms

go 1.21

require (
	github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common v1.3.48
	github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ess v1.3.48
	github.com/tencentess/ess-skills/toolkit/foundation v0.0.0
)

replace github.com/tencentess/ess-skills/toolkit/foundation => ../foundation

require (
	golang.org/x/sys v0.28.0 // indirect
	golang.org/x/term v0.27.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
