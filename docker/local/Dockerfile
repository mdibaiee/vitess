ARG bootstrap_version=33
ARG image="vitess/bootstrap:${bootstrap_version}-common"

FROM "${image}"

RUN apt-get update
RUN apt-get install -y sudo curl vim jq

# Install dependencies
COPY docker/utils/install_dependencies.sh /vt/dist/install_dependencies.sh
RUN /vt/dist/install_dependencies.sh mysql80

COPY docker/local/install_local_dependencies.sh /vt/dist/install_local_dependencies.sh
RUN /vt/dist/install_local_dependencies.sh
RUN echo "source /vt/common/env.sh" >> /etc/bash.bashrc

# Allows some docker builds to disable CGO
ARG CGO_ENABLED=0

# Re-copy sources from working tree.
COPY --chown=vitess:vitess . /vt/src/github.com/mdibaiee/vitess

# Build and install Vitess in a temporary output directory.
USER vitess

WORKDIR /vt/src/github.com/mdibaiee/vitess
RUN make install PREFIX=/vt/install

ENV VTROOT /vt/src/github.com/mdibaiee/vitess
ENV VTDATAROOT /vt/vtdataroot
ENV PATH $VTROOT/bin:$PATH
ENV PATH="/var/opt/etcd:${PATH}"

RUN mkdir /vt/local
COPY examples/local /vt/local

# Copy the vtadmin web app to the correct location and npm install
COPY --chown=vitess:vitess web /web
RUN npm install /web/vtadmin
RUN /web/vtadmin/build.sh

RUN mkdir /vt/common
COPY examples/common /vt/common

CMD cd /vt/local && ./101_initial_cluster.sh && /bin/bash
