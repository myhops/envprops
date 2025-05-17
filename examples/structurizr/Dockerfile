FROM ghcr.io/myhops/f12 AS f12-app

FROM structurizr/onpremises
COPY --from=f12-app /app/f12 /app/f12

ENV F12_NO_ENVPROPS=1
ENTRYPOINT ["/app/f12", "exec", "--" ]

CMD ["catalina.sh", "run"]

