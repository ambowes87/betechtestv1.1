**Tech Test - Alex Bowes**

**Prerequisities**

- GCC enabled on toolchain (sorry annoying I know, Docker doesn't work on my laptop so I couldn't do anything containerised)

**Build Instructions (bash shell in Windows)**

- Clone the repo locally
- Build the binary using `go build -o techtest.exe ./src/` from the root of the project
- The executable relies on a sqlite database named userstore.db in the same directory (my laptop doesn't support docker sadly). This is included in the project, yes I know this would never be done in a production setting :)

**Running Instructions on Windows**

- Launch the binary using `./techtest.exe -port 9090` to start the server
- Check the server is running by calling `http://localhost:9090/health` in a browser, you should get 200 status returned with no body

**Notes / Future Work**

- Unit tests are very light (in fact only one), usually I'd have much more coverage but ran out of time - I chose to do a http server one as that can be a bit more complex
- The pagination on Get users didn't work, this was probably due to an issue with my SQL query - I left the code in to show I know *how* it's done but this solution isn't working properly
- I've never used a pub/sub system before. I attempted one that just basically notifies a channel on a user update - but nothing is actually picking this up
