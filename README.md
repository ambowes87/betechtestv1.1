**Tech Test - Alex Bowes**

**Prerequisities**

- GCC enabled on toolchain

*Build Instructions (Windows)*

- Build the binary using `go build -o techtest.exe src/`
- The executable relies on a sqlite database named userstore.db in the same directory (my laptop doesn't support docker sadly)

*Running Instructions (Windows)*

- Launch the binary using `./techtest.exe -port 9090` to start the server
- Check the server is running by calling `http://localhost:9090/health` in a browser, you should get 200 returned and a blank page
