# Zstd Compression Package

This package provides zstd compression support with runtime detection of libzstd availability.

## Features

- **Runtime Detection**: Automatically detects and uses libzstd if available
- **Dual Implementation**: Falls back to pure Go implementation when libzstd is not available
- **High Compression Levels**: Supports compression levels up to 22 when using libzstd
- **Parallel Compression**: Automatically uses multiple CPU cores for faster compression

## Parallel Compression

The zstd compression automatically uses multiple CPU cores for faster compression:

- **Default**: Uses 75% of physical CPU cores
- **Configuration**: Set `ZSTD_WORKERS` environment variable to override

### Examples

```bash
# Use 4 compression workers
export ZSTD_WORKERS=4

# Use single-threaded compression
export ZSTD_WORKERS=1

# Use automatic detection (default)
unset ZSTD_WORKERS
```

### Performance Considerations

- **Pure Go Implementation**: Supports parallel compression with multiple workers
- **Gozstd (libzstd)**: Supports parallel compression with multiple workers
- Only one layer is compressed at a time, so resource usage is controlled

### Memory Usage

Memory usage scales with the number of workers:

- **Pure Go**: ~128KB per worker (e.g., 16 workers = 2MB)
- **libzstd**: ~1MB per worker (e.g., 16 workers = 16MB)

For memory-constrained environments, limit workers:
```bash
export ZSTD_WORKERS=4
```

## Compression Levels

- **Pure Go**: Levels 0-11 (uses klauspost/compress)
- **libzstd**: Levels 0-22 (via gozstd wrapper)

The implementation automatically selects the best available option based on the requested compression level and library availability.

## Testing

Run tests with specific tags:

```bash
# Basic tests
go test ./pkg/compression/zstd

# Benchmark parallel compression
go test -bench=BenchmarkParallelCompression ./pkg/compression/zstd

# Run all zstd tests including stress tests
make test-zstd-all
```
