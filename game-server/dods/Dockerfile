FROM lacledeslan/steamcmd:linux as dods-builder

# Download Day of Defeat: Source
RUN /app/steamcmd.sh +force_install_dir /output +login anonymous +app_update 232290 validate +quit;

#=======================================================================
FROM debian:bookworm-slim

ARG BUILDNODE=unspecified
ARG SOURCE_COMMIT=unspecified

HEALTHCHECK NONE

RUN dpkg --add-architecture i386 && \
    apt-get update \
    && apt-get install -y \
        ca-certificates lib32gcc-s1 libncurses5:i386 libsdl2-2.0-0:i386 libstdc++6 libstdc++6:i386 locales locales-all tmux && \
    apt-get clean && \
    rm -rf /tmp/* /var/lib/apt/lists/* /var/tmp/*;

ENV LANG=en_US.UTF-8 LANGUAGE=en_US.UTF-8 LC_ALL=en_US.UTF-8

# Set up Enviornment
RUN useradd --home /app --gid root --system DODS && \
    mkdir -p /app/.steam/sdk32 && \
    chown DODS:root -R /app;

COPY --chown=DODS:root --from=dods-builder /output /app

USER DODS

WORKDIR /app

CMD ["/bin/bash"]

ONBUILD USER root