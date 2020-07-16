// Describes the build process in Groovy

def err = null

try {
    node{
      stage 'Build and Test'
      env.PATH = "${tool 'maven-3.3.9'}/bin:${env.PATH}"
      checkout scm

      DOCKER_TAG = ""
      if ("${env.JOB_NAME}".contains("master"))
      {
        hash = sh(returnStdout: true, script: 'git log -n 1 --pretty=format:"%H"')
        DOCKER_TAG = hash.substring(0, 7)
      }
      else
      {
        DOCKER_TAG = branchName + "-SNAPSHOT"
      }
      echo "Building docker tag ${DOCKER_TAG}"

      stage 'Build and Test'
      buildPublishImage(env.JOB_NAME, DOCKER_TAG)

      currentBuild.result = "SUCCESS"

      cleanupDockerImages()
    }
} catch (caughtError) {
    err = caughtError
    currentBuild.result = "FAILURE"
} finally {
    if (currentBuild.result != "SUCCESS") {
        // Send slack notifications for failed or unstable builds.
        // currentBuild.result must be non-null for this step to work.
        if ("${env.JOB_NAME}".contains("master") || "${env.JOB_NAME}".contains("develop"))
        {
          slackSend channel: '#vector', color: 'bad', message: ":cry: Build failed - ${env.JOB_NAME} ${env.BUILD_NUMBER} - ${env.BUILD_URL}"
        }
    }

    /* Must re-throw exception to propagate error */
    if (err) {
        throw err
    }
}

def cleanupDockerImages()
{
  sh returnStatus: true, script: 'docker rm $(docker ps --all --quiet)'
  sh returnStatus: true, script: 'docker rmi --force $(docker images -q thelastpickle/*)'
  sh returnStatus: true, script: 'docker rmi --force $(docker images --quiet --filter "dangling=true")'
}

def buildPublishImage(jobName, dockerTag)
{
  configFileProvider([configFile(fileId: '721b2539-35ba-4612-81b2-f73d39ca08ad', variable: 'DOCKER_CONFIG_JSON')])
  {
    sh 'mkdir -p ~/.docker && cp $DOCKER_CONFIG_JSON ~/.docker/config.json'
    withCredentials([usernamePassword(credentialsId: '8daa33e8-c1cb-4e78-babd-2e1a06de1c90', passwordVariable: 'AWS_SECRET_ACCESS_KEY', usernameVariable: 'AWS_ACCESS_KEY_ID')])
    {
      sh "DOCKER_TAG=${dockerTag} build/build && DOCKER_TAG=${dockerTag} build/publish"
    }
  }
}
