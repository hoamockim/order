ARG GO_VERSION=1.16

FROM golang:${GO_VERSION}-alpine as builder

ARG SERVICE_TYPE
ENV SERVICE_NAME="order_$SERVICE_TYPE"

WORKDIR ./src/order/
COPY . .

RUN go mod tidy && \
    go build -o $SERVICE_NAME ./app/cmd/$SERVICE_TYPE/

RUN mkdir -p /app
RUN cp -p ./$SERVICE_NAME /app && \
    cp -p ./entrypoint.sh /app
RUN if [ "$SERVICE_TYPE" = "migration" ]; then cp -a $MIGRATION_RESOURCES /app/resources; fi

WORKDIR /app
RUN ls -a
RUN chmod +x ./entrypoint.sh
RUN chmod +x ./$SERVICE_NAME
CMD ["./entrypoint.sh"]
