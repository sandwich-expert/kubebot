// Describes the build process in Groovy

def err = null

try {
    node{
      stage 'Build and Test'
      //slackSend channel: '#tlp-sandwich-jenkins', color: 'good', message: "Build Started - ${env.JOB_NAME} ${env.BUILD_NUMBER} - ${env.BUILD_URL}"
      env.PATH = "${tool 'maven-3.3.9'}/bin:${env.PATH}"
      checkout scm

      if ("${env.JOB_NAME}".contains("master"))
      {
        // AMorton, push images to docker hub if building master, see TS-198
        sh 'build/build && build/publish'
      }
      else
      {
        sh 'build/build'
      }

      currentBuild.result = "SUCCESS"

      cleanupDockerImages()
      slackSend channel: '#tlp-sandwich-jenkins', color: 'good', message: ":smile: Build successful - ${env.JOB_NAME} ${env.BUILD_NUMBER} - ${env.BUILD_URL}"
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
          slackSend channel: '#tlp-sandwich', color: 'bad', message: ":cry: Build failed - ${env.JOB_NAME} ${env.BUILD_NUMBER} - ${env.BUILD_URL}"
        }
        slackSend channel: '#tlp-sandwich-jenkins', color: 'bad', message: ":cry: Build failed - ${env.JOB_NAME} ${env.BUILD_NUMBER} - ${env.BUILD_URL}"
    }

    /* Must re-throw exception to propagate error */
    if (err) {
        throw err
    }
}

def cleanupDockerImages(){
 sh returnStatus: true, script: 'docker rm $(docker ps --all --quiet)'
 sh returnStatus: true, script: 'docker rmi --force $(docker images -q thelastpickle/*)'
 sh returnStatus: true, script: 'docker rmi --force $(docker images --quiet --filter "dangling=true")'
}
