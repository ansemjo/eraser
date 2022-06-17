# eraser

Stupid simple Go binary to overwrite block devices and "wipe them clean".

The [idea][idea] here is that if you want to overwrite your disk with random data,
the `/dev/urandom` endpoint is pretty slow in supplying enough data quickly.
But because almost any modern computer has some form of hardware acceleration
for the AES algorithm, you can just encrypt a stream of zeroes with a random
key and you'll get pretty decent randomness. The need for 35 Gutman passes is
long gone, so this one pass ought the be enough, usually.

**NOTE:** This only reliably works on spinning disks, **not** flash disks like
SSDs! Use ATA Secure Erase in that case, which deletes the MEK on self-encrypting
drives and renders all data useless instantly. Some harddisks also have instant
secure erase (ISE) – you should prefer that.

[idea]: https://wiki.archlinux.org/index.php/Securely_wipe_disk#Random_data

# INSTALLATION

    go get github.com/ansemjo/eraser

# USAGE

    eraser { -rand | -zero } [-direct] [-note] blockdev

Use `-rand` for the encrypted zerostream described above or `-zero` to just
use zeroes instead.

The `-note` flag writes a little note with a timestamp to the first 32 bytes
of `blockdev` after successful deletion. You can then `head -1 blockdev` and
see when the disk was deleted later.

With `-direct` the disk is opened with `O_DIRECT`, which bypasses most caches and gives a more realistic speed.

The progress spinner calculates the estimated remaining time based on the average speed of the bytes *written so far*, as I've found that the current average speed is a very bad measure for accurate estimations.

# DISCLAIMER

I'm not a cryptographer. This is just a small utility I like to use because
I was fed up with copy-pasting a long `openssl enc` command. Don't trust me
with your data.
