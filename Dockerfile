FROM docker.bintray.io/jfrog/artifactory-pro
ENV LICENSE_DATA LICENSE_DATA
ENTRYPOINT ["sh", "-c", "echo ${LICENSE_DATA?No license data supplied} > /var/opt/jfrog/artifactory/etc/artifactory/artifactory.lic && /entrypoint-artifactory.sh"]
