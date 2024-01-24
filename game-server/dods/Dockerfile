# Use a base image with SteamCMD installed
FROM cm2network/steamcmd:latest

# Set the working directory
WORKDIR /home/steam/dods

# Install Day of Defeat: Source using SteamCMD
RUN ./steamcmd/steamcmd.sh +login anonymous +force_install_dir /home/steam/dods +app_update 232290 validate +quit

# Set the entry point for the container
ENTRYPOINT ["./srcds_run"]

# Default command, can be overridden -
CMD ["-game dod +map dod_avalanche +maxplayers 32"]
