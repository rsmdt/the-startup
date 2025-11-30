# Reference: Profiling Tools

## Overview

Comprehensive reference for profiling and performance measurement tools organized by language, platform, and profiling type.

## Quick Reference

| Category | Tool | Best For |
|----------|------|----------|
| Node.js CPU | `--prof`, `clinic.js` | V8 profiling, flame graphs |
| Node.js Memory | `--inspect`, `memwatch` | Heap snapshots, leak detection |
| Python CPU | `cProfile`, `py-spy` | Function timing, live profiling |
| Python Memory | `memory_profiler`, `tracemalloc` | Allocation tracking |
| Java/JVM | JFR, async-profiler | Low-overhead production profiling |
| Go | `pprof` | Built-in profiling |
| Browser | DevTools, Lighthouse | Client-side performance |
| Database | EXPLAIN, pg_stat | Query optimization |
| System | `perf`, `htop`, `iotop` | Resource utilization |

## Node.js Profiling

### CPU Profiling

**Built-in V8 Profiler**
```bash
# Generate V8 profiling log
node --prof app.js

# Process the log into readable format
node --prof-process isolate-*.log > profile.txt
```

**Clinic.js Suite**
```bash
# Install
npm install -g clinic

# Doctor - overall health check
clinic doctor -- node app.js

# Flame - flame graphs
clinic flame -- node app.js

# Bubbleprof - async operations
clinic bubbleprof -- node app.js
```

**0x - Flame Graphs**
```bash
npm install -g 0x
0x app.js
# Opens flame graph in browser
```

### Memory Profiling

**Chrome DevTools (via --inspect)**
```bash
node --inspect app.js
# Open chrome://inspect in Chrome
# Take heap snapshots, record allocations
```

**heapdump Module**
```javascript
const heapdump = require('heapdump');
heapdump.writeSnapshot('/tmp/heap-' + Date.now() + '.heapsnapshot');
```

**memwatch-next**
```javascript
const memwatch = require('memwatch-next');

memwatch.on('leak', (info) => {
  console.log('Memory leak detected:', info);
});

memwatch.on('stats', (stats) => {
  console.log('GC stats:', stats);
});
```

### APM Solutions

| Tool | Description |
|------|-------------|
| New Relic | Full-stack APM with Node.js agent |
| Datadog | APM with distributed tracing |
| Dynatrace | AI-powered APM |
| AppSignal | Lightweight APM for Node.js |

## Python Profiling

### CPU Profiling

**cProfile (Built-in)**
```bash
# Profile entire script
python -m cProfile -o output.prof script.py

# Analyze results
python -m pstats output.prof
```

```python
import cProfile
import pstats

profiler = cProfile.Profile()
profiler.enable()
# ... code to profile ...
profiler.disable()
stats = pstats.Stats(profiler).sort_stats('cumtime')
stats.print_stats(10)  # Top 10 functions
```

**py-spy (Low Overhead)**
```bash
pip install py-spy

# Sample running process
py-spy top --pid 12345

# Generate flame graph
py-spy record -o profile.svg --pid 12345

# Profile script directly
py-spy record -o profile.svg -- python script.py
```

**line_profiler**
```python
# Decorate functions to profile
@profile
def slow_function():
    pass

# Run with kernprof
# kernprof -l -v script.py
```

### Memory Profiling

**memory_profiler**
```python
from memory_profiler import profile

@profile
def memory_intensive():
    large_list = [i for i in range(1000000)]
    return large_list

# Run: python -m memory_profiler script.py
```

**tracemalloc (Built-in)**
```python
import tracemalloc

tracemalloc.start()
# ... code to profile ...
snapshot = tracemalloc.take_snapshot()
top_stats = snapshot.statistics('lineno')

for stat in top_stats[:10]:
    print(stat)
```

**objgraph**
```python
import objgraph

# Show most common types
objgraph.show_most_common_types()

# Find reference chains to objects
objgraph.show_backrefs(objgraph.by_type('MyClass')[0])
```

## Java/JVM Profiling

### Java Flight Recorder (JFR)

```bash
# Start recording
java -XX:StartFlightRecording=duration=60s,filename=recording.jfr MyApp

# Analyze with JDK Mission Control
jmc recording.jfr
```

### async-profiler

```bash
# Attach to running process
./profiler.sh -d 30 -f profile.html <pid>

# Profile specific events
./profiler.sh -e cpu -d 30 -f cpu.html <pid>
./profiler.sh -e alloc -d 30 -f alloc.html <pid>
```

### VisualVM

- Connect to local or remote JVM
- CPU and memory profiling
- Heap dump analysis
- Thread monitoring

### JVM Flags for Profiling

```bash
# GC logging
-Xlog:gc*:file=gc.log

# Heap dump on OOM
-XX:+HeapDumpOnOutOfMemoryError
-XX:HeapDumpPath=/tmp/heapdump.hprof

# JFR continuous recording
-XX:StartFlightRecording=disk=true,maxsize=500m
```

## Go Profiling

### pprof (Built-in)

```go
import (
    "net/http"
    _ "net/http/pprof"
)

func main() {
    go func() {
        http.ListenAndServe("localhost:6060", nil)
    }()
    // ... application code ...
}
```

```bash
# CPU profile
go tool pprof http://localhost:6060/debug/pprof/profile?seconds=30

# Memory profile
go tool pprof http://localhost:6060/debug/pprof/heap

# Interactive commands
(pprof) top10
(pprof) web        # Generates SVG
(pprof) list main  # Source annotation
```

### Benchmarking

```go
func BenchmarkFunction(b *testing.B) {
    for i := 0; i < b.N; i++ {
        FunctionToTest()
    }
}
```

```bash
# Run benchmarks with profiling
go test -bench=. -cpuprofile=cpu.prof -memprofile=mem.prof
```

## Browser/Frontend Profiling

### Chrome DevTools

**Performance Tab**
- Record page load or interaction
- Flame chart of JavaScript execution
- Layout and paint events
- Memory timeline

**Memory Tab**
- Heap snapshots
- Allocation timeline
- Allocation profiler

**Lighthouse**
- Core Web Vitals measurement
- Performance scoring
- Optimization recommendations

### Key Metrics

| Metric | Target | Description |
|--------|--------|-------------|
| LCP | < 2.5s | Largest Contentful Paint |
| FID | < 100ms | First Input Delay |
| CLS | < 0.1 | Cumulative Layout Shift |
| INP | < 200ms | Interaction to Next Paint |
| TTFB | < 800ms | Time to First Byte |

### Bundle Analysis

**webpack-bundle-analyzer**
```bash
npm install --save-dev webpack-bundle-analyzer

# In webpack config
const BundleAnalyzerPlugin = require('webpack-bundle-analyzer').BundleAnalyzerPlugin;
plugins: [new BundleAnalyzerPlugin()]
```

**source-map-explorer**
```bash
npm install -g source-map-explorer
source-map-explorer bundle.js
```

## Database Profiling

### PostgreSQL

**EXPLAIN ANALYZE**
```sql
EXPLAIN (ANALYZE, BUFFERS, FORMAT TEXT)
SELECT * FROM users WHERE email = 'user@example.com';
```

**pg_stat_statements**
```sql
-- Enable extension
CREATE EXTENSION pg_stat_statements;

-- View slow queries
SELECT query, calls, total_time, mean_time
FROM pg_stat_statements
ORDER BY total_time DESC
LIMIT 20;
```

**pg_stat_user_tables**
```sql
-- Check index usage
SELECT relname, seq_scan, idx_scan, n_tup_ins, n_tup_upd
FROM pg_stat_user_tables;
```

### MySQL

**EXPLAIN**
```sql
EXPLAIN FORMAT=JSON SELECT * FROM users WHERE email = 'user@example.com';
```

**Performance Schema**
```sql
-- Enable statement history
UPDATE performance_schema.setup_consumers
SET enabled = 'YES'
WHERE name = 'events_statements_history';

-- Query digest
SELECT DIGEST_TEXT, COUNT_STAR, AVG_TIMER_WAIT/1000000000 as avg_ms
FROM performance_schema.events_statements_summary_by_digest
ORDER BY AVG_TIMER_WAIT DESC
LIMIT 10;
```

**Slow Query Log**
```ini
# my.cnf
slow_query_log = 1
slow_query_log_file = /var/log/mysql/slow.log
long_query_time = 1
```

### MongoDB

**explain()**
```javascript
db.collection.find({email: "user@example.com"}).explain("executionStats")
```

**Profiler**
```javascript
// Enable profiling for slow queries (> 100ms)
db.setProfilingLevel(1, { slowms: 100 })

// View profiler output
db.system.profile.find().sort({ ts: -1 }).limit(5)
```

## System Profiling

### Linux Tools

**perf**
```bash
# CPU profiling
perf record -g ./app
perf report

# Flame graph generation
perf script | ./stackcollapse-perf.pl | ./flamegraph.pl > perf.svg
```

**htop**
- Interactive process viewer
- CPU, memory per process
- Thread view

**iotop**
```bash
# Disk I/O per process
iotop -o  # Only show processes doing I/O
```

**nethogs**
```bash
# Network usage per process
nethogs eth0
```

**vmstat**
```bash
# System resource statistics
vmstat 1  # Every 1 second
```

### Tracing

**strace**
```bash
# System calls
strace -c ./app          # Summary
strace -e open ./app     # Specific syscalls
strace -p <pid>          # Attach to process
```

**ltrace**
```bash
# Library calls
ltrace ./app
```

## Distributed Tracing

### OpenTelemetry

Standard for traces, metrics, and logs across services.

| Tool | Description |
|------|-------------|
| Jaeger | Open-source distributed tracing |
| Zipkin | Distributed tracing system |
| Tempo | Grafana's tracing backend |
| Honeycomb | Observability platform |

### Key Concepts

- **Trace**: End-to-end request journey
- **Span**: Individual operation within a trace
- **Context propagation**: Passing trace IDs across services

## Load Testing Tools

| Tool | Best For |
|------|----------|
| k6 | JavaScript-based, CI-friendly |
| Apache JMeter | Complex scenarios, GUI |
| Locust | Python-based, distributed |
| Gatling | Scala-based, detailed reports |
| wrk | Simple HTTP benchmarking |
| ab (ApacheBench) | Quick HTTP benchmarks |

### k6 Example

```javascript
import http from 'k6/http';
import { check, sleep } from 'k6';

export const options = {
  vus: 100,
  duration: '30s',
};

export default function () {
  const res = http.get('https://api.example.com/');
  check(res, { 'status was 200': (r) => r.status === 200 });
  sleep(1);
}
```

## External Resources

- [Brendan Gregg's Performance Tools](http://www.brendangregg.com/linuxperf.html) - Linux performance analysis
- [Flame Graphs](http://www.brendangregg.com/flamegraphs.html) - Visualization technique
- [web.dev Performance](https://web.dev/performance/) - Web performance guides
- [USE Method](http://www.brendangregg.com/usemethod.html) - Resource analysis methodology
- [Node.js Performance Guide](https://nodejs.org/en/docs/guides/dont-block-the-event-loop/) - Event loop optimization
