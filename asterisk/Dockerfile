FROM ubuntu:22.04

RUN export DEBIAN_FRONTEND=noninteractive && \
	apt-get update && \
	apt-get install asterisk -y

CMD [ "tail", "-F", "/dev/null" ]
