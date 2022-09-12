# syntax=docker/dockerfile:1

# Build with bin/build-container-image

FROM docker.io/library/golang:1.19-bullseye

# As good practice, use a non-privileged user instead of root.
# Notes on useradd args:
# -l = do not add user to the lastlog and faillog databases; based on note on
#      the podman-build man page.
# -m = create home directory in '/home' so that the Go tool can use it (e.g. for
#      the build cache).
RUN groupadd -g 1000 shenmeci && \
    useradd -u 1000 -lmg shenmeci shenmeci
USER shenmeci:shenmeci
RUN mkdir -p /home/shenmeci/go/bin \
             /home/shenmeci/go/pkg \
             /home/shenmeci/src
WORKDIR /home/shenmeci/src

# Override GOPATH set in golang:1.19-bullseye.
ENV GOPATH=/home/shenmeci/go
ENV PATH=$GOPATH/bin:$PATH

# For debugging only.
# COPY --chown=shenmeci:shenmeci . ./
# RUN find . -printf '%u:%g %m %Z %y %h/%f\n' | sort -k 5

# Download CC-CEDICT and write config file.
COPY --chown=shenmeci:shenmeci download_dict.sh ./
RUN ./download_dict.sh && \
    echo '{"Http":{"Host":"0.0.0.0","Port":8080},"CedictPath":"dict/cedict_1_0_ts_utf-8_mdbg.txt.gz"}' > config.json

# Set flag to enable SQLite's Full-Text Search engine.
ENV GOFLAGS="-tags=sqlite_fts5"

# Copy only what's necessary to download and cache dependencies.
# Not useful when using RUN --mount=type=cache.
# COPY --chown=shenmeci:shenmeci go.mod go.sum ./
# RUN go mod download

# Copy build context.
COPY --chown=shenmeci:shenmeci . ./

# Sanity check and build.
RUN --mount=type=cache,uid=1000,gid=1000,id=shenmeci_go_mod_cache,target=/home/shenmeci/go/pkg/mod,z \
    --mount=type=cache,uid=1000,gid=1000,id=shenmeci_go_build_cache,target=/home/shenmeci/.cache/go-build,z \
    go test -vet=all ./... && go install

# For debugging only.
# RUN --mount=type=cache,uid=1000,gid=1000,id=shenmeci_go_mod_cache,target=/home/shenmeci/go/pkg/mod,z \
#     ls -laZ /home/shenmeci/go/pkg/mod
# RUN --mount=type=cache,uid=1000,gid=1000,id=shenmeci_go_build_cache,target=/home/shenmeci/.cache/go-build,z \
#     ls -laZ /home/shenmeci/.cache/go-build

EXPOSE 8080

CMD ["shenmeci"]
