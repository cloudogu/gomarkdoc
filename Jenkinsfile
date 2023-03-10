#!groovy
@Library(['github.com/cloudogu/ces-build-lib@1.62.0'])
import com.cloudogu.ces.cesbuildlib.*

// Creating necessary git objects, object cannot be named 'git' as this conflicts with the method named 'git' from the library
git = new Git(this, "cesmarvin")
git.committerName = 'cesmarvin'
git.committerEmail = 'cesmarvin@cloudogu.com'
gitflow = new GitFlow(this, git)
github = new GitHub(this, git)
changelog = new Changelog(this)
Docker docker = new Docker(this)
gpg = new Gpg(this, docker)

// Configuration of repository
repositoryOwner = "cloudogu"
repositoryName = "gomarkdoc"
project = "github.com/${repositoryOwner}/${repositoryName}"

// Configuration of branches
productionReleaseBranch = "main"

node('docker') {
    timestamps {
        properties([
                // Keep only the last x builds to preserve space
                buildDiscarder(logRotator(numToKeepStr: '10')),
                // Don't run concurrent builds for a branch, because they use the same workspace directory
                disableConcurrentBuilds(),
        ])

        stage('Checkout') {
            checkout scm
            make 'clean'
        }

        stage('Build') {
            callInGoContainer{
                make 'compile'
            }
        }

        stageAutomaticRelease()
    }
}

void stageAutomaticRelease() {
    if (!gitflow.isReleaseBranch()) {
        return
    }

    String releaseVersion = git.getSimpleBranchName()

    stage('Build after Release') {
        git.checkout(releaseVersion)
        callInGoContainer{
            make 'clean compile checksum'
        }
    }
    
    stage('Finish Release') {
        gitflow.finishRelease(releaseVersion, productionReleaseBranch)
    }

    stage('Sign after Release'){
        gpg.createSignature()
    }

    stage('Add Github-Release') {
        releaseId=github.createReleaseWithChangelog(releaseVersion, changelog, productionReleaseBranch)
        github.addReleaseAsset("${releaseId}", "target/gomarkdoc")
        github.addReleaseAsset("${releaseId}", "target/gomarkdoc.sha256sum")
        github.addReleaseAsset("${releaseId}", "target/gomarkdoc.sha256sum.asc")
    }
}

void make(String makeArgs) {
    sh "make ${makeArgs}"
}

void callInGoContainer(Closure closure) {
    new Docker(this)
            .image('golang:1.19.7')
            .mountJenkinsUser()
            .inside("--volume ${WORKSPACE}:/go/src/${project} -w /go/src/${project}")
                    {
                        closure.call()
                    }
}
