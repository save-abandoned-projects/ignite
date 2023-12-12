:hand: [Weave Ignite](https://github.com/weaveworks/ignite) is a tool for easily managing firecracker vm in the container way. This is a great open source project, but it has been archived now.

### What Do Next

All developed under [charlie-init](https://github.com/save-abandoned-projects/ignite/tree/charlie-init) branch
### No test
- docker container runtime
- gitops
#### [Libgitops](https://github.com/save-abandoned-projects/libgitops)

libgitops is a project that depended by ignite, it used to store the ignite metadata, just like etcd.

- [x] Migrate project from weaveworks to save-abandoned-projects
- [x] Fix [bugs](https://github.com/save-abandoned-projects/libgitops/issues/2)
- [x] Merge changes to main branch
- [x] Upgrade golang and go module

#### [Ignite](https://github.com/save-abandoned-projects/ignite)

- [x] Migrate project from weaveworks to save-abandoned-projects
- [x] Upgrade golang and module
- [x] Compatible with the latest libgitops
- [ ] Test the modified ignite, include module test and e2e test
- [ ] Fix [bugs](https://github.com/save-abandoned-projects/ignite/issues)
- [ ] Merge changes to main branch
- [ ] Go on with the [roadmap](https://github.com/weaveworks/ignite/blob/main/docs/roadmap.md), but may not implement all the feature
  - [ ] Add Virtual Kubelet support to `ignited`
  - [ ] Use device-mapper Thin Provisioning for layering image -> kernel -> resize -> writable overlay
  - [ ] Add support for CSI volumes
- [ ] Support GPU？

If you are intresting about this, join and make it better!

## License

[Apache 2.0](file:///Users/charlie.liu/go/src/github.com/save-abandoned-projects/ignite/LICENSE)
