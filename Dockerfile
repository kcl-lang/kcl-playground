FROM plutolang/pluto
WORKDIR /
COPY . .
RUN python3 -m pip install -r requirements
CMD ["pluto", "run"]
EXPOSE 9443
