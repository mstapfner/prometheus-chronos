# Chronos
Chronos (the personification of time) allows to manipulate raw Prometheus TSDB data.

Chronos allows : 
- read Prometheus TSDB Data 
- export Prometheus TSDB data directory to human-readable JSON files
- change the timestamps (to future or past) of Prometheus TSDB data and create new Prometheus TDSB
- copy the Prometheus TSDB data to fill a bigger timewindows (for proper performance testing of Prometheus)

Flags: 
 
 ```
-exportTSDB         Boolean - Only import the Prometheus TDSB Data 
-importDir          String  - Import directory for Prometheus TSDB data
-outputDir          String  - Output directory for Prometheus TSDB data
-jsonOutput         Boolean - Output the TSDB data as JSON files
-jsonOutputDir      String  - Output directory for JSON files 
-redateStart        UnixTimestamp for the starting date of the copied prometheus data
-redateEnd          Unix timestamp for the ending data of the copied prometheus data 

```



