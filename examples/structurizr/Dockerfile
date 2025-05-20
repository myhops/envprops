FROM ghcr.io/myhops/f12 AS f12-app

FROM structurizr/onpremises
COPY --from=f12-app /app/f12 /app/f12

ARG OTEL_JAVA_AGENT=https://github.com/open-telemetry/opentelemetry-java-instrumentation/releases/latest/download/opentelemetry-javaagent.jar
ADD ${OTEL_JAVA_AGENT} /extras/

ENV F12_NO_ENVPROPS=1
ENTRYPOINT ["/app/f12", "exec", "--" ]

CMD ["catalina.sh", "run"]

