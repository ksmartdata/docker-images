FROM --platform=$TARGETPLATFORM maven:3.6.3-jdk-8-openj9 AS builder

RUN apt-get update && apt-get install -y unzip

RUN curl -Lo rocketmq-exporter.zip https://github.com/apache/rocketmq-exporter/archive/2f2feb752b2a1a4b3fce300b5a697b719a16cca8.zip

RUN unzip rocketmq-exporter.zip

RUN cd rocketmq-exporter && mvn clean package -Dmaven.test.skip=truedocker

FROM --platform=$TARGETPLATFORM openjdk:8-jre-alpine

RUN echo "Asia/Shanghai" > /etc/timezone

COPY --from=builder /rocketmq-exporter/target/rocketmq-exporter-0.0.3-SNAPSHOT-exec.jar /rocketmq-exporter.jar

EXPOSE 5557

ENTRYPOINT ["java","-jar","rocketmq-exporter.jar"]