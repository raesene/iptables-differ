# IPTables Differ

A very quite and dirty utility to compare two sets of IPTables rules

## Usage

To generate the input files, use the following iptables commands:
  Before changes: iptables-save > rules-before.txt
  After changes:  iptables-save > rules-after.txt

Usage:
  iptables-diff -before <before-file> -after <after-file>

Example:
  iptables-diff -before rules-before.txt -after rules-after.txt

Output will be color-coded:
  - Red for removed rules
  - Green for added rules
  - Yellow for table changes
