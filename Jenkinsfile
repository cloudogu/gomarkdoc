#!groovy
@Library(['github.com/cloudogu/ces-build-lib@1.62.0'])
import com.cloudogu.ces.cesbuildlib.*

// Creating necessary git objects, object cannot be named 'git' as this conflicts with the method named 'git' from the library
gitWrapper = new Git(this, "cesmarvin")
gitWrapper.committerName = 'cesmarvin'
gitWrapper.committerEmail = 'cesmarvin@cloudogu.com'
gitflow = new GitFlow(this, gitWrapper)
github = new GitHub(this, gitWrapper)
changelog = new Changelog(this)
Gpg gpg = new Gpg(this, docker)

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
            make 'compile'
        }

        stageAutomaticRelease()
    }
}

void stageAutomaticRelease() {
    if (!gitflow.isReleaseBranch()) {
        return
    }

    String releaseVersion = gitWrapper.getSimpleBranchName()

    stage('Finish Release') {
        gitflow.finishRelease(releaseVersion, productionReleaseBranch)
    }

    stage('Build after Release') {
        git.checkout(releaseVersion)
        make 'clean compile checksum'
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
