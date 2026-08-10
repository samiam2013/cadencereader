# CadenceReader

a blog syndication platform for reading in bits (built so the author can learn and faithfully claim experience in depth with Kubernetes)

## Dependencies

### k3s
The platform that runs the kubernetes manifests, a lightweight version of 'k8s' (abbreviation for kubernetes)

[instructions here](https://docs.k3s.io/quick-start)

and I had to add this line to my ./bashrc 
`export KUBECONFIG=$HOME/.kube/config`

and I had to fight a little to get kubectl working because when you do the k3s quick install `kubectl` is just a link to the k3s binary or something idk really

### Docker
because you have to build the containers for k3s to run

[Docker engine install instructions here](https://docs.docker.com/engine/install/)


### Go (aka Golang) 
the programming language,
* [instructions can be found here](https://go.dev/doc/install)
* [download a versoin for your arch here](https://go.dev/dl/)

### CNPG
the Cloud Native Postgres k8s operator which handles persistent data inside the cluster on your behalf, this is just a CRD

[the insecure quickstart guide instructions are here](https://cloudnative-pg.io/docs/devel/quickstart/#part-2-install-cloudnativepg) 

### Golang-migrate
[this is a schema migration tool](https://github.com/golang-migrate/migrate) you might not have to use directly, but it's good to have it for troubleshooting, it's what main webapp uses to build the pg tables and stuff
`go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest`

and I may be adding [flux for gitops](https://fluxcd.io/) so that this repo's main branch, which is protected, automatically redeploys to the cluster when a change is merged

## How to Run

### Local Development
Because I'm sick of waiting on docker to take ~30 seconds on any platform to rebuild the container, I set up `/localdev` with a script that starts a postgres container `/dockerdb.sh` and a sample env file so that you can source environment variables with a run script and get near-instant startup

You don't need k3s, CNPG, flux or golang-migrate to get started locally.

on MacOS, because that's the only place I've worked on development of this,
```bash
cd localdev/
cp sample.env .env
./dockerdb.sh
```
should be all it takes to get bootstrapped (but no promises), 
then
```bash
./run.sh
```
should open the app in your broswer

### k3s cluster

once all the dependencies are set up and working
```bash
./build.sh
```
in the main directory should do the trick. This script keeps track of a hash of the files that could change necessitating a rebuild of the docker continer(s) so you can modify the k8s manifests without waiting on a full rebuild to test those changes