pipeline {
    agent any

    tools {
        go 'go-1.26'
    }

    stages {
        stage('Test') {
            steps {
                sh 'go build ./...'
                sh 'go vet ./...'
            }
        }

        stage('Release') {
            when {
                buildingTag()
            }
            environment {
                REGISTRY = credentials('harbor-bot-credentials')
            }
            steps {
                sh 'echo $REGISTRY_PSW | docker login junhyung.cloud -u $REGISTRY_USR --password-stdin'
                sh 'docker buildx inspect mcprobe-builder >/dev/null 2>&1 || docker buildx create --use --name mcprobe-builder'
                sh 'goreleaser release --clean'
            }
        }
    }
}
