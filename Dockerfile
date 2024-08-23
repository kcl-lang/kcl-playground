FROM plutolang/pluto
WORKDIR /
COPY . .
# Install Pluto and dev dependencies
RUN npm install
# Build frontend web application
RUN cd web && npm install && npm run build
# Install backend dependencies
RUN python3 -m pip install -U -r ./requirements.txt
# Run
CMD ["pluto", "run"]
EXPOSE 9443
