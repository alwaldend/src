FROM ${IMAGE} as git-web
RUN apk add cgit git nginx dumb-init
COPY cgitrc /etc/cgitrc
ENTRYPOINT ["/usr/bin/dumb-init", "--"]
CMD ["nginx"]

ARG IMAGE BASE
FROM ${IMAGE} as git-server
RUN apk add git openssh shadow dumb-init && \
    useradd -s /usr/bin/git-shell -p "*" -m git && \
    mkdir -p /srv/git /home/git/.ssh && \
    chmod 700 /home/git/.ssh /srv/git && \
    chown -R git:git /srv/git /home/git && \
    touch /etc/ssh/revoked_keys
COPY sshd_config /etc/ssh/sshd_config
COPY entrypoint /usr/bin
EXPOSE 22
ENTRYPOINT ["/usr/bin/dumb-init", "--"]
CMD ["/usr/bin/entrypoint"]
