# count

[![CircleCI](https://circleci.com/gh/huttotw/count/tree/master.svg?style=svg)](https://circleci.com/gh/huttotw/count/tree/master)

Count is a tool to count high volume things and persisting the aggregate to a store at some interval.

**Potential use cases:**
- Counting the number page views every minute
- Counting the number of messages processed by a queue every second.

## Example
```go
db := newDatabaseWriter() // implements count.Writer
c := count.New(time.Minute * 5, db)

c.Increment("user:1", 1)
```

In this example, we set up count to write the number of times we have seen user 1 to the database every 5 minutes.
If your database writer is configured to write to `user_pageviews` table, your rows may look something like this:

| user_id | count | timestamp            |
|---------|-------|----------------------|
| 1       | 17    | 2019-03-02T15:30:00Z |
| 1       | 5     | 2019-03-02T15:35:00Z |

Once your data is in this form, you can run queries like: _how many times has user 1 visited in the past 10 minutes?_
Or in the context of message queues, _how many messages have I processed for this user in the past 10 minutes?_

## Benchmarks
The benchmark can be run with this command:
```bash
go test -v -run none -bench . -benchmem -benchtime 3s
```

| name                          | ops       | ns/op      | B/op   | allocs/op   |
|-------------------------------|-----------|------------|--------|-------------|
| BenchmarkCounter_Increment-12 | 100000000 | 49.8 ns/op | 0 B/op | 0 allocs/op |

## License

Copyright &copy; 2019 Trevor Hutto

Licensed under the Apache License, Version 2.0 (the "License"); you may not use this work except in compliance with the License. You may obtain a copy of the License in the LICENSE file, or at:

http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the specific language governing permissions and limitations under the License.
