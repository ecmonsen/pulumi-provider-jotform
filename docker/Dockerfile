FROM pulumi/pulumi

RUN mkdir -p /tmp/pc && \
    cd /tmp/pc && \
    wget https://github.com/pulumi/pulumictl/releases/download/v0.0.43/pulumictl-v0.0.43-linux-amd64.tar.gz && \
    tar xzvf *.gz && \
    mv pulumictl /usr/local/bin

RUN export PULUMI_CONFIG_PASSPHRASE= && \
    pulumi login file:///root && \
    mkdir -p /root/bin && \
    echo 'export PATH=$PATH:$HOME/bin' >> /root/.bashrc \
    echo 'export PULUMI_CONFIG_PASSPHRASE=' >> /root/.bashrc

COPY run-jotform-example /root/bin/


