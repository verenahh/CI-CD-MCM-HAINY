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