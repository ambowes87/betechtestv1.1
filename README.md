**Tech Test - Alex Bowes**

**Prerequisities**

- GCC enabled on toolchain

*Build Instructions (Windows)*

- Clone the repo locally
- Build the binary using `go build -o techtest.exe src/` from the root of the project
- The executable relies on a sqlite database named userstore.db in the same directory (my laptop doesn't support docker sadly). This is included in the project, yes I know this would never be done in a production setting :)

*Running Instructions (Windows)*

- Launch the binary using `./techtest.exe -port 9090` to start the server
- Check the server is running by calling `http://localhost:9090/health` in a browser, you should get 200 returned and a blank page
