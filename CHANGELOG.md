# v0.7.3
- update to go1.23
- update github workflows
- add noOp function to correctly return With functions that only set app data
- documentation fixes

# v0.7.2
- https://github.com/skeletonkey/lib-instance-gen-go/issues/18
  - update linters
  - update workflow for running the linters

# v0.7.1
- https://github.com/skeletonkey/lib-instance-gen-go/issues/8
  - add `WithCodeOwners` method to add simple CODEOWNERS file
  - standardized on os.OpenFile when overwriting files
- https://github.com/skeletonkey/lib-instance-gen-go/issues/20
  - removed the changelog workflow
    - NOTE: changelog workflow should be manually deleted
  - added an example pull_request_template.md which can be manually added to projects

# v0.7.0
- https://github.com/skeletonkey/lib-instance-gen-go/issues/4
  - add `SetupApp` so that With* call can happen in any order
  - add `Generate` to create application files

# v0.6.1
- forgot to add the template file for app level config

# v0.6.0
- add WithConfig

# v0.5.0
- add changelog watch workflow
- fix linter workflow
- update go version to 1.22
- update README.md

# v0.4.0
- add WithGoVersion

# v0.3.0
- add CHANGELOG.md
- add .golangci.toml to unify linting
