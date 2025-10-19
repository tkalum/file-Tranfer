# Local Network File Transfer

Simple command-line file transfer tool that uses mDNS for device discovery on local networks.

## Features
- Automatic peer discovery using mDNS/Bonjour
- Real-time transfer progress display
- No configuration needed - just run and transfer
- Supports files of any size
- Simple command-line interface

## Prerequisites
- Go 1.25 or higher
- Local network that allows mDNS traffic
- Firewall access for ports 24242 (discovery) and 24243 (transfer)

## Installation
```bash
git clone github.com/tkalum/file-Trasfer.git
cd fileTransfer
go build ./cmd/file-transfer
```

## Usage

### Send a file
```bash
./file-transfer -s filename.txt
```
The sender will:
1. Announce its presence on the network
2. Wait for a receiver to connect
3. Display transfer progress

### Receive a file
```bash
./file-transfer
```
The receiver will:
1. Search for available senders on the network
2. Connect to the first discovered sender
3. Save received file with "receive_" prefix
4. Show download progress

## Technical Details
- Discovery Protocol: mDNS (Bonjour)
- Service Name: _filetransfer._tcp
- Discovery Port: 24242
- Transfer Port: 24243
- Transfer Buffer: 64KB chunks

## Troubleshooting
- Ensure both devices are on the same local network
- Check firewall settings if devices can't discover each other
- Look for files with "receive_" prefix in the working directory
- Make sure ports 24242 and 24243 are not in use

## Notes
- Files are saved in the current working directory
- Transfer speed depends on local network performance
- Large files are supported but may take longer to transfer
- Only one file can be transferred at a time
