FROM scratch

ARG releaseName=layover.amd64.linux.tar.gz
RUN wget https://github.com/kjcodeacct/layover/releases/download/v1.0.0/${releaseName} \
&& tar zxvpf ${releaseName} && rm ${releaseName}

# Copy our static executable.
COPY layover /

# Run the layover binary.
ENTRYPOINT ["/layover"]

