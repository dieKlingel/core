FROM ubuntu:22.04

RUN apt-get update && \
	apt-get -y install \
	gawk \
	wget \
	git \
	diffstat \
	unzip \
	texinfo \
	build-essential \
	chrpath \
	socat \
	cpio \
	python3 \
	python3-pip \
	python3-pexpect \
	xz-utils \
	debianutils \
	iputils-ping \
	libsdl1.2-dev \
	xterm \
	tar \
	locales \
	file \
	lz4 \
	zstd

RUN rm /bin/sh && ln -s bash /bin/sh
RUN locale-gen en_US.UTF-8 && \
	update-locale LC_ALL=en_US.UTF-8 LANG=en_US.UTF-8

RUN git clone -b dunfell git://git.yoctoproject.org/poky

RUN groupadd -g 1001 dieklingel && \
	useradd -g 1001 -m -s /bin/bash -u 1001 dieklingel
RUN chown -R dieklingel:dieklingel /poky
USER dieklingel
WORKDIR /poky
RUN git clone -b dunfell git://git.openembedded.org/meta-openembedded &&\
	git clone -b dunfell git://git.yoctoproject.org/meta-raspberrypi

RUN source ./oe-init-build-env && \
	bitbake-layers add-layer /poky/meta && \
	bitbake-layers add-layer /poky/meta-poky && \
	bitbake-layers add-layer /poky/meta-yocto-bsp && \
	bitbake-layers add-layer /poky/meta-raspberrypi && \
	bitbake-layers add-layer /poky/meta-openembedded/meta-oe && \
	bitbake-layers add-layer /poky/meta-openembedded/meta-python && \
	bitbake-layers add-layer /poky/meta-openembedded/meta-networking

ENTRYPOINT [ "tail" , "-F", "/dev/null" ]
