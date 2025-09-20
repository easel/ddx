---
tags: [adr, architecture, search, discovery, indexing, ddx]
template: false
version: 1.0.0
---

# ADR-006: Asset Discovery and Search Architecture

**Date**: 2025-01-14
**Status**: Proposed
**Deciders**: DDX Development Team
**Technical Story**: Design the system for discovering, indexing, and searching DDX assets across local and remote repositories

## Context

### Problem Statement
DDX needs an efficient asset discovery and search system that:
- Enables fast discovery of relevant templates, patterns, and prompts
- Supports multiple search strategies (keyword, semantic, tag-based)
- Works offline with cached assets
- Scales to thousands of assets without performance degradation
- Provides relevant results ranked by quality and applicability
- Maintains search index consistency across updates

### Forces at Play
- **Speed**: Search must return results in < 500ms
- **Relevance**: Results must be accurately ranked by usefulness
- **Offline Support**: Must work without network connectivity
- **Scalability**: Must handle growing asset libraries efficiently
- **Freshness**: Index must stay synchronized with asset changes
- **Storage**: Index size must remain reasonable
- **Complexity**: Search logic must be maintainable

### Constraints
- No external database dependencies
- Index must be portable across platforms
- Search must work with partial/fuzzy matches
- Must support incremental index updates
- Cannot require significant memory (< 100MB)
- Must integrate with git-based distribution

## Decision

### Chosen Approach
Implement a hybrid search architecture combining:
1. **Local file-based index** using BoltDB for metadata storage
2. **Full-text search** using Bleve search library
3. **Tag-based categorization** with hierarchical taxonomy
4. **Smart caching** with TTL and lazy updates
5. **Ranking algorithm** based on multiple relevance factors

### Architecture Components

```
┌─────────────────────────────────────┐
│         Search Interface            │
├─────────────────────────────────────┤
│         Query Parser                │
├──────────┬──────────┬───────────────┤
│  Text    │   Tag    │   Metadata    │
│  Search  │  Filter  │   Filter      │
├──────────┴──────────┴───────────────┤
│       Ranking Engine                │
├─────────────────────────────────────┤
│       Index Manager                 │
├──────────┬──────────────────────────┤
│  BoltDB  │      Bleve               │
│  Store   │   Text Index             │
└──────────┴──────────────────────────┘
```

### Index Structure
```go
type AssetIndex struct {
    ID          string    // Unique identifier
    Type        string    // template|pattern|prompt|config
    Name        string    // Human-readable name
    Path        string    // Filesystem path
    Description string    // Brief description
    Tags        []string  // Categorization tags
    Keywords    []string  // Extracted keywords
    Language    []string  // Programming languages
    Framework   []string  // Frameworks/libraries
    Metadata    map[string]interface{}
    Content     string    // Indexed content
    Popularity  int       // Usage count
    Quality     float64   // Quality score
    Updated     time.Time // Last update time
    Hash        string    // Content hash
}
```

### Search Query Language
```
# Simple search
ddx search "react component"

# Tag filtering
ddx search --tags "web,typescript" "authentication"

# Type filtering
ddx search --type template "nextjs"

# Language filtering
ddx search --language go "cli"

# Combined query
ddx search --tags security --language python "oauth"
```

### Ranking Algorithm
```
Score = (
    0.3 * TextRelevance +      # Full-text search score
    0.2 * TagMatch +            # Tag similarity score
    0.2 * Popularity +          # Usage frequency
    0.15 * Recency +            # How recently updated
    0.1 * QualityScore +        # Community quality rating
    0.05 * LanguageMatch        # Language preference match
)
```

### Rationale
- **BoltDB**: Embedded key-value store, no dependencies, fast
- **Bleve**: Pure Go full-text search, no CGO dependencies
- **File-based**: Portable, version-controllable, no server required
- **Hybrid Approach**: Combines structured and unstructured search
- **Local-first**: Works offline, syncs when connected

## Alternatives Considered

### Option 1: SQLite with FTS5
**Description**: Use SQLite with full-text search extension

**Pros**:
- Mature and battle-tested
- Excellent query capabilities
- ACID compliance
- Good performance
- Wide tool support

**Cons**:
- CGO dependency complicates builds
- Platform-specific FTS implementations
- Larger binary size
- More complex deployment

**Why rejected**: CGO dependency conflicts with simple distribution goal

### Option 2: Elasticsearch/OpenSearch
**Description**: Use external search service

**Pros**:
- Powerful search capabilities
- Excellent relevance tuning
- Scalable to millions of documents
- Rich query language
- Advanced analytics

**Cons**:
- Requires running service
- Complex deployment
- Network dependency
- Resource intensive
- Overkill for use case

**Why rejected**: External service dependency violates zero-dependency principle

### Option 3: Simple JSON Index
**Description**: Maintain search index as JSON files

**Pros**:
- Dead simple implementation
- Human readable
- Version controllable
- No dependencies
- Easy debugging

**Cons**:
- Poor performance at scale
- Limited query capabilities
- No full-text search
- Large memory footprint
- Slow updates

**Why rejected**: Insufficient performance and capability for production use

### Option 4: Sonic Search Backend
**Description**: Use Sonic lightweight search backend

**Pros**:
- Very fast search
- Low memory usage
- Simple protocol
- Good for autocomplete

**Cons**:
- Requires separate service
- Limited query capabilities
- No complex ranking
- Additional deployment complexity

**Why rejected**: External service requirement and limited features

### Option 5: Custom Trie-based Index
**Description**: Build custom prefix-trie search index

**Pros**:
- Optimized for our use case
- Fast prefix matching
- Low memory usage
- No dependencies

**Cons**:
- Significant development effort
- Limited to prefix search
- No fuzzy matching
- Maintenance burden
- Reinventing the wheel

**Why rejected**: Development effort not justified given available solutions

## Consequences

### Positive Consequences
- **Fast Local Search**: Sub-second results even offline
- **Rich Queries**: Support complex filtering and ranking
- **Scalability**: Handles thousands of assets efficiently
- **Portability**: Index travels with repository
- **Incremental Updates**: Only reindex changed assets
- **Extensibility**: Can add new ranking factors easily

### Negative Consequences
- **Index Maintenance**: Must keep index synchronized
- **Storage Overhead**: Index adds ~10-50MB to repository
- **Complexity**: Multiple components to maintain
- **Initial Indexing**: First index build takes time
- **Memory Usage**: Search operations use 50-100MB RAM

### Neutral Consequences
- **Update Lag**: Index updates not immediate
- **Ranking Tuning**: Requires iteration to optimize
- **Search Syntax**: Users must learn query language
- **Cache Management**: Periodic cleanup needed

## Implementation

### Required Changes
1. Integrate BoltDB for metadata storage
2. Integrate Bleve for full-text search
3. Build index management layer
4. Implement query parser
5. Create ranking engine
6. Build cache management
7. Add index maintenance commands
8. Create search UI/output formatting

### Index Lifecycle
```
1. Initial Index Build
   - Scan all asset directories
   - Extract metadata from files
   - Build full-text index
   - Store in BoltDB

2. Incremental Updates
   - Monitor file changes
   - Update affected entries
   - Reindex modified content
   - Update rankings

3. Search Operation
   - Parse query
   - Execute searches in parallel
   - Merge results
   - Apply ranking
   - Format output

4. Cache Management
   - TTL-based expiration
   - LRU eviction
   - Background refresh
   - Consistency checks
```

### Success Metrics
- **Search Speed**: < 500ms for 95% of queries
- **Index Size**: < 10% of asset size
- **Relevance**: > 80% user satisfaction with results
- **Index Build**: < 30s for 1000 assets
- **Memory Usage**: < 100MB during search
- **Cache Hit Rate**: > 70% for common queries

## Compliance

### Security Requirements
- Sanitize search queries to prevent injection
- Validate index integrity with checksums
- No execution of indexed content
- Secure handling of file paths
- Rate limiting for search operations

### Performance Requirements
- Search response < 500ms
- Index update < 100ms per file
- Memory usage < 100MB
- CPU usage < 50% during index
- Concurrent search support

### Regulatory Requirements
- No indexing of sensitive data
- Respect .gitignore patterns
- GDPR compliance for any metadata
- No telemetry without consent

## Monitoring and Review

### Key Indicators to Watch
- Search query performance distribution
- Index size growth rate
- Cache hit/miss ratios
- Memory usage patterns
- Failed query rate
- User search patterns

### Review Date
Q2 2025 - After initial user feedback

### Review Triggers
- Search performance degrades > 50%
- Index size exceeds 100MB
- User complaints about relevance
- Memory usage exceeds 200MB
- New search technology emerges

## Related Decisions

### Dependencies
- ADR-001: Defines asset types to index
- ADR-002: Git subtree affects index distribution
- ADR-003: Go implementation determines library choices
- ADR-005: Configuration defines search preferences

### Influenced By
- Modern search user expectations
- Offline-first architecture principles
- Go ecosystem library availability
- Performance requirements

### Influences
- Asset metadata standards
- Contribution quality guidelines
- CLI user interface design
- Caching strategies

## References

### Documentation
- [Bleve Search Library](https://blevesearch.com/)
- [BoltDB Documentation](https://github.com/etcd-io/bbolt)
- [Search Relevance Tuning](https://opensourceconnections.com/blog/2014/06/10/what-is-search-relevancy/)
- [Information Retrieval](https://nlp.stanford.edu/IR-book/)

### External Resources
- [Building Search Engines](https://www.algolia.com/blog/engineering/inside-the-algolia-engine/)
- [Search UX Best Practices](https://www.nngroup.com/articles/search-visible-and-simple/)
- [Indexing Strategies](https://www.elastic.co/guide/en/elasticsearch/guide/current/index-design.html)

### Discussion History
- Search requirements gathering session
- Performance benchmarking results
- User feedback on discovery pain points
- Evaluation of search libraries

## Notes

The search architecture aligns with DDX's medical metaphor - like a medical database that helps doctors quickly find relevant cases, treatments, and research. The ranking algorithm acts like clinical decision support, surfacing the most relevant information based on multiple factors.

Key insight: Local-first search with smart caching provides the best balance of performance, reliability, and simplicity. Users get fast results whether online or offline, while the system remains maintainable.

Implementation tip: Start with basic text search and gradually add sophistication based on user needs. Resist over-engineering the ranking algorithm before gathering real usage data.

The hybrid approach (structured + unstructured search) mirrors medical information retrieval - combining structured data (patient records) with unstructured text (clinical notes) for comprehensive results.

---

**Last Updated**: 2025-01-14
**Next Review**: 2025-04-14