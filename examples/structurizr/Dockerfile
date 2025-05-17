# Ensure to download opentelemetry-javaagent.jar
# curl -O -L https://github.com/open-telemetry/opentelemetry-java-instrumentation/releases/latest/download/opentelemetry-javaagent.jar
# or use ADD

FROM ghcr.io/myhops/f12:latest AS f12-app

FROM structurizr/onpremises
COPY --from=f12-app /app/f12 /app/f12

ARG javaagent="https://github.com/open-telemetry/opentelemetry-java-instrumentation/releases/latest/download/opentelemetry-javaagent.jar"
# Extra
ADD ${javaagent} /extras/opentelemetry-javaagent.jar

ENV F12_NO_ENVPROPS=1
ENTRYPOINT ["/app/f12", "exec", "--" ]
CMD ["catalina.sh", "run"]