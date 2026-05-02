## Explain each stage of the multi-stage build

### Build stage

Sets up a full environment with all the heavy tools (the Go compiler, libraries, and your source code). Then compiles your code into a single, "ready-to-eat" file called a binary.

### Runtime Stage

Starts with a completely empty, tiny box (the Alpine image) that has almost nothing in it. The runtime stage reaches back to the Build Stage and grabs only that finished binary file. It ships a tiny, lightweight package that contains the app and nothing else. This makes it faster to move across the internet and gives hackers fewer "tools" to work with if they ever got inside.

## What does CGO_ENABLED=0 do? Why is it important?

When you set `CGO_ENABLED=0`, you are telling the Go compiler not to use any C code but to use Go’s internal versions of these tools instead. "Build a static binary." → A file that contains absolutely everything it needs to run inside one package.

Why is iz important? 

- If you build a binary with `CGO` enabled on a standard machine and move it to Alpine, it will look for `glibc`, fail to find it, and give you a confusing **"file not found"** error—even though the file is right there.
- `CGO_ENABLED=0` fixes this by making the binary independent of the host's libraries.
- If your binary is static because of `CGO_ENABLED=0`, you can use a Docker image called `FROM scratch`. This is an empty 0MB image. Your final container size will just be the size of your compiled code.
- By disabling CGO, you reduce the "attack surface." You aren't relying on external C libraries that might have their own security vulnerabilities

## Compare final image size vs single-stage build

| **Build Type** | **Final Image Size** | **What's part of it?** |
| --- | --- | --- |
| **Single-stage** | **~300MB+** | OS + Go Compiler + Source Code + Build Tools + Your Binary |
| **Multi-stage** | **~20MB** | Minimal OS + Your Binary |

If I’m not copmletely mistaken, this comparison is about the Multi-stage final image size vs. the Single-stage (final) image size? In this case, the Multi-stage approach would be much more efficient as it contains a lot less things in the build.

## Task 3- Testing

### Test operations with curl

#### Create 3+ products and list them afterwarts
Result: [{"id":1,"name":"test product 1","price":10},{"id":2,"name":"test product 2","price":20},{"id":3,"name":"test product 3","price":30},{"id":4,"name":"test product 4","price":40}]
#### Update price of product 1 and 2 and list them

Result:[{"id":1,"name":"test product 1","price":100},{"id":2,"name":"test product 2","price":200},{"id":3,"name":"test product 3","price":30},{"id":4,"name":"test product 4","price":40}]

#### Delete products 2 and three
Result: [{"id":1,"name":"test product 1","price":100},{"id":4,"name":"test product 4","price":40}]

### Verify data persistance
#### Running docker compose down&&docker compose up.
Result: WARN[0000] /Users/verenah/Documents/CICD/CI-CD-MCM-HAINY/docker-compose.yml: the attribute `version` is obsolete, it will be ignored, please remove it to avoid potential confusion 
[+] down 3/3
 ✔ Container ci-cd-mcm-hainy-api-1 Removed                                                                                                      0.2s
 ✔ Container ci-cd-mcm-hainy-db-1  Removed                                                                                                      0.2s
 ✔ Network ci-cd-mcm-hainy_default Removed                                                                                                      0.2s
WARN[0000] /Users/verenah/Documents/CICD/CI-CD-MCM-HAINY/docker-compose.yml: the attribute `version` is obsolete, it will be ignored, please remove it to avoid potential confusion 
[+] up 3/3
 ✔ Network ci-cd-mcm-hainy_default Created                                                                                                      0.0s
 ✔ Container ci-cd-mcm-hainy-db-1  Created                                                                                                      0.0s
 ✔ Container ci-cd-mcm-hainy-api-1 Created   
#### Open new terminal and list products.
Result: [{"id":1,"name":"test product 1","price":100},{"id":4,"name":"test product 4","price":40}]