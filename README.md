# ftp2p

Simple peer-to-peer file transfer protocol inspired by the BitTorrent protocol, but instead of custom TCP protocols everything is built on top of QUIC (and therefore can be protected by TLS easily).

## How does it work?

There are three types of actors in a functioning `ftp2p` environment. First of all, there are the **taverns**. These are rendevouz servers, which introduce network peers to each other. The first type of peer is the **seeder**. They provide their bandwidth and storage for others to use and download file chunks from. Then there is the **fetcher**. Its sole purpose is to contact **taverns**, ask around for peers for a given file hash and contact these **seeders**, fetching file chunks and assembling the final file at the end.

## What does (not) work yet?

* [x] Simple file sharing using tavern, fetcher and seeder
* [x] Announcement timeouts
* [x] Better chunk download scheduling
* [x] Add support to resume fetches by checking for already downloaded chunks (and verifying them)
* [x] Use UDP-based QUIC as transport layer with Buf Connect on top of it (with this comes UDP hole punching)
* [ ] Encryption of tavern traffic using Let's Encrypt certificates
* [ ] Use distributed hash tables for decentralized peer discovery
* [ ] Split `.tracker` into public- and private trackers, allowing for file encryption
* [ ] Add tests for fetcher, seeder, tavern and tracker