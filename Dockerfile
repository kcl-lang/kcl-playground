FROM plutolang/pluto
WORKDIR /
COPY . .
RUN npm install
RUN cd web && npm run build
RUN python3 -m pip install -U -r ./requirements.txt
CMD ["pluto", "run"]
EXPOSE 9443
