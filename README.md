# Gosh (GitOps Shell)

![Build](https://github.com/ndriessen/gosh/actions/workflows/build.yml/badge.svg?branch=master)
![CodeQL](https://github.com/ndriessen/gosh/actions/workflows/codeql-analysis.yml/badge.svg?branch=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/ndriessen/gosh)](https://goreportcard.com/report/github.com/ndriessen/gosh)

Gosh is an opinionated way of implementing a GItOps based deployment flow.

It assumes a certain structure based on the [Gosh Repository Template] and offer a convenient CLI to interact with it.

> NOTE: This document currently references some paint points from our current deployment tools and processes. This document should involve into a more objective/general version of the concept

# Opinionated GitOps

If you want to get a general introduction into GitOps, head over to https://www.gitops.tech/ for a good starting collection of resources.

## Idea

Manage everything as code, including your infrastructure and application deployment configurations.
Tools like ArgoCD can watch GIT repositories and apply changes in them to Kubernetes clusters.

There are many advantages, as any change is traceable in history, and a rollback is as simple as reverting a commit
Next to that, one can take advantage of GIT features to manage deployment workflows instead of writing custom scripting.

This project focusses on the GIT repository structure and contents to provide that deployment configuration so we can leverage it in our Kubernetes deployments

## Requirements

- Designed from a continuous deployment mindset first
- Support on-premise releases, including hotfixes for older versions
- Reproducible releases: every release needs to be fully reproducible (with the assumption that application artifacts (container images etc) are stored safely somewhere)
- Support for development lifecycle: let teams decide which versions to make available to other teams and which ones to release for production
- Flexible environments, tracked by code: support auto-provisioning environments using specific versions by committing code
- Easy to use from project pipelines: ideally provide a dedicated CLI to handle common things like adding an application, updating a version etc.
- Integrate with central configuration management based on GIT for things like platform configuration, feature toggles, etc.

As a transition period, we still need to support our current processes, so additionally
- Provide an easy way to integrate bamboo specs with the new way of versioning

## Solution Concept

### The basic flow

From a high-level point of view, a project only needs to have deployment descriptors in its sources (or potentially published as output)
and a container image that is published to a registry.

To deploy this version somewhere (or make it available for deployment at least), the project's pipeline commits a change to the deployment GIT repo
to update the produced version.

```
project repo    -> sources: Deployment Descriptors
                -> output: Container Image

project pipeline -> output: image:version -> commit version to deployment repo -> version available for deployment <- ArgoCD Operator applies changes to cluster(s)
                
```

### Lifecycle Stages

Today we use different 'stages' in the release dashboard to determine which versions are in which 'stage' of the 'lifecycle'
Although there are limitations with the current implementation, it does enable some useful features
- teams can get all (or a subset of) application versions in a certain 'state' easily deployed on their environments for testing
- versions can be 'promoted' along the stages during different stages in their CD pipeline, or not, if they fail some gate
- teams can control when promotion happens, and decide when to make a version available for "production use" in their pipelines

However, today some limitations are
- fixed stages, not easy to add them
- releases are tightly coupled to stages, making it hard to do multiple releases at once
- "overriding" a certain version, means updating other teams pipelines to also always deploy to a certain env, this is not flexible and doesn't give the team that owns the environment full control over what is being deployed

Lifecycle stages will now be a completely standalone concept:

> A lifecycle stage is a **collection of versions** for applications, that can be reused.

Lifecycle stages do not *do* anything in itself, in order to use the versions, you need to also define a release or environment that uses them

### Releases

Releases contain, similar to stages, a list of versions of the applications in the release,
Other than stages, however, they do apply them to the actual deployment config.

> A release is a deployment configuration that specifies which applications and application versions to deploy

While a releases defines effectively *what* is going to be deployed, it does not actually *trigger* a deployment,
to trigger deployments, you still need to define targets

### Targets

Targets define actual deployments that are going to be executed whenever anything changes.
While you can use targets to deploy a specific release, you do not have to.
You can define any deployment configuration in a target, but this is typically rather inconvenient, so most targets will:
- use a release to deploy a specific release
- use a stage to deploy apps in a certain lifecycle stage, and override certain versions
- use a combination of a release and stage and specific overrides
- use a combination of different stages to mix versions
- use a release as a base and specific overrides to develop a hotfix e.g.
- ...

> A target defines an actual deployment and can reuse stages, releases, or define a completely custom deployment configuration.

### Environment classes

> An environment class can be used as a blueprint for common types of targets, e.g. dev, production, etc. ...

Environment classes are just convenience, they allow to define common deployment configuration that are not release but target specific.
These are commonly used to e.g.:
- enable/disable certain features for dev targets, but not for staging or production.
- apply certain debug configurations for dev targets
- etc...

You are not obliged to use an environment class, but it will definitely help avoid repetition and keep your code DRY

### Applications and application groups

> An application group is a convenient way to deploy a set of applications

An application group is simple a way to group applications, it doesn't control version for the contained apps,
but it makes it easy to deploy a complex application composed out of many components without needing to list all of them for each deployment

> An application defines what to deploy

An application is the component that produces the actual output. This generates manifest and other files needed to actually deploy your application
While the version being deployed is controlled by all the other concepts, the application itself always defines what
is being deployed

## Release Flows

### Types of releases

Our solution has been designed with GitOps and continuous deployment as the end-goal in mind. However, reality is that we do need to support on-premise releases.

In our solution, this distinction does not really matter, we could define a `saas-production` release that all targets for our SaaS environments use,
pipelines of individual apps could promote to this release after all (automated) gates have been passed e.g.
On the other hand, one can define a new release `R2021-R3` and either start from the `production` stage, or even just manually list which versions to include
For hotfixes, a release `R2021-R1.1` could be created that starts from `R2021-R1` and updates only affected apps

#### SaaS Releases

SaaS releases will (eventually) follow a Continuous Deployment flow, where we will have a single `saas-production` release that is used to control our deployments.

#### Product Releases

Product releases are our *on-premise releases*, the flow could go like this

1. Create new release based on `production` stage versions
1. Re-use or add a target to deploy your new release somewhere
1. Release validation happens
1. If issues are found
    1. Team fixes issues forward and make sure the patched version get 'released' to the `production` stage
    2. Team can use branches or any desired workflow to accomplish this depending on the needs
    1. The release is updated with the new fixed version or any newly released version of other apps

This does mean, you do not *freeze* versions when you create a release, which might be useful during validation, but when the release is finalized they should become static snapshots.

To accomplish this, one could
1. Create a tag in GIT to indicate which commit defined the released snapshot
1. Just rebuild the release from that tag to reproduce
1. Make sure all versions are listed statically and **do not use a stage** anymore

> The solution should provide a CLI for these type of tasks

### Product HotFix releases

These would work similar as product release, although from a development point of view, affected application will also need to start from the original version of the code to create a patched version.
Second, these versions should never be pushed to stages as they are not forward fixing. The project pipeline should update hotfix versions directly in the hotfix release

> The solution should provide a CLI for these type of tasks


# Usage

## Pre-requisites

- Install GIT
- Install Docker

## Repository structure
| Path          | Path      | Path              | Path      | Description |     
|:----           |:----       |:----               |:----       |:----       |
| `inventory`   |           |                   |           |           |
|               | `classes` |                   |           |           |
|               |           | `stages`          |           | This folder contains lifecycle `stages`, every file is the name of a `stage` |
|               |           | `releases`        |           | This folder contains `releases` |
|               |           |                   | `hotfix`  | This folder contains hotfix `releases` |
|               |           |                   | `product` | This folder contains product `releases` |
|               |           |                   | `stage`   | This folder contains CD `releases` based on `stages` |
|               |           | `env`             |           | This folder contains `environment classes` |
|               |           | `apps`            |           | This folder contains the `applications` deployment configuration (everything except the versions) |
|               | `targets` |                   |           | This folder contains all `targets`, each file is a target |

## Applications

### Create an application

> TODO: update with apps and app groups examples

In order to deploy anything, the deployment repo needs to output deployment descriptors.
Every application needs to create a file in `inventory/classes/apps/APP_GROUP` that defines what resources to deploy (other than the version)

Run the following commands:
```shell script
gosh create app your-app-name -g app-group
```
This will use the default template to setup your app. You can specify other templates as well to fit your needs. 
See [#App Templates]

## Versions

> TODO: the in-app command line help is much more up-to-date obviously, should we document commands here?

### List versions

Use `gosh list versions` to list specific version, you can filter on stages, releases, app groups and apps
See the CLI help for more information

_Example:_ List all *alpha* versions for all apps in group *my_app_group* 
```shell script
gosh list versions --stage alpha -g my_app_group 
```

### Update a version

Use `gosh update version`. See the CLI help for more information

*Example:* Update the stable version for app my-app
```shell
gosh update version --stage stable my-app 1.9.5
```

### List artifacts

Use `gosh list artifacts`

## Compiling the output

In order to compile the output, simply run

```shell script
./gitops compile
```

## Deploying

### Local

The following commands will run a single node K3S cluster purely in docker using K3D and will deploy a target to it

```shell script
brew update && brew install helm@3 k3d kubectl
k3d cluster create k3s-cluster -p9999:80@loadbalancer
./gitops compile -t tm-toggles-next
cd compiled/tm-toggles-next
kubectl apply --recursive -f .
```
