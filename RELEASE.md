## Steps to release a new version

- Update CHANGELOG.md with changes made since last release
- Create a signed tag: ```git tag --sign -m "Release 1.1.1" 1.1.1```
- Push commits: ```git push```
- Push tag: ```git push --tags```
- Update github release page: https://github.com/cloudtools/troposphere/releases
- Upload build artifacts
