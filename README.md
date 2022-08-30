# sieve_date_offset_filter
Read an e-mail and exit with error status if the date is too far in the future.

## Why?
I've got a lot of spam where they put the Date header way in the future.
Perhaps this is to make sure the e-mail shows up at the top of whatever folder
it's in for a long time.

In any event, the date header should be set to the current time when the e-mail
is generated, so the receiving server should expect e-mails at the current time
if they were delivered very quickly, ranging to some time in the past, but
never from the future. This makes for a consistent way to detect this type of
spam.

Most of this spam I've seen has the header set at least one year in the future.
So, this program takes a fairly lenient approach of allowing e-mails up to
48-hours in the future. This is a wide enough margin to allow messages from
any legitimate senders who have misconfigured their server's timezone.

## Building
Only the Go standard library is used. Building should be as simple as:
```
go build .
```

## Usage
There are no command-line options. The program reads e-mail data including
headers from STDIN. It will exit as soon as a line with the `Date:` header
is read, or if EOF is reached. An exit status of `0` indicates that the date
is acceptable. An exit status of `1` indicates that the date is too far in
the future.

Note that in order to prevent an e-mail server incorrectly discarding messages,
the program will exit with status `0` if an internal error occurs. Information
about internal errors is output to STDOUT.

## Integrating with your e-mail server
This program is designed to be used with Dovecot Pigeonhole's `extprograms`
plugin, using an `execute :pipe` condition in your sieve file. Exact
configuration depends on your server.

More information on extprograms can be found in
https://doc.dovecot.org/configuration_manual/sieve/plugins/extprograms/

### Memory consideration
Dovecot places a memory limit on subprocesses, and the Go runtime may exceed
the default for this limit. If the process crashes when Dovecot runs it,
you may need to increase the `default_vsz_limit` in Dovecot's config.

## Testing
For testing convenience, you can paste some emails in the `test_data` directory
named like `bad1.eml`, `bad2.eml`, etc. and `good1.eml`, `good2.eml`, etc.
`test.sh` will run the program on all files in that directory, expecting all
the `good`... files to pass the date filter and the `bad`... files to fail it.
The script exits on the first failed expectation.