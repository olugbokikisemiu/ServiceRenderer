# Getting started

Hello, nice having you aboard🤝🕺🏼

## Read System Design Document

[System Requiremt Document](./System_Requirement.md)
It contains a high-level overview of the system, what we expect the system to achieve

## New to Golang?

- [Go Intro](https://tour.golang.org)

- [Learning Go](https://medium.com/go-go-go/learning-go-golang-47127a796323)

- [8 Insights](https://medium.com/go-go-go/golang-8-insights-of-the-first-weeks-of-the-real-usage-f01290811b8b)

- [How to Write Go Code](https://golang.org/doc/code.html)

## Setting-up development environment

### Step 1. Install Golang

[Golang](https://golang.org/doc/install)

Make sure your Go version is at least 1.10.

### Step 2. Install Docker

We use **Community Edition**.

[Install Docker](https://docs.docker.com/engine/installation/)

### Step 5. Put sources in the correct folder

```
$ mkdir -p github.com.sleekservices

$ cd github.com.sleekservices

$ git clone https://github.com/olugbokikisemiu/ServiceRenderer.git

$ cd ServiceRenderer

$ go mod init
```

### Checking that everything works

#### 1. Build and run everything in docker

```
$ make env
$ docker ps (you should see the mongoDB image running)
```

### IDEs, Editors

Atom
Vim
Emacs

Gogland
Visual Code

## Workflows

### Dockerized workflow

The simplest way to get things running, with the least dependencies, is the dockerized workflow.
This keeps all of the building and running inside docker containers, so you don't need to do more setup.
For now only mongoDB is dockerized which you can start by running 

```
$ make env
```

### Choosing a task

Task can be picked on trello, an invite will be sent to the email you provide

### Implementing new task
Kindly follow this process when implementing your task

 - Checkout master
 - branch out a new branch from master and use this format for the branch name (your-name/taskname/version e.g (semiu/docs/0))
 - Add, Commit all neccessary implementation and push to github (git push --set-upstream  your-branch-name)
 - Vist github to create Pull Request on the branch you pushed

 ### Note: 
 Please don't make changes on master and push directly to github. Also, don't merge your PR to master until a review has been done on it and you receive a green light to go ahead.
