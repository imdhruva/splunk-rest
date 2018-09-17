FROM ubuntu:latest

RUN apt-get update -y \
    && apt install wget -y \
    && wget -O splunk-7.1.2-a0c72a66db66-Linux-x86_64.tgz 'https://www.splunk.com/bin/splunk/DownloadActivityServlet?architecture=x86_64&platform=linux&version=7.1.2&product=splunk&filename=splunk-7.1.2-a0c72a66db66-Linux-x86_64.tgz&wget=true' \
    && useradd splunk \
    && mkdir /opt/splunk \
    && tar -xzvf splunk-7.1.2-a0c72a66db66-Linux-x86_64.tgz -C /opt/ \
    && chown -R splunk:splunk /opt/splunk/

COPY --chown=splunk:splunk ./opt/splunk/etc/splunk-launch.conf /opt/splunk/etc/
COPY --chown=splunk:splunk ./etc/system/local/user-seed.conf /etc/system/local/

RUN su - splunk \
    && /opt/splunk/bin/splunk start --answer-yes --no-prompt --accept-license
    
RUN /opt/splunk/bin/splunk enable boot-start

ENTRYPOINT /etc/init.d/splunk start && tail -f /opt/splunk/var/log/splunk/splunkd.log

EXPOSE 8000
