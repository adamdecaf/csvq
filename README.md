# csvq

Extract first_name and last_name columns (in that order). Sort results.
```
csvq -keep first_name,last_name ~/Downloads/report.csv | sort -u
```

Change delimiter used in `report.csv`.
```
csvq -d';' user_id,dob,email ~/Downloads/report.csv
```

Output CSV columns in a table.
```
csvq -keep first_name,last_name -format table
```

Combine multiple files.
```
csvq -keep user_id,email ~/Downloads/report1.csv ~/Downloads/report2.csv
```

## License

Apache License 2.0 - See [LICENSE](LICENSE) for details.
