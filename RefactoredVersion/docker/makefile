postgres:
	@ docker run --name auth-postgres \
       -e POSTGRES_USER=authuser \
       -e POSTGRES_PASSWORD=authpass \
       -e POSTGRES_DB=authdb \
       -p 5432:5432 \
       -d postgres:latest \
       -c 'shared_preload_libraries=uuid-ossp'


docker:
	@dco