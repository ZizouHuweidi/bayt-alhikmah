FROM devopsfaith/krakend:2.7.0 as builder

# Copy configuration file
COPY krakend.json /etc/krakend/krakend.json

# Validate configuration
RUN krakend check -c /etc/krakend/krakend.json

# Production stage
FROM devopsfaith/krakend:2.7.0

# Copy validated configuration
COPY --from=builder /etc/krakend/krakend.json /etc/krakend/krakend.json

# Expose port
EXPOSE 8080

# Run KrakenD
CMD ["krakend", "run", "-c", "/etc/krakend/krakend.json"]