FROM ubuntu:xenial

ENV PATH=/usr/lib/go-1.9/bin:$PATH

RUN \
  apt-get update && apt-get upgrade -q -y && \
  apt-get install -y --no-install-recommends golang-1.9 git make gcc libc-dev ca-certificates && \
  git clone --depth 1 --branch release/1.8 https://gitlab.neji.vm.tc/marconi/go-methereum-lite && \
  (cd go-methereum-lite && make geth) && \
  cp go-methereum-lite/build/bin/geth /geth && \
  apt-get remove -y golang-1.9 git make gcc libc-dev && apt autoremove -y && apt-get clean && \
  rm -rf /go-methereum-lite

EXPOSE 8545
EXPOSE 30303

ENTRYPOINT ["/geth"]
