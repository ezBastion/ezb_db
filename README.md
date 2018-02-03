# ezBastion DataBase access microservice.

**ezb_db** is a *rest to sql* microservice. It was use by ezBastion modules to interact
with the configuration DataBase.


## No drivers

For performance and memory foot print, each ezb_db embeds a native sql driver. Use
the binary corresponding to your sql engine (see git branch), ezb_db was compiled for:
- **MSSql** SQL Server 2005 or newer, Azure SQL Database https://github.com/denisenkom/go-mssqldb
- **Mysql** (4.1+), MariaDB, Percona Server, Google CloudSQL or Sphinx (2.2.3+) https://github.com/go-sql-driver/mysql)
- **Postgres** https://github.com/lib/pq
- **Sqlite** https://github.com/mattn/go-sqlite3

## No setup

You can install ezb_db as much as you need on Linux, Mac and Windows. As it uses
Rest (http) infrastructure, put it behind a load balancer and be elastic.

- **Copy** binary and json files somewhere (folder, docker, saas).
- **Configure** json file with you DataBase information and choose listen port.
- In a console, **start** the binary (see options below).


## ezb_db Options

- **init** Create or update DataBase tables and views.
- **install** Create Deamon/Service on the computer.
- **remove** Delete Deamon or service.
- **debug** Start ezb_db on the console and output all request made (verbose).
- **start** Start the Deamon or Service.
- **stop** Stop the Deamon or Service.
