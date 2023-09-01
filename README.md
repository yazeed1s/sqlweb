# sqlweb

## About The Tool
This is a DB web client that enables users to seamlessly connect to relational databases via a user-friendly web interface. It offers a comprehensive set of features designed to enhance your database management experience.

## ✨ Features
- Connect to relational databases, whether hosted locally or remotely.
- A built-in SQL editor.
- Export table data in CSV or JSON formats.
- Generate raw SQL for database schema objects.
- Editable table cells
- Save connection info locally
- Data visualization
- Relatively small and efficient executable < 15MB
- Memory usage < 25MB
- Rely on minimal external dependencies, mostly utilizing the standard library.

## ⭐️ Screenshots
<p align="left">
  <img src="https://github.com/Yazeed1s/sqlweb/blob/master/screenshot/dbit1.png" width="1000">
</p>
<p align="left">
  <img src="https://github.com/Yazeed1s/sqlweb/blob/master/screenshot/dbit2.png" width="1000">
</p>
<p align="left">
  <img src="https://github.com/Yazeed1s/sqlweb/blob/master/screenshot/dbit3.png" width="1000">
</p>
<p align="left">
  <img src="https://github.com/Yazeed1s/sqlweb/blob/master/screenshot/dbit4.png" width="1000">
</p>

## 📦 Installation

### You can grab your executable from the [releases](https://github.com/Yazeed1s/sqlweb/releases) page

### Using the install script (mac/linux)
```bash
curl -s https://raw.githubusercontent.com/Yazeed1s/sqlweb/master/install.sh | sudo bash
```

### build and install from source
#### Dependencies:
	1- go 
	2- vite
	3- yarn
	4- git

1- clone the repo
```bash
git clone https://github.com/Yazeed1s/sqlweb.git
```

2- cd into sqlweb/ui and install the ui dependencies
```bash
cd sqlweb/ui && yarn install
```

3- go back to the parent dir and run make build
```bash
cd .. && make build
```

4- install the binary
```bash
sudo make install
```

TODO: brew, yay, windows

## 🚀 Usage

- Just run `sqlweb` from the terminal, then open the browser on localhost:3000
```bash
$ sqlweb
2023/09/01 13:41:04 Listening...3000
```
- There are some flags that can be passed to `sqlweb`
```bash
$ sqlweb -h
  Help information:
  USAGE: sqlweb [OPTION]
     OPTION:
	   -p <port>   	Set the port number (default: 3000)
	   -h          	Display help information
	   -v          	Display version
```

## ✅  TODO:
- [x] Add support for MySQL
- [x] Add support for PostgreSQL
- [ ] Add support for SQLite 
- [ ] Add support for MariaDB
- [x] Editable cells
- [x] Table pagination
- [x] Display columns info (field name, key, type, referenced column, referenced table, constraints name)
- [x] Export table to csv
- [x] Export table to json
- [x] Export schema objects to raw sql (for those ORM users who didn't design/write the schema)
- [ ] Data visualization
- [ ] Support multiple sessions
- [ ] Add `-o` flag to open up the browser on localhost:port
- [ ] Display database constraints, size, number of tables... etc
- [ ] Manage/add/remove users and their permissions
- [ ] Query history

## 🔥 Contributing

If you would like to add a feature or to fix a bug please feel free to send a PR.
