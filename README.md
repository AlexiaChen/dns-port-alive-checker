# DNS Port Check Tool

This tool checks if the DNS port (53) is open for a list of IP addresses provided in a CSV file. The results are written to an output CSV file.

## Usage

- -input: Path to the input CSV file (default: ip_list.csv)
- -output: Path to the output CSV file (default: output.csv)

##  Example:

```bash
./dns_port_check -input=ip_list.csv -output=output.csv
```

## Input File Format


The input CSV file should have the following format:

```csv
IP,Server ID
192.168.1.1,server1
192.168.1.2,server2
...
```
   