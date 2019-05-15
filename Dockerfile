FROM busybox
COPY ./resource/zoneinfo.zip /usr/local/go/lib/time/zoneinfo.zip
COPY resource/ /resource
COPY apidoc/ /apidoc/
RUN ls -la /apidoc/*
CMD ["/resource/main", "-env", "prod", "-config", "/resource/cfg.yml"]
EXPOSE 80
