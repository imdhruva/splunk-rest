FROM ubuntu:latest

RUN apt-get update -y \
    && apt-get install wget -y \
    && wget -O splunk-7.1.2-a0c72a66db66-Linux-x86_64.tgz 'https://www.splunk.com/bin/splunk/DownloadActivityServlet?architecture=x86_64&platform=linux&version=7.1.2&product=splunk&filename=splunk-7.1.2-a0c72a66db66-Linux-x86_64.tgz&wget=true' \
    && useradd -s /bin/bash -m splunk \
    && tar -xzvf splunk-7.1.2-a0c72a66db66-Linux-x86_64.tgz -C /opt/

COPY --chown=splunk:splunk ./opt/splunk/etc/splunk-launch.conf /opt/splunk/etc/
COPY --chown=splunk:splunk ./etc/system/local/user-seed.conf /etc/system/local/


RUN /opt/splunk/bin/splunk start --answer-yes --no-prompt --accept-license \
    && /opt/splunk/bin/splunk enable boot-start -user splunk


COPY --chown=root:root ./etc/init.d/splunk /etc/init.d/

RUN chown -R splunk:splunk /opt/splunk/

# Define mountable directories.
VOLUME ["/opt/splunk/var/"]

# Define working directory.
WORKDIR /opt/splunk

ENTRYPOINT /etc/init.d/splunk start && tail -f /opt/splunk/var/log/splunk/splunkd.log

EXPOSE 8089
