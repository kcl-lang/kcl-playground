run:
	cd web && npm run build
	pluto run

fmt:
	python3 -m pip install black && python3 -m black .

deps:
	python3 -m pip install -r requirements
	npm install
	npm install -g pluto
