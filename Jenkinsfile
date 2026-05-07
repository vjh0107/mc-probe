pipeline {
    agent any

    stages {
        stage('Test') {
            agent {
                docker {
                    image 'golang:1.26'
                    reuseNode true
                }
            }
            steps {
                sh 'go build ./...'
                sh 'go vet ./...'
            }
        }

        stage('Login & Buildx') {
            when { buildingTag() }
            environment {
                REGISTRY = credentials('harbor-bot-credentials')
            }
            steps {
                sh 'echo "$REGISTRY_PSW" | docker login junhyung.cloud -u "$REGISTRY_USR" --password-stdin'
                sh 'docker buildx inspect mcprobe-builder >/dev/null 2>&1 || docker buildx create --use --name mcprobe-builder'
            }
        }

        stage('Release') {
            when { buildingTag() }
            agent {
                docker {
                    image 'goreleaser/goreleaser:v2.15.4'
                    args '-v /var/run/docker.sock:/var/run/docker.sock -v $HOME/.docker/config.json:/root/.docker/config.json:ro'
                    reuseNode true
                }
            }
            steps {
                sh 'goreleaser release --clean'
            }
        }
    }
}
