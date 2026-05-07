FROM scratch
COPY mcprobe /mcprobe
ENTRYPOINT ["/mcprobe"]
