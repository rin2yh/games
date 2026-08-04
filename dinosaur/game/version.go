package game

// Version is rewritten by tagpr to match the git tag (see .tagpr), so it
// names the last release, not the working tree. The tag covers the whole
// repository, but dinosaur/ is the only package here that is importable
// from outside it, so this is where a version number carries meaning. It is
// here for the firmware to report; nothing reads it yet.
const Version = "0.0.1"
